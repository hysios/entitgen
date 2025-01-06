package main

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/hysios/entitgen/gen"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

type EntitGen struct {
	PkgName      string
	ProjectDir   string
	Output       string
	ProtobufType string
	ModelType    string
	ProtoGenPath string
	ProtoPath    string
	Options      Options
	PbPkgName    string
	Debug        bool
}

func (g *EntitGen) init() {
	if g.Debug {
		log = initLog(-1)
	}
}

func (g *EntitGen) Run() error {
	g.init()

	protoPaths := strings.Split(g.ProtoGenPath, ",")

	pkgs, err := g.parsePackage(protoPaths...)
	if err != nil {
		return err
	}

	typeinfos, err := g.buildTypeInfos(pkgs)
	if err != nil {
		return err
	}

	modelPkgs, err := g.parsePackage(g.Output)
	if err != nil {
		return err
	}

	modelTypesInfos, err := g.buildTypeInfos(modelPkgs)
	if err != nil {
		return err
	}

	typ, ok := typeinfos[g.ProtobufType]
	if !ok {
		return fmt.Errorf("not found protobuf type %s", g.ProtobufType)
	}

	mtyp, ok := modelTypesInfos[g.ModelType]
	var modelFile string
	if ok && g.Options.EntityFileDetect {
		modelFile = mtyp.GetFilename()
		// modelFile has .entity.go suffix
		if !strings.HasSuffix(modelFile, ".entity.go") {
			return fmt.Errorf("found model type `%s` in non .entity.go file: %s", g.ModelType, modelFile)
		}
	} else {
		// create new model file
		modelFile = g.Output + "/" + g.ModelType + ".entity.go"
	}

	gens, err := g.convertPbToModel(typ, typeinfos, modelTypesInfos, g.ModelType)
	if err != nil {
		return err
	}

	for typ, gentyp := range gens {
		modelFile = filepath.Join(g.Output, strings.ToLower(typ)+".entity.go")
		_ = os.MkdirAll(filepath.Dir(modelFile), 0755)
		log.Debugf("write model to %s", modelFile)
		f, err := os.Create(modelFile)
		if err != nil {
			return err
		}

		if err := g.writeTypeTo(gentyp, f); err != nil {
			return err
		}
	}

	return exec.Command("goimports", "-w", ".").Run()
}

// parsePackage parses the package in the given directory.
func (g *EntitGen) parsePackage(dir ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedExportFile |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes,
	}

	dirs := readSubDir(dir)

	pkgs, err := packages.Load(cfg, dirs...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		// return nil, err
	}

	return pkgs, nil
}

func readSubDir(dirs []string) []string {
	var (
		subDirs []string
		rel     bool
		err     error
	)
	for _, dir := range dirs {
		rel = strings.HasPrefix(dir, ".")
		if rel {
			dir, err = filepath.Abs(dir)
			if err != nil {
				continue
			}
		}
		// subDirs = append(subDirs, dir)
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

			if info.IsDir() {
				subDirs = append(subDirs, path)
			}
			return nil
		})
	}
	return subDirs
}

// buildTypeInfos builds the typeinfos for the given packages.
func (g *EntitGen) buildTypeInfos(pkgs []*packages.Package) (typeinfos map[string]*typeInfo, err error) {
	var (
		// namesByType   typeutil.Map // value is []string
		funcsByType   typeutil.Map // value is []string
		structsByType typeutil.Map
		// othersByType  typeutil.Map
		// structs   []*typeInfo
		nameTypes []*typeInfo
	)

	typeinfos = make(map[string]*typeInfo)

	for _, pkg := range pkgs {

		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			if unicode.IsUpper(([]rune)(name)[0]) {
				var (
					object = scope.Lookup(name)
					T      = object.Type()
				)
				log.Debugf("object %s by %T", name, T)
				switch T.(type) {
				case *types.Signature:
					names, _ := funcsByType.At(T).([]string)
					log.Debugf("func %s %T", name, T)
					names = append(names, name)
					funcsByType.Set(T, names)
				case *types.Struct:
					names, _ := structsByType.At(T).([]string)
					names = append(names, name)
					structsByType.Set(T, names)
					nameTypes = append(nameTypes, &typeInfo{
						pkg:     object.Pkg(),
						name:    name,
						pkgName: pkg.Name,
						// alias: getAlias(object.Pkg()),
						typ: T,
					})
				case *types.Named:
				NameAgain:
					switch x := T.Underlying().(type) {
					case *types.Signature:
						names, _ := funcsByType.At(T).([]string)
						names = append(names, name)
						log.Debugf("func %s %T", name, T)
						funcsByType.Set(T, names)
					case *types.Struct:
						names, _ := structsByType.At(T).([]string)
						names = append(names, name)
						structsByType.Set(x, names)
						typeinfos[name] = &typeInfo{
							pkg:     object.Pkg(),
							pkgName: pkg.Name,
							name:    name,
							typ:     x,
							scope:   scope,
							methods: typeutil.IntuitiveMethodSet(T, nil),
						}
					case *types.Named:
						T = x.Underlying()
						goto NameAgain
					case *types.Slice:
						log.Debugf("slice type: %s %T", x, x)
					case *types.Interface:
					case *types.Basic:
						log.Debugf("alias type: %s %T", x, x)
						names, _ := structsByType.At(T).([]string)
						names = append(names, name)
						structsByType.Set(x, names)
						typeinfos[name] = &typeInfo{
							pkg:     object.Pkg(),
							pkgName: pkg.Name,
							name:    name,
							typ:     x,
							scope:   scope,
							methods: typeutil.IntuitiveMethodSet(T, nil),
						}
					default:
						log.Debugf("unknown type: %s %T", x, x)
					}
				case *types.Map:
				default:
					log.Debugf("unknown type: %s", T)
				}
			}
		}
	}

	return
}

func (g *EntitGen) convertPbToModel(typ *typeInfo, globalInfos map[string]*typeInfo, modelInfos map[string]*typeInfo, outType string) (gens map[string]*GenType, err error) {
	// typ.

	var (
		typeInfos []*typeInfo
		genType   *GenType
		cycleIdx  = make(map[*typeInfo]bool)
		modelType = outType
	)
	// typeInfos = append(typeInfos, typ)
	gens = make(map[string]*GenType)

	genOne := func(typ *typeInfo, modelType string) (*GenType, error) {
		cycleIdx[typ] = true
		genType = &GenType{
			Name:      typ.name,
			ModelName: modelType,
			PkgName:   filepath.Base(g.Output),
			PbPkgName: typ.pkgName,
			Options:   g.Options,
			models:    modelInfos,
		}

		genType.AddImport(gen.Pkg{
			Fullname: typ.pkg.Path(),
			Alias:    typ.pkgName,
		})

		for _, field := range typ.Fields() {
			if !field.Exported() {
				continue
			}

			n := field.Name()
			_ = n
			if g.suppressInclude(field.Name(), modelType) {
				// if g.Options.Suppress.Include(field.Name(), modelType) {
				continue
			}

			modTyp, found := modelInfos[modelType]
			var modFieldTyp *types.Var
			if found {
				if modFieldTyp, found = modTyp.FieldByName(field.Name()); !found {
					log.Infof("not found field %s in model %s", field.Name(), modelType)
					continue
				}
			}

			var (
				convType string
				isPtr    bool = isPointer(field.Type())
				ptr      bool
				ok       bool
			)

			_ = isPtr
			switch {
			case isScalarType(field.Type().String()):
				convType, ok = conventionType(field, modFieldTyp)
				if ok {
					if strings.HasPrefix(convType, "sql.") { // sql Null Values
						convKey := fieldkey(field.Name(), convType)
						fieldType := field.Type().String()
						genType.AddConv(convKey, gen.TypeConv(
							fromNullType(fieldType),
							toNullType(fieldType),
						))
						genType.AddImport(gen.Pkg{
							Fullname: "github.com/hysios/entitgen/null",
						})
						genType.AddImport(gen.Pkg{
							Fullname: "database/sql",
						})
					} else if convType == "decimal.Decimal" {
						convKey := fieldkey(field.Name(), convType)
						genType.AddConv(convKey, &externConvert{
							toProto: func(in string) string {
								return "decimal.NewFromFloat(" + in + ")"
							},
							fromProto: func(in string) string {
								return in + ".InexactFloat64()"
							},
						})
						genType.AddImport(gen.Pkg{
							Fullname: "github.com/shopspring/decimal",
						})
					} else {
						convKey := fieldkey(field.Name(), convType)
						genType.AddConv(convKey, gen.TypeConv(field.Type().String(), convType))
					}
				}
			case isSliceType(field.Type()):
				// convType = field.Type().String()

				elemType := field.Type().(*types.Slice).Elem()
				if g.isEmbeddedField(field.Name()) {
					convType = gormSliceType(elemType)
					genType.AddImport(gen.Pkg{
						Fullname: "gorm.io/datatypes",
					})
				} else {
					if typ, ok := convertModelType(modelInfos, elemType); ok {
						convType = "[]*" + typ.name
						conv := getSliceProtoConv(genType, typ.name)
						convKey := fieldkey(field.Name(), convType)
						genType.AddConv(convKey, conv)
					} else {
						convType = getSliceType(elemType)
					}
				}
			case isMapType(field.Type()):
				key, value := getMapType(field.Type())
				// convType = "map[" + key.String() + "]" + getTypeName(value)
				convType = gormMapType(key, value)
				genType.AddImport(gen.Pkg{
					Fullname: "gorm.io/datatypes",
				})
			case isExternalType(field.Type().String(), modFieldTyp):
				extTyp, ok := getExternalType(field.Type().String(), modFieldTyp)
				if !ok {
					return nil, fmt.Errorf("not found external type: %s", field.Type().String())
				}
				genType.AddImport(gen.Pkg{
					Fullname: extTyp.Type.PkgName,
				})
				convType = extTyp.Type.PureType()
				ptr = extTyp.Type.Pointer
				conv := getExternalConvert(field.Type().String(), modFieldTyp)
				if conv == nil {
					return nil, fmt.Errorf("not found convert: %s", field.Type().String())
				}

				convKey := fieldkey(field.Name(), convType)
				genType.AddConv(convKey, conv)
				for _, imp := range extTyp.Imports {
					genType.AddImport(imp)
				}
			case isAliasType(field.Type()):
				t, ok := globalInfos[pureName(getTypeName(field.Type()))]
				if !ok {
					convType = field.Type().Underlying().String()
					// fromType := field.Type().String()
					convKey := fieldkey(field.Name(), convType)
					genType.AddConv(convKey, gen.TypeConv(getAliasType(field.Type()), convType))
				} else {
					if checkEnumMethod(t.methods) && modFieldTyp != nil {
						convType = field.Type().Underlying().String()
						if kind(modFieldTyp.Type().String()) != kind(convType) {
							convKey := fieldkey(field.Name(), convType)
							genType.AddConv(convKey, newEnumMapConv(getAliasType(field.Type()), convType))
						} else {
							convType = field.Type().Underlying().String()
							// fromType := field.Type().String()
							convKey := fieldkey(field.Name(), convType)
							genType.AddConv(convKey, gen.TypeConv(getAliasType(field.Type()), convType))
						}
					} else {
						convType = field.Type().Underlying().String()
						// fromType := field.Type().String()
						convKey := fieldkey(field.Name(), convType)
						genType.AddConv(convKey, gen.TypeConv(getAliasType(field.Type()), convType))
					}
				}
			case isStructType(field.Type()):
				// convType = field.Type().String()
				name, stucTyp := getStructType(field.Type())
				_, _ = name, stucTyp
				// convType = getTypeName(field.Type())
				convType = name
				ptr = true
				typeInfos = append(typeInfos, globalInfos[name])
				// modelType = name
				convKey := fieldkey(field.Name(), convType)
				genType.AddConv(convKey, gen.ProtoConv(name))
			default:
				convType = field.Type().String()
				// return nil, fmt.Errorf("unknown type: %s", field.Type().String())
			}

			genType.Fields = append(genType.Fields, &gen.Field{
				ID:      field.Name(),
				Name:    convertName(field.Name()),
				PbName:  field.Name(),
				Pkg:     field.Pkg().Path(),
				PbType:  field.Type().String(),
				Type:    convType,
				Pointer: ptr,
			})

		}

		return genType, nil
	}

	genType, err = genOne(typ, modelType)
	if err != nil {
		return nil, err
	}
	gens[modelType] = genType

	for i := 0; i < len(typeInfos); i++ {
		typ := typeInfos[i]
		if typ == nil {
			continue
		}

		if cycleIdx[typ] {
			continue
		}

		genType, err = genOne(typ, typ.name)
		if err != nil {
			return nil, err
		}

		gens[typ.name] = genType
	}
	return
}

func getTypeName(value types.Type) string {
	switch x := value.(type) {
	case *types.Basic:
		return x.String()
	case *types.Named:
		return x.Obj().Pkg().Name() + "." + x.Obj().Name()
	case *types.Pointer:
		return "*" + getTypeName(x.Elem())
	case *types.Slice:
		return "[]" + getTypeName(x.Elem())
	case *types.Map:
		return "map[" + getTypeName(x.Key()) + "]" + getTypeName(x.Elem())
	case *types.Struct:
		n, _ := getStructType(x.Underlying())
		return n
	default:
		log.Debugf("unknown type: %s %T", value, value)
		return ""
	}
}

func getAliasType(typ types.Type) string {
	obj := typ.(*types.Named).Obj()
	return obj.Pkg().Name() + "." + obj.Name()
}

// isPointer
func isPointer(typ types.Type) bool {
	_, ok := typ.(*types.Pointer)
	return ok
}

// writeTypeTo writes the type to the given file.
func (g *EntitGen) writeTypeTo(typ *GenType, w io.Writer) error {
	tmpl, err := template.New("type").Parse(typeTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, typ)
}

// suppressInclude
func (g *EntitGen) suppressInclude(field, model string) bool {
	for _, sup := range g.Options.Suppress {
		if sup.Model != "" {
			if sup.Field == field && sup.Model == model {
				return true
			}
		} else {
			if sup.Field == field {
				return true
			}
		}
	}
	return false
}

// isEmbeddedField
func (g *EntitGen) isEmbeddedField(field string) bool {
	if g.Options.NoEmbed == nil {
		return true
	}

	for _, noembed := range g.Options.NoEmbed {
		if noembed == field {
			return false
		}
	}
	return true
}

// GetFilename returns the filename of the type.
func (t *typeInfo) GetFilename() string {
	var b bytes.Buffer
	t.scope.WriteTo(&b, 0, false)
	log.Debugf("scope %s", b.String())
	return ""
}

func fieldkey(field, typ string) string {
	return field + "_" + typ
}

// nullType returns the null type of the given type.
func nullType(typ string) string {
	switch typ {
	case "*string", "string":
		return "sql.NullString"
	case "*int", "*int32", "*int64", "int":
		return "sql.NullInt64"
	case "*uint", "*uint32", "*uint64", "uint":
		return "sql.NullInt64"
	case "*float64", "*float32", "float64":
		return "sql.NullFloat64"
	case "*bool", "bool":
		return "sql.NullBool"
	case "*byte", "byte":
		return "sql.NullByte"
	case "uint32", "uint64":
		return "sql.NullInt64"
	case "int32", "int64":
		return "sql.NullInt64"
	default:
		return ""
	}
}

// toNullType returns the null type of the given type.
func toNullType(typ string) string {
	return "null.ToNull" + strings.TrimPrefix(nullType(typ), "sql.Null")
}

// fromNullType returns the null type of the given type.
func fromNullType(typ string) string {
	return "null.NullTo" + strings.TrimPrefix(nullType(typ), "sql.Null")
}

// convertModelType
func convertModelType(infos map[string]*typeInfo, typ types.Type) (*typeInfo, bool) {
	typName := pureName(getTypeName(typ))

	if tt, ok := infos[typName]; ok {
		return tt, true
	}

	return nil, false
}

// pureName
func pureName(name string) string {
	ss := strings.Split(name, ".")
	return ss[len(ss)-1]
}

func checkMethod(methods []*types.Selection, name, retur string) bool {
	for _, m := range methods {
		fun, ok := m.Obj().(*types.Func)
		if !ok {
			continue
		}

		if fun.Name() != name {
			return false
		}

		sig, ok := fun.Type().(*types.Signature)
		if !ok {
			continue
		}

		ret := sig.Results().At(0).String()
		if strings.TrimSpace(strings.TrimPrefix(ret, "var ")) == retur {
			return true
		}
	}
	return false
}

func checkEnumMethod(methods []*types.Selection) bool {
	return checkMethod(methods, "Descriptor", "google.golang.org/protobuf/reflect/protoreflect.EnumDescriptor")
}

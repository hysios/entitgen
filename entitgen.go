package main

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
	"log"
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
	Output       string
	ProtobufType string
	ModelType    string
	ProtoGenPath string
	ProtoPath    string
	Options      Options
	PbPkgName    string
}

func (g *EntitGen) Run() error {
	protoPaths := strings.Split(g.ProtoGenPath, ";")

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
		log.Printf("write model to %s", modelFile)
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
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps | packages.NeedExportFile | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypesSizes}
	pkgs, err := packages.Load(cfg, dir...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, err
	}

	return pkgs, nil
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
				switch T.(type) {
				case *types.Signature:
					names, _ := funcsByType.At(T).([]string)
					names = append(names, name)
					funcsByType.Set(T, names)
				case *types.Struct:
					names, _ := structsByType.At(T).([]string)
					names = append(names, name)
					structsByType.Set(T, names)
					nameTypes = append(nameTypes, &typeInfo{
						pkg:  object.Pkg(),
						name: name,
						// alias: getAlias(object.Pkg()),
						typ: T,
					})
				case *types.Named:
					switch x := T.Underlying().(type) {
					case *types.Signature:
						names, _ := funcsByType.At(T).([]string)
						names = append(names, name)
						funcsByType.Set(T, names)

					case *types.Struct:
						names, _ := structsByType.At(T).([]string)
						names = append(names, name)
						structsByType.Set(T, names)
						typeinfos[name] = &typeInfo{
							pkg:   object.Pkg(),
							name:  name,
							typ:   x,
							scope: scope,
						}

					default:
						log.Printf("unknown type: %s %T", T, T)
					}
				case *types.Map:
				default:
					log.Printf("unknown type: %s", T)
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
	typeInfos = append(typeInfos, typ)
	gens = make(map[string]*GenType)

	for i := 0; i < len(typeInfos); i++ {
		typ := typeInfos[i]
		if cycleIdx[typ] {
			continue
		}
		cycleIdx[typ] = true

		genType = &GenType{
			Name:      modelType,
			PkgName:   filepath.Base(g.Output),
			PbPkgName: "pb",
			Options:   g.Options,
		}

		genType.AddImport(gen.Pkg{
			Fullname: typ.pkg.Path(),
			Alias:    "pb",
		})

		for _, field := range typ.Fields() {
			if !field.Exported() {
				continue
			}

			if g.suppressInclude(field.Name(), modelType) {
				// if g.Options.Suppress.Include(field.Name(), modelType) {
				continue
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
				convType, ok = conventionType(field.Type().String())
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
					} else {
						convKey := fieldkey(field.Name(), convType)
						genType.AddConv(convKey, gen.TypeConv(field.Type().String(), convType))
					}
				}
			case isSliceType(field.Type()):
				// convType = field.Type().String()
				elemType := field.Type().(*types.Slice).Elem()
				convType = gormSliceType(elemType)
				genType.AddImport(gen.Pkg{
					Fullname: "gorm.io/datatypes",
				})
			case isMapType(field.Type()):
				key, value := getMapType(field.Type())
				// convType = "map[" + key.String() + "]" + getTypeName(value)
				convType = gormMapType(key, value)
				genType.AddImport(gen.Pkg{
					Fullname: "gorm.io/datatypes",
				})
				// convType = field.Type().String()
			case isExternalType(field.Type().String()):
				extTyp, ok := getExternalType(field.Type().String())
				if !ok {
					return nil, fmt.Errorf("not found external type: %s", field.Type().String())
				}

				genType.AddImport(gen.Pkg{
					Fullname: extTyp.Type.PkgName,
				})
				convType = extTyp.Type.PureType()
				ptr = extTyp.Type.Pointer
				conv := getExternalConvert(field.Type().String())
				if conv == nil {
					return nil, fmt.Errorf("not found convert: %s", field.Type().String())
				}

				convKey := fieldkey(field.Name(), convType)
				genType.AddConv(convKey, conv)
				for _, imp := range extTyp.Imports {
					genType.AddImport(imp)
				}
			case isAliasType(field.Type()):
				convType = field.Type().Underlying().String()
				// fromType := field.Type().String()
				convKey := fieldkey(field.Name(), convType)
				genType.AddConv(convKey, gen.TypeConv(getAliasType(field.Type()), convType))
			case isStructType(field.Type()):
				// convType = field.Type().String()
				name, stucTyp := getStructType(field.Type())
				_, _ = name, stucTyp
				// convType = getTypeName(field.Type())
				convType = name
				ptr = true
				typeInfos = append(typeInfos, globalInfos[name])
				modelType = name
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
	default:
		log.Printf("unknown type: %s %T", value, value)
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

// GetFilename returns the filename of the type.
func (t *typeInfo) GetFilename() string {
	var b bytes.Buffer
	t.scope.WriteTo(&b, 0, false)
	log.Printf("scope %s", b.String())
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
	case "*float64", "*float32", "float64":
		return "sql.NullFloat64"
	case "*bool", "bool":
		return "sql.NullBool"
	case "*byte", "byte":
		return "sql.NullByte"
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

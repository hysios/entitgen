package main

import (
	"fmt"
	"go/types"
	"regexp"
	"strings"
	"sync"

	"github.com/hysios/entitgen/gen"
)

const SlicePkg = "github.com/akrennmair/slice"

type ExternalType struct {
	Imports []gen.Pkg
	Type    gen.Var
	Convert gen.Converter
}

// GenType
func (ext *ExternalType) GenType(modType *types.Var) ExternalType {
	out := ExternalType{
		Imports: ext.Imports,
		Type:    ext.Type,
		Convert: ext.Convert,
	}

	g := genericReg.FindAllStringSubmatch(modType.Type().String(), -1)
	out.Type = parseType(g[0][1])
	return out
}

// parseType
func parseType(typ string) gen.Var {
	var isPointer = strings.HasPrefix(typ, "*")
	typ = strings.TrimPrefix(typ, "*")

	after := typ[strings.Index(typ, "/"):]
	ss := strings.Split(after, ".")
	// typ = strings.TrimPrefix(typ, "[")
	return gen.Var{
		PkgName: typ[:strings.Index(typ, "/")] + ss[0],
		Pointer: isPointer,
		Type:    ss[1],
	}
}

var (
	// ExternalTypes
	externalTypes = map[string]ExternalType{
		"google.golang.org/protobuf/types/known/timestamppb.Timestamp": ExternalType{
			Imports: []gen.Pkg{
				gen.Pkg{
					Fullname: "google.golang.org/protobuf/types/known/timestamppb",
				},
			},
			Type: gen.Var{
				Name:    "t",
				Type:    "Time",
				PkgName: "time",
			},
			Convert: &externConvert{
				toProto: func(in string) string {
					return in + ".AsTime()"
				},
				fromProto: func(in string) string {
					return "timestamppb.New(" + in + ")"
				},
			},
		},
		"google.golang.org/protobuf/types/known/timestamppb.Timestamp=>time.Time": ExternalType{
			Imports: []gen.Pkg{
				gen.Pkg{
					Fullname: "google.golang.org/protobuf/types/known/timestamppb",
				},
			},
			Type: gen.Var{
				Name:    "t",
				Type:    "Time",
				PkgName: "time",
			},
			Convert: &externConvert{
				toProto: func(in string) string {
					return in + ".AsTime()"
				},
				fromProto: func(in string) string {
					return "timestamppb.New(" + in + ")"
				},
			},
		},
		"google.golang.org/protobuf/types/known/timestamppb.Timestamp=>database/sql.NullTime": ExternalType{
			Imports: []gen.Pkg{
				gen.Pkg{
					Fullname: "google.golang.org/protobuf/types/known/timestamppb",
				},
				gen.Pkg{
					Fullname: "github.com/hysios/entitgen/null",
				},
			},
			Type: gen.Var{
				Name:    "t",
				Type:    "NullTime",
				PkgName: "sql",
			},
			Convert: &externConvert{
				toProto: func(in string) string {
					return "null.PbtimeToSQLTime(" + in + ")"
				},
				fromProto: func(in string) string {
					return "null.SQLTimeToPbtime(" + in + ")"
				},
			},
		},
		"gorm.io/datatypes.JSONType[*]": ExternalType{
			Imports: []gen.Pkg{
				{
					Fullname: "gorm.io/datatypes",
				},
			},
			Convert: &externConvert{
				toProto: func(in string) string {
					return "datatypes.NewJSONType(" + in + ")"
				},
				fromProto: func(in string) string {
					return in + ".Data()"
				},
			},
		},
	}
)

type externConvert struct {
	fromProto func(in string) string
	toProto   func(in string) string
}

// From
func (c *externConvert) From(in string) string {
	return c.fromProto(in)
}

// To
func (c *externConvert) To(in string) string {
	return c.toProto(in)
}

var genericReg = regexp.MustCompile(`\[(.*)\]`)

func isExternalType(typ string, modelTyp *types.Var) bool {
	// if typ has star
	typ = strings.TrimPrefix(typ, "*")
	if _, ok := externalTypes[typ]; ok {
		return true
	}

	if modelTyp == nil {
		return false
	}

	modTyp := modelTyp.Type().String()
	// typ is generic type
	if genericReg.MatchString(modTyp) {
		modTyp = genericReg.ReplaceAllString(modTyp, "[*]")
		if _, ok := externalTypes[modTyp]; ok {

			return true
		}
	}
	return false
}

func getExternalType(typ string, modelTyp *types.Var) (ExternalType, bool) {
	typ = strings.TrimPrefix(typ, "*")
	if modelTyp != nil {
		typ += "=>" + modelTyp.Type().String()
	}
	t, ok := externalTypes[typ]
	if !ok {
		modTyp := modelTyp.Type().String()
		// typ is generic type
		if genericReg.MatchString(modTyp) {
			modTyp = genericReg.ReplaceAllString(modTyp, "[*]")
			t, ok = externalTypes[modTyp]
			if ok {
				t = t.GenType(modelTyp)
			}
		}
	}
	return t, ok
}

// getExternalConvert
func getExternalConvert(typ string, modelTyp *types.Var) gen.Converter {
	ext, ok := getExternalType(typ, modelTyp)
	if !ok {
		return nil
	}

	return ext.Convert
}

type sliceProtoConv struct {
	Type string
	sync.Once
}

func getSliceProtoConv(g *GenType, typ string) gen.Converter {
	c := &sliceProtoConv{
		Type: typ,
	}

	c.Once.Do(func() {
		g.AddImport(gen.Pkg{
			Fullname: SlicePkg,
		})
	})

	return c
}

// From implements gen.Converter.
func (s *sliceProtoConv) From(in string) string {
	return "slice.Map(" + in + ", " + s.Type + "ToProto)"
}

// To implements gen.Converter.
func (s *sliceProtoConv) To(in string) string {
	return "slice.Map(" + in + ", " + s.Type + "FromProto)"
}

type enumMapConv struct {
	Type     string
	alisType string
}

func newEnumMapConv(typ, alias string) *enumMapConv {
	return &enumMapConv{
		Type:     typ,
		alisType: alias,
	}
}

// From implements gen.Converter.
func (e *enumMapConv) From(in string) string {
	return fmt.Sprintf("%s(%s_value[%s])", e.Type, e.Type, in)
}

// pb.Role_name[int32(pUser.Role)],
// To implements gen.Converter.
func (e *enumMapConv) To(in string) string {
	return fmt.Sprintf("%s_name[%s(%s)]", e.Type, e.alisType, in)
}

var (
	_ gen.Converter = (*enumMapConv)(nil)
	_ gen.Converter = (*sliceProtoConv)(nil)
)

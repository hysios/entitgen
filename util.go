package main

import (
	"go/types"
	"strings"

	pluralize "github.com/gertd/go-pluralize"
)

var conventionTypes = map[string]string{
	"int32":    "int",
	"uint32":   "uint",
	"float32":  "float64",
	"double":   "float64",
	"*string":  "sql.NullString",
	"*uint":    "sql.NullInt64",
	"*uint32":  "sql.NullInt64",
	"*uint64":  "sql.NullInt64",
	"*int":     "sql.NullInt64",
	"*int32":   "sql.NullInt32",
	"*int64":   "sql.NullInt64",
	"*float32": "sql.NullFloat64",
	"*bool":    "sql.NullBool",
	"float64":  "decimal.Decimal",
}

func conventionType(protoType, modelType *types.Var) (string, bool) {
	var typ = protoType.Type().String()
	if modelType == nil {
		t, ok := conventionTypes[typ]
		if ok {
			return t, true
		}
		return typ, false
	}

	modTyp := getTypeName(modelType.Type())
	if modTyp == typ {
		return typ, false
	}

	t, ok := conventionTypes[typ]
	if ok {
		if t == modTyp {
			return t, true
		}
	}

	switch kind(typ) {
	case "int":
		switch kind(modTyp) {
		case "int", "uint":
			return modTyp, true
		case "NullInt":
			return modTyp, true
		default:
			return typ, false
		}
	case "uint":
		switch kind(modTyp) {
		case "int", "uint":
			return modTyp, true
		case "NullInt":
			return modTyp, true
		default:
			return typ, false
		}
	case "float":
		if kind(modTyp) == "float" {
			return modTyp, true
		}
		return typ, false
	}

	return typ, false
}

func kind(typ string) string {
	switch typ {
	case "int", "int8", "int16", "int32", "int64":
		return "int"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "uint"
	case "float32", "float64":
		return "float"
	case "bool":
		return "bool"
	case "string":
		return "string"
	case "byte", "rune":
		return "byte"
	case "sql.NullInt32", "sql.NullInt64":
		return "NullInt"
	case "sql.NullFloat64":
		return "NullFloat"
	case "sql.NullBool":
		return "NullBool"
	case "sql.NullString":
		return "NullString"
	}

	return ""
}

type matchFunc func(string) bool

type match struct {
	Match matchFunc
	Conv  func(string) string
}

var conventionNames = []match{}

func registerConv(m match) {
	conventionNames = append(conventionNames, m)
}

func convertName(name string) string {
	for _, m := range conventionNames {
		if m.Match(name) {
			return m.Conv(name)
		}
	}
	return name
}

var (
	plur = pluralize.NewClient()
)

func init() {
	registerConv(match{
		Match: func(name string) bool {
			return strings.HasSuffix(name, "Id")
		},
		Conv: func(name string) string {
			return name[:len(name)-2] + "ID"
		},
	})

}

// isScalaType returns true if typ is a scalar type.
func isScalarType(typ string) bool {
	if _, ok := conventionTypes[typ]; ok {
		return true
	}

	switch typ {
	case "string", "bool", "int", "int64", "uint", "uint64", "float64":
		return true
	case "int32", "uint32", "float32":
		return true
	}

	return false
}

func isSliceType(typ types.Type) bool {
	_, ok := typ.(*types.Slice)
	return ok
}

func isGormExtType(typ types.Type, modFieldTyp types.Type) bool {
	if modFieldTyp == nil {
		return false
	}

	typName := modFieldTyp.String()
	return strings.HasPrefix(typName, "gorm.io/datatypes.JSONType[")
}

func isMapType(typ types.Type) bool {
	_, ok := typ.(*types.Map)
	return ok
}

func getMapType(typ types.Type) (key types.Type, value types.Type) {
	m, ok := typ.(*types.Map)
	if !ok {
		return
	}
	key = m.Key()
	value = m.Elem()
	return
}

// isStructType returns true if typ is a struct type.
func isStructType(typ types.Type) bool {
	p, ok := typ.(*types.Pointer)
	if !ok {
		return false
	}
	n, ok := p.Elem().(*types.Named)
	if !ok {
		return false
	}
	_, ok = n.Underlying().(*types.Struct)
	return ok
	// _, ok := typ.(*types.Struct)
	// return ok
}

// getStructType
func getStructType(typ types.Type) (name string, stuc *types.Struct) {
	p, ok := typ.(*types.Pointer)
	if !ok {
		return
	}
	n, ok := p.Elem().(*types.Named)
	if !ok {
		return
	}
	name = n.Obj().Name()
	s, ok := n.Underlying().(*types.Struct)
	if !ok {
		return
	}
	stuc = s
	return
}

// isAliasType returns true if typ is a alias type.
func isAliasType(typ types.Type) bool {
	_, ok := typ.(*types.Named)
	return ok
}

// gormSliceType
func gormSliceType(typ types.Type) string {
	return "datatypes.JSONSlice[" + getTypeName(typ) + "]"
}

// getSliceType
func getSliceType(typ types.Type) string {

	return "[]" + getTypeName(typ)
}

// gormMapType
func gormMapType(key, value types.Type) string {
	mapType := "map[" + key.String() + "]" + getTypeName(value)
	return "datatypes.JSONType[" + mapType + "]"
}

func shortType(typ types.Type) string {
	t := typ.String()
	if strings.Contains(t, "/") {
		t = t[strings.LastIndex(t, "/")+1:]
	}
	return t
}

// pluralize
func plural(name string) string {
	return plur.Plural(name)
}

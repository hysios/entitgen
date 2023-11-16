package main

import (
	"go/ast"
	"go/types"
	"strings"
)

type typeInfo struct {
	pkg   *types.Package
	name  string
	alias string
	typ   types.Type
	scope *types.Scope
	funcs []*ast.FuncDecl
}

// Fields
func (t *typeInfo) Fields() []*types.Var {
	var fields []*types.Var
	if s, ok := t.typ.(*types.Struct); ok {
		for i := 0; i < s.NumFields(); i++ {
			fields = append(fields, s.Field(i))
		}
	}
	return fields
}

func (t *typeInfo) Tags() []string {
	var tags []string
	if s, ok := t.typ.(*types.Struct); ok {
		for i := 0; i < s.NumFields(); i++ {
			tags = append(tags, s.Tag(i))
		}
	}
	return tags
}

// Methods
func (t *typeInfo) Methods() []*types.MethodSet {
	var msets []*types.MethodSet
	celsius := t.pkg.Scope().Lookup(t.name).Type()
	for _, t := range []types.Type{celsius, types.NewPointer(celsius)} {
		mset := types.NewMethodSet(t)
		msets = append(msets, mset)
	}
	return msets
}

// FieldByName
func (t *typeInfo) FieldByName(name string) (*types.Var, bool) {
	if v, ok := fieldAlias[name]; ok {
		name = v
	}

	if alias, ok := hasSuffixAlias(name); ok {
		name = alias
	}

	if s, ok := t.typ.(*types.Struct); ok {
		for i := 0; i < s.NumFields(); i++ {
			if s.Field(i).Name() == name {
				return s.Field(i), true
			}
		}
	}

	return nil, false
}

var (
	fieldAlias = map[string]string{
		"Id": "ID",
	}
)

func hasSuffixAlias(name string) (string, bool) {
	for k, v := range fieldAlias {
		if strings.HasSuffix(name, k) {
			return strings.TrimSuffix(name, k) + v, true
		}
	}

	return "", false
}

func hasVariadic(fieldList *ast.FieldList) bool {
	if fieldList == nil {
		return false
	}

	for _, field := range fieldList.List {
		if field == nil {
			continue
		}

		if field.Names == nil {
			continue
		}

		for _, name := range field.Names {
			if name == nil {
				continue
			}

			if name.Name == "..." {
				return true
			}
		}
	}

	return false
}

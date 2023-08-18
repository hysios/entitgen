package main

import (
	"go/ast"
	"go/types"
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

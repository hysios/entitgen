package main

import (
	"strconv"
	"strings"

	"github.com/akrennmair/slice"
	"github.com/hysios/entitgen/gen"
)

type GenMethodContext struct {
	gen.Named

	Name    string
	VarName string
	Rece    *gen.Var
	Outputs []*gen.Var
	Inputs  []*gen.Var
	Fields  []*gen.Field
	types   *GenType
}

type Nameds map[string]bool

// ShortName returns the short name of the type.
// Example: User
func (ctx *GenTypeContext) ShortName() string {
	short := func() string {
		return strings.ToLower(string([]rune(ctx.Name)[0]))
	}
	return ctx.suggestName(short())
}

// suggestName returns a name that is not already in use.
func (ctx *GenTypeContext) suggestName(name string) string {
	if ctx.nameds[name] {
		return ctx.suggestName(name + strconv.Itoa(len(ctx.nameds)))
	}

	ctx.nameds[name] = true
	return name
}

// Bind
func (ctx *GenMethodContext) Bind() {
	ctx.Fields = bindFields(ctx.Fields, &ctx.types.TypeConverts)
}

// Receive
func (m *GenMethodContext) Receive() string {
	return m.ShortName() + " " + m.Rece.PureType()
}

func (m *GenMethodContext) ShortName() string {
	short := func() string {
		return strings.ToLower(string([]rune(m.Rece.Name)[0]))
	}
	return m.SuggestName("ShortName", short())
}

// Return
func (m *GenMethodContext) Return() string {
	var s = strings.Join(slice.Map(m.Outputs, func(arg *gen.Var) string {
		return arg.String()
	}), ",")

	if len(m.Outputs) > 1 {
		return "(" + s + ")"
	} else {
		return m.Outputs[0].PureType()
	}
}

// OutputType
func (m *GenMethodContext) OutputType() string {
	return m.Outputs[0].Address()
}

// InputArgs
func (m *GenMethodContext) InputArgs() string {
	return strings.Join(slice.Map(m.Inputs, func(arg *gen.Var) string {
		return arg.String()
	}), ",")
}

// InputVals
func (m *GenMethodContext) InputVals() string {
	return strings.Join(slice.Map(m.Inputs, func(arg *gen.Var) string {
		return arg.Name
	}), ",")
}

// InputName
func (m *GenMethodContext) InputName() string {
	return m.Inputs[0].Name
}

// UpdateChain
func (m *GenMethodContext) UpdateChain() string {
	return "updates ...UpdateChain"
}

func (ctx *GenMethodContext) FieldToProto(field *gen.Field) string {
	return field.ConvertPbType(ctx.ShortName() + "." + field.Name)
}

func (ctx *GenMethodContext) FieldToModel(field *gen.Field) string {
	return field.ConvertType(ctx.InputName() + "." + field.PbName)
}

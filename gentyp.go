package main

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/hysios/entitgen/gen"
)

type GenType struct {
	gen.Global
	gen.TypeConverts

	Name      string
	ModelName string
	PkgName   string
	PbPkgName string
	Fields    []*gen.Field
	Options   Options

	models map[string]*typeInfo
}

type TypeConvert struct {
}

type GenTypeContext struct {
	*GenType
	nameds Nameds
}

func (t *GenType) AddField(f *gen.Field) {
	t.Fields = append(t.Fields, f)
}

// Structs returns the struct types.
// func (t *GenType) Structs() []*gen.Struct {
// 	var structs []*gen.Struct
// 	for _, f := range t.Fields {
// 		if s, ok := f.Type.(*gen.Struct); ok {
// 			structs = append(structs, s)
// 		}
// 	}
// }

// ToProtoMethod returns the protobuf type name for the given type.
// Example: ToProto() *pb.Type
//
//	func (u *User) ToProto() *pb.User {
//		return &pb.User{
//			Id:        u.Id,
//			FirstName: u.FirstName,
//			LastName:  u.LastName,
//			Email:     u.Email,
//			Password:  null.NullToString(u.Password),
//			Phone:     u.Phone,
//			Role:      u.Role,
//			Active:    u.Active,
//			CreatedAt: timestamppb.New(u.CreatedAt),
//			UpdatedAt: timestamppb.New(u.UpdatedAt),
//			Member:   u.Member.ToProto(),
//		}
//	}
func (t *GenType) ToProtoMethod() string {
	var (
		tmpl = template.Must(template.New("toProto").Parse(toProtoTemplate))
		buf  = &bytes.Buffer{}
		ctx  = &GenMethodContext{
			Rece: &gen.Var{
				Name:    t.Name,
				Pointer: true,
				Type:    t.ModelName,
			},
			Fields: t.Fields,
			types:  t,
			Outputs: []*gen.Var{
				{
					Name:    "pb" + t.Name,
					Pointer: true,
					Type:    t.PbPkgName + "." + t.Name,
				},
			},
		}
	)

	ctx.Bind()

	if err := tmpl.Execute(buf, ctx); err != nil {
		panic(err)
	}

	return buf.String()
}

// bindFileds
func bindFields(fields []*gen.Field, types *gen.TypeConverts) []*gen.Field {
	for _, f := range fields {
		f.Bind(types)
	}
	return fields
}

// FromProtoMethod returns the protobuf type name for the given type.
// Example: FromProto(*pb.Type) *Type
//
//	func (u *User) FromProto(pbUser *pb.User) *User {
//		return &User{
//			Id:        pbUser.Id,
//			FirstName: pbUser.FirstName,
//			LastName:  pbUser.LastName,
//			Email:     pbUser.Email,
//			Password:  null.ToNullString(pbUser.Password),
//			Phone:     pbUser.Phone,
//			Role:      pbUser.Role,
//			Active:    pbUser.Active,
//			Member:    pbUser.Member.FromProto(),
//			CreatedAt: pbUser.CreatedAt.AsTime(),
//			UpdatedAt: pbUser.UpdatedAt.AsTime(),
//		}
//	}
func (t *GenType) FromProtoMethod() string {
	var (
		tmpl = template.Must(template.New("fromProto").Parse(fromProtoTemplate))
		buf  = &bytes.Buffer{}
		ctx  = &GenMethodContext{
			Rece: &gen.Var{
				Name:    t.Name,
				Pointer: true,
				Type:    t.ModelName,
			},
			Fields: t.Fields,
			types:  t,
			Outputs: []*gen.Var{
				{
					Name:    strings.ToLower(t.Name),
					Pointer: true,
					Type:    t.ModelName,
				},
			},
			Inputs: []*gen.Var{
				&gen.Var{
					Name:    "p" + t.Name,
					Pointer: true,
					Type:    t.PbPkgName + "." + t.Name,
				},
			},
		}
	)

	if err := tmpl.Execute(buf, ctx); err != nil {
		panic(err)
	}

	return buf.String()
}

// ModelFromProtoMethod
func (t *GenType) ModelFromProtoMethod() string {
	var (
		tmpl = template.Must(template.New("modelFromProto").Parse(modelFromProtoTemplate))
		buf  = &bytes.Buffer{}
		ctx  = &GenMethodContext{
			Name:      t.Name,
			ModelName: t.ModelName,
			Rece: &gen.Var{
				Name:    t.Name,
				Pointer: true,
				Type:    t.ModelName,
			},
			Fields: t.Fields,
			types:  t,
			Outputs: []*gen.Var{
				{
					Name:    strings.ToLower(t.Name),
					Pointer: true,
					Type:    t.ModelName,
				},
			},
			Inputs: []*gen.Var{
				&gen.Var{
					Name:    "p" + t.Name,
					Pointer: true,
					Type:    t.PbPkgName + "." + t.Name,
				},
			},
		}
	)

	if err := tmpl.Execute(buf, ctx); err != nil {
		panic(err)
	}

	return buf.String()
}

// ModelToProtoMethod
func (t *GenType) ModelToProtoMethod() string {
	var (
		tmpl = template.Must(template.New("modelToProtoMethod").Parse(modelToProtoTemplate))
		buf  = &bytes.Buffer{}
		ctx  = &GenMethodContext{
			Name:      t.Name,
			ModelName: t.ModelName,
			Rece: &gen.Var{
				Name:    t.Name,
				Pointer: true,
				Type:    t.Name,
			},
			VarName: t.suggestName(strings.ToLower(t.Name)),
			Fields:  t.Fields,
			types:   t,
			Outputs: []*gen.Var{
				{
					Name:    strings.ToLower(t.Name),
					Pointer: true,
					Type:    t.PbPkgName + "." + t.Name,
				},
			},
			Inputs: []*gen.Var{
				&gen.Var{
					Name:    t.suggestName(strings.ToLower(t.Name)),
					Pointer: true,
					Type:    t.ModelName,
				},
			},
		}
	)

	if err := tmpl.Execute(buf, ctx); err != nil {
		panic(err)
	}

	return buf.String()
}

var reversedNames = []string{
	"import", "for", "if",
}

// suggestName
func (t *GenType) suggestName(name string) string {
	for _, n := range reversedNames {
		if n == name {
			return name + "x"
		}
	}
	return name
}

// ModelListFromMethod
func (t *GenType) ModelListFromMethod() string {
	var (
		tmpl = template.Must(template.New("modelListFromMethod").Parse(modelListFromTemplate))
		buf  = &bytes.Buffer{}
		ctx  = &GenMethodContext{
			Name:      t.Name,
			ModelName: t.ModelName,
			Rece: &gen.Var{
				Name:    t.Name,
				Pointer: true,
				Type:    t.Name,
			},
			VarName: strings.ToLower(plural(t.Name)),
			Fields:  t.Fields,
			types:   t,
			Outputs: []*gen.Var{
				{
					Name:    strings.ToLower(t.Name),
					Pointer: true,
					Type:    t.ModelName,
					Slice:   []string{""},
				},
			},
			Inputs: []*gen.Var{
				&gen.Var{
					Name:    strings.ToLower(plural(t.Name)),
					Pointer: true,
					Type:    t.PbPkgName + "." + t.Name,
					Slice:   []string{""},
				},
			},
		}
	)

	if err := tmpl.Execute(buf, ctx); err != nil {
		panic(err)
	}

	return buf.String()
}

// ModelListToMethod
func (t *GenType) ModelListToMethod() string {
	var (
		tmpl = template.Must(template.New("modelListToMethod").Parse(modelListToTemplate))
		buf  = &bytes.Buffer{}
		ctx  = &GenMethodContext{
			Name:      t.Name,
			ModelName: t.ModelName,
			Rece: &gen.Var{
				Name:    t.Name,
				Pointer: true,
				Type:    t.Name,
			},
			VarName: strings.ToLower(plural(t.Name)),
			Fields:  t.Fields,
			types:   t,
			Outputs: []*gen.Var{
				{
					Name:    strings.ToLower(t.Name),
					Pointer: true,
					Type:    t.PbPkgName + "." + t.Name,
					Slice:   []string{""},
				},
			},
			Inputs: []*gen.Var{
				&gen.Var{
					Name:    strings.ToLower(plural(t.Name)),
					Pointer: true,
					Type:    t.ModelName,
					Slice:   []string{""},
				},
			},
		}
	)

	if err := tmpl.Execute(buf, ctx); err != nil {
		panic(err)
	}

	return buf.String()
}

// NoModel
func (m *GenType) NoModel() bool {
	for _, f := range m.Options.NoModels {
		if m.Name == f {
			return true
		}
	}

	if _, ok := m.models[m.ModelName]; ok {
		return true
	}

	return false
}

// GenSlice
func (m *GenType) GenSlice() bool {
	return m.Options.GenSlice
}

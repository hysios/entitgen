package gen

type Field struct {
	ID       string
	Name     string
	PbName   string
	Pkg      string
	PbType   string
	Type     string
	Pointer  bool
	Children []*Field
	Tags     []string
	Comment  string

	types *TypeConverts
}

// DeclareField
func (v *Field) DeclareField() string {
	if v.Pointer {
		return v.Name + " *" + v.Type
	}
	return v.Name + " " + v.Type
}

// Bind
func (v *Field) Bind(types *TypeConverts) {
	v.types = types
}

// ConvertType()
func (v *Field) ConvertType(in string) string {
	if v.types == nil {
		return in
	}

	conv := v.types.GetConv(v.ID + "_" + v.Type)
	if conv == nil {
		return in
	}

	return conv.To(in)
}

// GetPbType()
func (v *Field) ConvertPbType(in string) string {
	if v.types == nil {
		return in
	}

	conv := v.types.GetConv(v.ID + "_" + v.Type)
	if conv == nil {
		return in
	}

	return conv.From(in)
}

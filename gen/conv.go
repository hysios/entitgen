package gen

type TypeConverts struct {
	converts map[string]Converter
}

type Converter interface {
	From(in string) string
	To(in string) string
}

// AddConv
func (v *TypeConverts) AddConv(name string, conv Converter) {
	if v.converts == nil {
		v.converts = make(map[string]Converter)
	}
	v.converts[name] = conv
}

// GetConv
func (v *TypeConverts) GetConv(name string) Converter {
	if v.converts == nil {
		return nil
	}
	return v.converts[name]
}

// // ConvFromType
// func (v *TypeConverts) ConvFromType(fromType string, in string) string {
// 	if v.converts == nil {
// 		return in
// 	}
// 	return v.converts[fromType].From(fromType)
// }

// // ConvToType
// func (v *TypeConverts) ConvToType(toType string, in string) string {
// 	if v.converts == nil {
// 		return in
// 	}

// 	return v.converts[toType].To(toType)
// }

func TypeConv(from, to string) Converter {
	return &simpleConv{from, to}
}

type simpleConv struct {
	from, to string
}

func (v *simpleConv) From(in string) string {
	return v.from + "(" + in + ")"
}

func (v *simpleConv) To(in string) string {
	return v.to + "(" + in + ")"
}

// ProtoConv
func ProtoConv(model string) Converter {
	return &protoConv{model}
}

type protoConv struct {
	model string
}

func (v *protoConv) From(in string) string {
	return in + ".ToProto()"
}

func (v *protoConv) To(in string) string {
	return "(*" + v.model + ")(nil).FromProto(" + in + ")"
}

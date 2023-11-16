package gen

type Var struct {
	Name    string
	PkgName string
	Pointer bool
	Type    string
	Slice   []string
}

// String
func (v *Var) String() string {
	var s string

	if v.Name != "" {
		s = v.Name + " "
	}

	return s + v.PureType()
}

// PureType
func (v *Var) PureType() string {
	var s string

	if v.isSlice() {
		return v.SliceType()
	} else if v.Pointer {
		s += "*" + v.TypeString()
	} else {
		s += v.TypeString()
	}

	return s
}

// isSlice
func (v *Var) isSlice() bool {
	return len(v.Slice) > 0
}

// SliceType
func (v *Var) SliceType() string {
	var s string

	if len(v.Slice) > 0 {
		for _, v := range v.Slice {
			s += "["
			s += v
			s += "]"
		}
	}
	if v.Pointer {
		s += "*"
	}
	s += v.TypeString()
	return s
}

// TypeString
func (v *Var) TypeString() string {
	var s string

	if v.PkgName != "" {
		s = v.PkgName + "."
	}

	s += v.Type
	return s
}

// Address
func (v *Var) Address() string {
	return "&" + v.TypeString()
}

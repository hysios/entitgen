package gen

type Var struct {
	Name    string
	PkgName string
	Pointer bool
	Type    string
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

	if v.Pointer {
		s += "*" + v.TypeString()
	} else {
		s += v.TypeString()
	}

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

package gen

import "strconv"

type Pkg struct {
	Fullname string
	Alias    string
}

// String pkg name stringer
func (p Pkg) String() string {
	if p.Alias != "" {
		return p.Alias + " " + strconv.Quote(p.Fullname)
	}
	return strconv.Quote(p.Fullname)
}

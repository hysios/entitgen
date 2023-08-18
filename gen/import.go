package gen

type Imports struct {
	pkgs []Pkg
}

func (i *Imports) AddImport(pkg Pkg) {
	// check if already added
	for _, p := range i.pkgs {
		if p.Fullname == pkg.Fullname && p.Alias == pkg.Alias {
			return
		}
	}

	i.pkgs = append(i.pkgs, pkg)
}

// Imports
func (i *Imports) GoImports() string {
	if len(i.pkgs) == 0 {
		return ""
	}

	var s string
	for _, pkg := range i.pkgs {
		s += pkg.String() + "\n"
	}

	return "import (\n" + s + ")\n"
}

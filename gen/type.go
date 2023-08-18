package gen

type Type struct {
	*Global

	Name    string
	PkgName string
	Fields  []*Field
}

package out

//go:generate entitgen -type Friend
type Friend struct {
	ID       uint
	Name     string
	Username string
	Nickname string
}

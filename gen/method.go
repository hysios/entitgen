package gen

type Method struct {
	Name    string
	Rece    Var
	Args    []*Var
	Results []*Var
}

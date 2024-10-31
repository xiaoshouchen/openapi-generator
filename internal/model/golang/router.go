package golang

type RouterItem struct {
	Method    string
	Path      string
	ShortPath string
	FuncName  string
}

type Router struct {
	Items   []RouterItem
	Imports []string
}

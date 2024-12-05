package typescript

type ApiStruct struct {
	Imports   []string
	Functions []Function
}

type Function struct {
	Path   string
	Method string
	Name   string
}

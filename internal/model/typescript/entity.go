package typescript

type EntityStruct struct {
	Name string
	Rows []EntityRow
}

type EntityRow struct {
	Name       string
	DataType   string
	Annotation string
}

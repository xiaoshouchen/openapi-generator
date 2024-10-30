package golang

type ResponseStruct struct {
	Name string
	Rows []ResponseRow
}

type ResponseRow struct {
	DataType    string
	Name        string
	Description string
}

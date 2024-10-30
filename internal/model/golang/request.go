package golang

type RequestStruct struct {
	Name string
	Rows []RequestRow
}

type RequestRow struct {
	DataType    string
	Name        string
	Validate    string
	Description string
	BindType    string
}

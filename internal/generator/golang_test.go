package generator

import (
	goModel "github.com/xiaoshouchen/openapi-generator/internal/model/golang"
	"html/template"
	"testing"
)

func TestGoGenerator_Request(t *testing.T) {
	type args struct {
		path string
		file string
		f    template.FuncMap
		data map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test",
			args{
				"./",
				"gogo.go",
				nil,
				map[string]interface{}{
					"packageName": "test",
					"structs": []goModel.RequestStruct{{
						"Test",
						[]goModel.RequestRow{
							{
								"int",
								"age",
								"required",
								"fsdfsdfa",
							},
							{
								"int64",
								"age1",
								"required",
								"测试",
							},
						},
					}},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GoGenerator{}
			if err := g.Request(tt.args.path, tt.args.file, tt.args.f, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Response() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

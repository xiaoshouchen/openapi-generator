package generator

import (
	"testing"
)

func Test_write(t *testing.T) {
	type args struct {
		source    []byte
		path      string
		file      string
		overwrite bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"测试写功能",
			args{
				[]byte("templates/go/controller.go.tmpl"),
				"./",
				"test.txt",
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			write(tt.args.source, tt.args.path, tt.args.file, tt.args.overwrite)
		})
	}
}

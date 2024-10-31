package pkg

import "testing"

func TestFormatDescription(t *testing.T) {
	type args struct {
		desc string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			"test", args{desc: "binding:\"required,test\";这是普通的字符串"},
			"required,test",
			"这是普通的字符串",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FormatDescription(tt.args.desc)
			if got != tt.want {
				t.Errorf("FormatDescription() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FormatDescription() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetTopLevelName(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"api/user/login",
			args{path: "api/user/login"},
			"api",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTopLevelName(tt.args.path); got != tt.want {
				t.Errorf("GetTopLevelName() = %v, want %v", got, tt.want)
			}
		})
	}
}

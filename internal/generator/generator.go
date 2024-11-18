package generator

import (
	"bytes"
	"github.com/xiaoshouchen/openapi-generator/internal/model"
	"github.com/xiaoshouchen/openapi-generator/pkg"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Generator interface {
	Request(path, file string, f template.FuncMap, data map[string]interface{}) error
	Response(path, file string, f template.FuncMap, data map[string]interface{}) error
	Router(path, file string, f template.FuncMap, data map[string]interface{}) error
	Controller(path, file string, f template.FuncMap, data map[string]interface{}) error
	Service(path, file string, f template.FuncMap, data map[string]interface{}) error
}

func NewGenerator(conf model.Config) Generator {
	var generator Generator
	switch conf.AimType {
	case "go":
		generator = NewGoGenerator()
	default:
		log.Fatal("不支持的生成类型")
	}
	return generator
}

func genTpl(tpl string, f template.FuncMap, data map[string]interface{}) []byte {
	t, err := template.New("tplFile").Funcs(f).Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}
	var buf = new(bytes.Buffer)
	err = t.ExecuteTemplate(buf, "tplFile", data)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func write(source []byte, path, file string, overwrite bool) {

	absPath := strings.TrimRight(path, "/") + "/" + file

	dir := filepath.Dir(absPath)

	// 创建多层目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("创建目录失败 %s: %v", dir, err)
	}

	if !overwrite {
		exists, err := pkg.FileExists(absPath)
		if err != nil {
			log.Println(err)
		}
		if exists {
			return
		}
	}
	err := os.WriteFile(absPath, source, 0664)
	if err != nil {
		log.Println(err)
	}
}

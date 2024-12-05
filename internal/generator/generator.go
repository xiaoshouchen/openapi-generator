package generator

import (
	"bytes"
	"github.com/xiaoshouchen/openapi-generator/internal/model"
	"github.com/xiaoshouchen/openapi-generator/pkg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type genConfig struct {
	path      string // 文件路径
	file      string //文件名称
	overwrite bool   // 是否可以覆盖
	tpl       string // 模板地址
}

type Generator struct {
	generator I
}

// Generate 生成代码文件
// genType 是生成的类型，比如 request, response
// path 是生成的文件的相对路径,比如service/controller/api/user/list.go
func (g *Generator) Generate(genType string, path string, f template.FuncMap, data map[string]interface{}) error {
	config := g.generator.parseGenType(genType, path)
	source := genTpl(config.tpl, f, data)
	source = g.generator.formatSource(source)
	write(source, config.path, config.file, config.overwrite)
	return nil
}

type I interface {
	parseGenType(genType string, path string) genConfig
	formatSource(source []byte) []byte
}

func NewGenerator(conf model.Config) Generator {
	var generator I
	switch conf.AimType {
	case "go":
		generator = NewGoGenerator(conf)
	case "ts":
		generator = NewTSGenerator(conf)
	default:
		log.Fatal("不支持的生成类型")
	}
	return Generator{
		generator: generator,
	}
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

package process

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/xiaoshouchen/openapi-generator/internal/generator"
	"github.com/xiaoshouchen/openapi-generator/internal/model"
	"github.com/xiaoshouchen/openapi-generator/pkg"
)

// Processor Process Raw Data and produce final data
type Processor interface {
	Process(schema *model.OpenAPISchema, generator generator.Generator)
}

func NewProcessor(conf model.Config) Processor {
	var processor Processor
	switch conf.AimType {
	case "go":
		processor = NewGolang(conf)
	case "ts":
		processor = NewTypescript(conf)
	default:
		log.Fatal("不支持的类型")
	}
	return processor
}

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"lowCamel": pkg.LineToLowCamel,
		"upCamel":  pkg.LineToUpCamel,
		"inline":   pkg.Inline,
	}
}

// 支持模糊匹配删除文件
func deleteFile(path string) {
	files, err := filepath.Glob(path)
	if err != nil {
		log.Println(err)
	}
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			log.Println(err)
		}
	}
}

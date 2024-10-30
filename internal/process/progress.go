package process

import (
	"github.com/xiaoshouchen/openapi-go-generator/internal/generator"
	"github.com/xiaoshouchen/openapi-go-generator/internal/model"
	"github.com/xiaoshouchen/openapi-go-generator/pkg"
	"html/template"
	"log"
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

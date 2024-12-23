package generator

import (
	_ "embed"
	"go/format"
	"log"

	"github.com/xiaoshouchen/openapi-generator/internal/enum"
	"github.com/xiaoshouchen/openapi-generator/internal/model"
)

type GoGenerator struct {
	config model.Config
}

func NewGoGenerator(config model.Config) *GoGenerator {
	return &GoGenerator{
		config: config,
	}
}

func (g *GoGenerator) formatSource(source []byte) []byte {
	sourceFormat, err := format.Source(source)
	if err != nil {
		log.Printf("request 格式化失败,%s,%v", source, err)
		return source
	}
	return sourceFormat
}

//go:embed templates/go/request.go.tmpl
var requestTpl string

//go:embed templates/go/response.go.tmpl
var responseTpl string

//go:embed templates/go/controller.go.tmpl
var controllerTpl string

//go:embed templates/go/service.go.tmpl
var serviceTpl string

//go:embed templates/go/router.go.tmpl
var routerTpl string

func (g *GoGenerator) parseGenType(genType string, path string) genConfig {
	switch genType {
	case enum.GeneratorGoRequest:
		return genConfig{
			path:      g.config.OutPath,
			file:      path,
			overwrite: true,
			tpl:       requestTpl,
		}
	case enum.GeneratorGoResponse:
		return genConfig{
			path:      g.config.OutPath,
			file:      path,
			overwrite: true,
			tpl:       responseTpl,
		}
	case enum.GeneratorGoController:
		return genConfig{
			path:      g.config.OutPath,
			file:      path,
			overwrite: true,
			tpl:       controllerTpl,
		}
	case enum.GeneratorGoRouter:
		return genConfig{
			path:      g.config.OutPath,
			file:      path,
			overwrite: true,
			tpl:       routerTpl,
		}
	case enum.GeneratorGoService:
		return genConfig{
			path:      g.config.OutPath,
			file:      path,
			overwrite: false,
			tpl:       serviceTpl,
		}
	default:
		return genConfig{}
	}
}

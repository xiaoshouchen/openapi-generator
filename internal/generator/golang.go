package generator

import (
	_ "embed"
	"go/format"
	"log"
	"os"
	"path/filepath"

	"github.com/xiaoshouchen/openapi-generator/internal/enum"
	"github.com/xiaoshouchen/openapi-generator/internal/model"
	"github.com/xiaoshouchen/openapi-generator/pkg"
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

//go:embed templates/go/response_func.go.tmpl
var responseFuncTpl string

func getTpl(embeddedTpl, tplPath string) string {
	localTplPath := filepath.Join(".templates", tplPath)
	exists, err := pkg.FileExists(localTplPath)
	if err != nil || !exists {
		return embeddedTpl
	}
	content, err := os.ReadFile(localTplPath)
	if err != nil {
		return embeddedTpl
	}
	return string(content)
}

func (g *GoGenerator) parseGenType(genType string, path string) genConfig {
	switch genType {
	case enum.GeneratorGoRequest:
		return genConfig{
			path:      g.config.OutPath,
			file:      path,
			overwrite: true,
			tpl:       getTpl(requestTpl, "go/request.go.tmpl"),
		}
	case enum.GeneratorGoResponse:
		return genConfig{
			path:      g.config.OutPath,
			file:      path,
			overwrite: true,
			tpl:       getTpl(responseTpl, "go/response.go.tmpl"),
		}
	case enum.GeneratorGoController:
		return genConfig{
			path:      g.config.OutPath,
			file:      path,
			overwrite: true,
			tpl:       getTpl(controllerTpl, "go/controller.go.tmpl"),
		}
	case enum.GeneratorGoRouter:
		return genConfig{
			path:      g.config.OutPath,
			file:      path,
			overwrite: true,
			tpl:       getTpl(routerTpl, "go/router.go.tmpl"),
		}
	case enum.GeneratorGoService:
		return genConfig{
			path:      g.config.OutPath,
			file:      path,
			overwrite: false,
			tpl:       getTpl(serviceTpl, "go/service.go.tmpl"),
		}
	case enum.GeneratorGoResponseFunc:
		return genConfig{
			path:      g.config.OutPath,
			file:      path,
			overwrite: false,
			tpl:       getTpl(responseFuncTpl, "go/response_func.go.tmpl"),
		}
	default:
		return genConfig{}
	}
}

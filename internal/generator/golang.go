package generator

import (
	_ "embed"
	"go/format"
	"html/template"
	"log"
)

type GoGenerator struct {
}

func NewGoGenerator() *GoGenerator {
	return &GoGenerator{}
}

//go:embed templates/go/request.go.tmpl
var requestTpl string

func (g *GoGenerator) Request(path, file string, f template.FuncMap, data map[string]interface{}) error {
	bytes := genTpl(requestTpl, f, data)
	source, err := format.Source(bytes)
	if err != nil {
		source = bytes
		log.Printf("request 格式化失败,%s,%v", file, err)
	}
	write(source, path, file, true)
	return nil
}

//go:embed templates/go/response.go.tmpl
var responseTpl string

func (g *GoGenerator) Response(path, file string, f template.FuncMap, data map[string]interface{}) error {
	bytes := genTpl(responseTpl, f, data)
	source, err := format.Source(bytes)
	if err != nil {
		source = bytes
		log.Printf("response 格式化失败,%s,%v", file, err)
	}
	write(source, path, file, true)
	return nil
}

//go:embed templates/go/controller.go.tmpl
var controllerTpl string

func (g *GoGenerator) Controller(path, file string, f template.FuncMap, data map[string]interface{}) error {
	bytes := genTpl(controllerTpl, f, data)
	source, err := format.Source(bytes)
	if err != nil {
		source = bytes
		log.Printf("controller 格式化失败,%s,%v", file, err)
	}
	write(source, path, file, true)
	return nil
}

//go:embed templates/go/service.go.tmpl
var serviceTpl string

func (g *GoGenerator) Service(path, file string, f template.FuncMap, data map[string]interface{}) error {
	bytes := genTpl(serviceTpl, f, data)
	source, err := format.Source(bytes)
	if err != nil {
		source = bytes
		log.Printf("service 格式化失败,%s,%v", file, err)
	}
	write(source, path, file, false)
	return nil
}

//go:embed templates/go/router.go.tmpl
var routerTpl string

func (g *GoGenerator) Router(path, file string, f template.FuncMap, data map[string]interface{}) error {
	bytes := genTpl(routerTpl, f, data)
	source, err := format.Source(bytes)
	if err != nil {
		source = bytes
		log.Printf("router 格式化失败,%s,%v", file, err)
	}
	write(source, path, file, true)
	return nil
}

package generator

import (
	_ "embed"

	"github.com/xiaoshouchen/openapi-generator/internal/enum"
	"github.com/xiaoshouchen/openapi-generator/internal/model"
)

type TsGenerator struct {
	config model.Config
}

func NewTSGenerator(config model.Config) *TsGenerator {
	return &TsGenerator{
		config: config,
	}
}

//go:embed templates/ts/api.ts.tmpl
var tsApiTpl string

//go:embed templates/ts/entity.ts.tmpl
var tsEntityTpl string

func (t *TsGenerator) parseGenType(genType string, path string) genConfig {
	switch genType {
	case enum.GeneratorTsApi:
		return genConfig{
			path:      t.config.OutPath,
			file:      path,
			overwrite: true,
			tpl:       getTpl(tsApiTpl, "ts/api.ts.tmpl"),
		}
	case enum.GeneratorTsEntity:
		return genConfig{
			path:      t.config.OutPath,
			file:      path,
			overwrite: true,
			tpl:       getTpl(tsEntityTpl, "ts/entity.ts.tmpl"),
		}
	}
	return genConfig{}
}

func (t *TsGenerator) formatSource(source []byte) []byte {
	return source
}

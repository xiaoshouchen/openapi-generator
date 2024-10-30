package parser

import (
	"encoding/json"
	"github.com/xiaoshouchen/openapi-go-generator/internal/model"
	"gopkg.in/yaml.v3"
)

type Parser struct {
	data   []byte
	format string
}

func NewParser(data []byte, format string) *Parser {
	parse := &Parser{
		data:   data,
		format: format,
	}
	return parse
}

func (p *Parser) Parse() (*model.OpenAPISchema, error) {
	var schema = new(model.OpenAPISchema)
	if p.format == "json" {
		err := json.Unmarshal(p.data, schema)
		if err != nil {
			return nil, err
		}
	}
	if p.format == "yaml" {
		err := yaml.Unmarshal(p.data, schema)
		if err != nil {
			return nil, err
		}
	}
	return schema, nil
}

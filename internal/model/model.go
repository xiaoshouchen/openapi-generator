package model

import (
	"encoding/json"
	"fmt"
)

// OpenAPISchema represents the root of the OpenAPI document
type OpenAPISchema struct {
	OpenAPI    string              `json:"openapi"`
	Info       Info                `json:"info"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
}

// Info provides metadata about the API
type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

// PathItem describes the operations available on a single path
type PathItem struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
	// Add other HTTP methods as needed
}

// Operation describes a single API operation on a path
type Operation struct {
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	OperationID string              `json:"operationId,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses"`
}

// Parameter describes a single operation parameter
type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"` // query, header, path or cookie
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

// RequestBody describes a single request body
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content"`
	Required    bool                 `json:"required,omitempty"`
}

// Response describes a single response from an API Operation
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

// MediaType provides schema and examples for the media type identified by its key
type MediaType struct {
	Schema SchemaOrArray `json:"schema,omitempty"`
}

type SchemaOrArray struct {
	Schema  *Schema
	Schemas []Schema
}

func (soa *SchemaOrArray) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a single Schema
	var schema Schema
	if err := json.Unmarshal(data, &schema); err == nil {
		soa.Schema = &schema
		return nil
	}

	// If it's not a single Schema, try to unmarshal as an array of Schemas
	var schemas []Schema
	if err := json.Unmarshal(data, &schemas); err == nil {
		soa.Schemas = schemas
		return nil
	}

	return fmt.Errorf("schema must be either an object or an array")
}

// Components holds a set of reusable objects for different aspects of the OAS
type Components struct {
	Schemas map[string]Schema `json:"schemas,omitempty"`
}

// Schema describes the structure of a type in OpenAPI
type Schema struct {
	Type                 string           `json:"type,omitempty"`
	Format               string           `json:"format,omitempty"`
	Properties           SchemaProperties `json:"properties,omitempty"`
	Items                *SchemaOrArray   `json:"items,omitempty"`
	Ref                  string           `json:"$ref,omitempty"`
	Required             []string         `json:"required,omitempty"`
	Description          string           `json:"description,omitempty"`
	AdditionalProperties *SchemaOrArray   `json:"additionalProperties,omitempty"`
}

type SchemaProperties map[string]Schema

// UnmarshalJSON custom unmarshaler for SchemaProperties
func (sp *SchemaProperties) UnmarshalJSON(data []byte) error {
	// First, try to unmarshal as an object
	var objMap map[string]json.RawMessage
	err := json.Unmarshal(data, &objMap)
	if err == nil {
		*sp = make(SchemaProperties)
		for k, v := range objMap {
			var schema Schema
			if err := json.Unmarshal(v, &schema); err != nil {
				return err
			}
			(*sp)[k] = schema
		}
		return nil
	}

	// If it's not an object, try to unmarshal as an array
	var arrayProps []json.RawMessage
	err = json.Unmarshal(data, &arrayProps)
	if err == nil {
		*sp = make(SchemaProperties)
		for _, v := range arrayProps {
			var propObj map[string]Schema
			if err := json.Unmarshal(v, &propObj); err != nil {
				return err
			}
			for k, schema := range propObj {
				(*sp)[k] = schema
			}
		}
		return nil
	}

	return fmt.Errorf("properties must be either an object or an array")
}

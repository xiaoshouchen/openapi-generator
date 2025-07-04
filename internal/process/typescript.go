package process

import (
	"log"
	"sort"
	"strings"

	"github.com/xiaoshouchen/openapi-generator/internal/enum"
	"github.com/xiaoshouchen/openapi-generator/internal/generator"
	"github.com/xiaoshouchen/openapi-generator/internal/model"
	tsModel "github.com/xiaoshouchen/openapi-generator/internal/model/typescript"
	"github.com/xiaoshouchen/openapi-generator/pkg"
)

type Typescript struct {
	config model.Config
}

func NewTypescript(config model.Config) *Typescript {
	return &Typescript{
		config: config,
	}
}

func (t *Typescript) GoTypeMap(typeName string) string {
	var m = map[string]string{
		"integer": "number",
		"number":  "number",
		"boolean": "boolean",
		"string":  "string",
	}
	if str, ok := m[typeName]; ok {
		return str
	}
	return typeName
}

func (t *Typescript) Process(schema *model.OpenAPISchema, generator generator.Generator) {
	var apiMap = make(map[string]tsModel.ApiStruct)
	var entityMap = make(map[string][]tsModel.EntityStruct)
	for path, item := range schema.Paths {
		prefix := strings.Split(strings.Trim(path, "/"), "/")[0]
		if pkg.ArrayContains(t.config.OmitPrefixPath, prefix) {
			continue
		}
		uniKey := t.getApiUniqueKey(path)
		funcName := t.getFuncName(path)
		apiFunc := tsModel.Function{
			Path:   path,
			Method: "get",
			Name:   funcName,
		}

		if item.Get != nil {
			entityMap[uniKey] = append(entityMap[uniKey], t.ProcessGetRequest(t.getEntityName(path), item.Get.Parameters))
			for k, v := range item.Get.Responses {
				if k == "200" {
					if jsonData, ok := v.Content["application/json"]; ok {
						if sch := jsonData.Schema.Schema; sch != nil {
							entityMap[uniKey] = append(entityMap[uniKey], t.processResponse(t.getEntityRespName(path), sch.Properties, []tsModel.EntityStruct{})...)
						}
					}
				}
			}
		}
		if item.Post != nil {
			apiFunc.Method = "post"
			// Request 逻辑
			if item.Post.RequestBody == nil {
				log.Println("no request body", path)
				continue
			}
			if item.Post.RequestBody.Content == nil {
				log.Println("no request body", path)
				continue
			}
			sch := item.Post.RequestBody.Content["application/json"].Schema
			if schSch := sch.Schema; schSch != nil {
				entityMap[uniKey] = append(entityMap[uniKey], t.processPostRequest(t.getEntityName(path), schSch.Properties, []tsModel.EntityStruct{}, schSch.Required)...)
			}
			for k, v := range item.Post.Responses {
				if k == "200" {
					if jsonData, ok := v.Content["application/json"]; ok {
						if sch := jsonData.Schema.Schema; sch != nil {
							entityMap[uniKey] = append(entityMap[uniKey], t.processResponse(t.getEntityRespName(path), sch.Properties, []tsModel.EntityStruct{})...)
						}
					}
				}
			}
		}
		f := apiMap[uniKey]
		f.Functions = append(f.Functions, apiFunc)
		f.Imports = append(f.Imports, pkg.LineToUpCamel(funcName)+"Req", pkg.LineToUpCamel(funcName)+"Resp")
		apiMap[uniKey] = f
	}

	for k, api := range apiMap {
		functions := api.Functions
		// 排序
		sort.Slice(functions, func(i, j int) bool {
			// 根据name进行排序
			return functions[i].Name < functions[j].Name
		})
		imports := api.Imports
		sort.Slice(imports, func(i, j int) bool {
			// 根据name进行排序
			return imports[i] < imports[j]
		})
		_ = generator.Generate(enum.GeneratorTsApi, pkg.LineToLowCamel("api/modules/"+k+".ts"), FuncMap(), map[string]interface{}{
			"functions":  functions,
			"imports":    imports,
			"importPath": pkg.LineToLowCamel(k),
		})
	}

	for k, entities := range entityMap {
		// 排序
		for index, entity := range entities {
			sort.Slice(entity.Rows, func(i, j int) bool {
				// 根据name进行排序
				return entity.Rows[i].Name < entity.Rows[j].Name
			})
			entities[index] = entity
		}
		sort.Slice(entities, func(i, j int) bool {
			// 根据name进行排序
			return entities[i].Name < entities[j].Name
		})
		_ = generator.Generate(enum.GeneratorTsEntity, pkg.LineToLowCamel("types/"+k+".ts"), FuncMap(), map[string]interface{}{
			"entities": entities,
		})
	}
}

func (t *Typescript) getApiUniqueKey(path string) string {
	params := strings.Split(strings.Trim(path, "/"), "/")
	if len(params) == 3 {
		return params[0] + "/" + params[1]
	}
	log.Fatal("api path error", path)
	return path
}

func (t *Typescript) getFuncName(path string) string {
	params := strings.Split(strings.Trim(path, "/"), "/")
	if len(params) == 3 {
		return params[1] + "_" + params[2]
	}
	log.Fatal("api path error", path)
	return path
}

func (t *Typescript) getEntityName(path string) string {
	params := strings.Split(strings.Trim(path, "/"), "/")
	if len(params) == 3 {
		return pkg.LineToUpCamel(params[1] + "_" + params[2] + "Req")
	}
	return path
}

func (t *Typescript) getEntityRespName(path string) string {
	params := strings.Split(strings.Trim(path, "/"), "/")
	if len(params) == 3 {
		return pkg.LineToUpCamel(params[1] + "_" + params[2] + "Resp")
	}
	return path
}

func (t *Typescript) ProcessGetRequest(name string, parameters []model.Parameter) tsModel.EntityStruct {
	var st tsModel.EntityStruct
	st.Name = name
	for _, param := range parameters {
		var req tsModel.EntityRow
		req.DataType = t.GoTypeMap(param.Schema.Type)
		req.Name = param.Name
		if !param.Required {
			req.Name += "?"
		}
		req.Annotation = param.Description
		st.Rows = append(st.Rows, req)
	}

	return st
}

func (t *Typescript) processPostRequest(name string, schema model.SchemaProperties, structs []tsModel.EntityStruct, required []string) []tsModel.EntityStruct {
	var st tsModel.EntityStruct
	in, err := pkg.NewQuickInArray(required)
	if err != nil {
		log.Println(err)
	}
	st.Name = name
	for k, v := range schema {
		var req tsModel.EntityRow
		req.DataType = t.GoTypeMap(v.Type)
		if req.DataType == "object" {
			req.DataType = name + pkg.LineToUpCamel(k)
			structs = t.processPostRequest(name+pkg.LineToUpCamel(k), v.Properties, structs, v.Required)
		}
		if req.DataType == "array" {
			var items *model.Schema
			req.DataType, items = t.arrayType(v.Items, name+pkg.LineToUpCamel(k), "[]")
			if items != nil {
				structs = t.processPostRequest(name+pkg.LineToUpCamel(k), items.Properties, structs, items.Required)
			}
		}
		req.Name = k
		if !in.InArray(k) {
			req.Name += "?"
		}
		req.Annotation = v.Description
		st.Rows = append(st.Rows, req)
	}
	structs = append(structs, st)
	return structs
}

func (t *Typescript) processResponse(name string, schema model.SchemaProperties, structs []tsModel.EntityStruct) []tsModel.EntityStruct {
	var st tsModel.EntityStruct
	st.Name = name
	for k, v := range schema {
		var resp tsModel.EntityRow
		resp.DataType = t.GoTypeMap(v.Type)
		if resp.DataType == "object" {
			resp.DataType = name + pkg.LineToUpCamel(k)
			structs = t.processResponse(name+pkg.LineToUpCamel(k), v.Properties, structs)
		}
		if resp.DataType == "array" {
			var items *model.Schema
			resp.DataType, items = t.arrayType(v.Items, name+pkg.LineToUpCamel(k), "[]")
			if items != nil {
				structs = t.processResponse(name+pkg.LineToUpCamel(k), items.Properties, structs)
			}
		}
		resp.Name = k
		resp.Annotation = v.Description
		st.Rows = append(st.Rows, resp)
	}
	structs = append(structs, st)
	return structs
}

func (t *Typescript) arrayType(schOrArr *model.SchemaOrArray, parentName string, suffix string) (string, *model.Schema) {
	if schOrArr == nil {
		return "any[]", nil
	}
	if schOrArr.Schema == nil {
		return "any[]", nil
	}
	if schOrArr.Schema.Type == "object" {
		return parentName + suffix, schOrArr.Schema
	}
	if schOrArr.Schema.Type == "array" {
		return t.arrayType(schOrArr.Schema.Items, parentName, suffix+"[]")
	}
	return t.GoTypeMap(schOrArr.Schema.Type) + "[]", nil
}

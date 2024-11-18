package process

import (
	"github.com/xiaoshouchen/openapi-generator/internal/generator"
	"github.com/xiaoshouchen/openapi-generator/internal/model"
	goModel "github.com/xiaoshouchen/openapi-generator/internal/model/golang"
	"github.com/xiaoshouchen/openapi-generator/pkg"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

type Golang struct {
	config model.Config
}

func NewGolang(config model.Config) *Golang {
	return &Golang{config: config}
}

// Process the json data
// TODO support formData
func (g *Golang) Process(schema *model.OpenAPISchema, generator generator.Generator) {
	var routers = make(map[string]*goModel.Router)
	for path, item := range schema.Paths {
		var router goModel.RouterItem
		var method string
		// packageName
		packName := pkg.GetPackageName(path)
		reqImportPath := g.config.ProjectName + "/" + g.config.OutPath + "/" + "server/request" + filepath.Dir(path+".go")
		respImportPath := g.config.ProjectName + "/" + g.config.OutPath + "/" + "server/response" + filepath.Dir(path+".go")
		svcImportPath := g.config.ProjectName + "/" + g.config.OutPath + "/" + "service" + filepath.Dir(path+".go")
		controllerImportPath := g.config.ProjectName + "/" + g.config.OutPath + "/" + "server/controller" + filepath.Dir(path+".go")
		reqShortPath := packName + "_request"
		respShortPath := packName + "_response"
		svcShortPath := packName + "_service"

		var getRequest goModel.RequestStruct
		var postRequests []goModel.RequestStruct
		var structs []goModel.RequestStruct

		var responseStructs []goModel.ResponseStruct
		if item.Post != nil {
			method = "POST"
			// Request 逻辑
			sch := item.Post.RequestBody.Content["application/json"].Schema
			if schSch := sch.Schema; schSch != nil {
				postRequests = g.processPostRequest("json", pkg.GetRequestName(path), schSch.Properties, []goModel.RequestStruct{}, schSch.Required)
			}
			if len(item.Post.Parameters) > 0 {
				temp := g.ProcessGetRequest(pkg.GetRequestName(path), item.Post.Parameters)
				for j, post := range postRequests {
					if getRequest.Name == post.Name {
						postRequests[j].Rows = append(postRequests[j].Rows, temp.Rows...)
					}
				}
			}
			structs = postRequests

			// Response逻辑
			for k, v := range item.Post.Responses {
				if k == "200" {
					if jsonData, ok := v.Content["application/json"]; ok {
						if sch := jsonData.Schema.Schema; sch != nil {
							responseStructs = g.processResponse(pkg.GetResponseName(path), sch.Properties, []goModel.ResponseStruct{})
						}
					}
				}
			}

			// Controller逻辑

		}
		if item.Get != nil {
			method = "GET"
			getRequest = g.ProcessGetRequest(pkg.GetRequestName(path), item.Get.Parameters)
			structs = append(structs, getRequest)

			for k, v := range item.Get.Responses {
				if k == "200" {
					if jsonData, ok := v.Content["application/json"]; ok {
						if sch := jsonData.Schema.Schema; sch != nil {
							responseStructs = g.processResponse(pkg.GetResponseName(path), sch.Properties, []goModel.ResponseStruct{})
						}
					}
				}
			}
		}

		router.Path = path
		router.Method = method
		router.FuncName = pkg.GetFuncName(path)
		router.ShortPath = packName
		if routers[pkg.GetTopLevelName(path)] == nil {
			routers[pkg.GetTopLevelName(path)] = &goModel.Router{
				Items:   make([]goModel.RouterItem, 0),
				Imports: make([]string, 0),
			}
		}
		routers[pkg.GetTopLevelName(path)].Items = append(routers[pkg.GetTopLevelName(path)].Items, router)
		routers[pkg.GetTopLevelName(path)].Imports = append(routers[pkg.GetTopLevelName(path)].Imports, controllerImportPath)
		sort.Slice(structs, func(i, j int) bool {
			return structs[i].Name < structs[j].Name
		})
		for k, v := range structs {
			var rows = v.Rows
			sort.Slice(rows, func(i, j int) bool {
				return rows[i].Name < rows[j].Name
			})
			structs[k].Rows = rows
		}
		sort.Slice(responseStructs, func(i, j int) bool {
			return responseStructs[i].Name < responseStructs[j].Name
		})

		for k, v := range responseStructs {
			var rows = v.Rows
			sort.Slice(rows, func(i, j int) bool {
				return rows[i].Name < rows[j].Name
			})
			responseStructs[k].Rows = rows
		}
		generator.Request(g.config.OutPath, "server/request/"+path+".go", FuncMap(), map[string]interface{}{
			"structs":     structs,
			"packageName": packName,
		})

		generator.Response(g.config.OutPath, "server/response/"+path+".go", FuncMap(), map[string]interface{}{
			"structs":     responseStructs,
			"packageName": packName,
		})

		generator.Controller(g.config.OutPath, "server/controller/"+path+".go", FuncMap(), map[string]interface{}{
			"importData": map[string]string{
				reqShortPath: reqImportPath,
				svcShortPath: svcImportPath,
			},
			"reqShortPath": reqShortPath,
			"svcShortPath": svcShortPath,
			"packageName":  packName,
			"reqName":      pkg.GetRequestName(path),
			"funcName":     pkg.GetFuncName(path),
		})

		generator.Service(g.config.OutPath, "service/"+path+".go", FuncMap(), map[string]interface{}{
			"importData": map[string]string{
				reqShortPath:  reqImportPath,
				respShortPath: respImportPath,
			},
			"reqShortPath":  reqShortPath,
			"packageName":   packName,
			"funcName":      pkg.GetFuncName(path),
			"respShortPath": respShortPath,
		})
	}

	for k, v := range routers {
		var items = v.Items
		sort.Slice(items, func(i, j int) bool {
			return items[i].Path < items[j].Path
		})
		generator.Router(g.config.OutPath, "server/router/"+k+".go", FuncMap(), map[string]interface{}{
			"routerName": k,
			"importData": v.Imports,
			"routers":    v.Items,
		})
	}

}

func (g *Golang) GoTypeMap(t string) string {
	var m = map[string]string{
		"integer": "int64",
		"number":  "float64",
		"boolean": "bool",
		"string":  "string",
	}
	if str, ok := m[t]; ok {
		return str
	}
	return t
}

// ProcessGetRequest 处理get请求
// TODO 支持GET中的数组类型
func (g *Golang) ProcessGetRequest(name string, parameters []model.Parameter) goModel.RequestStruct {
	var st goModel.RequestStruct
	st.Name = name
	for _, param := range parameters {
		var req goModel.RequestRow
		req.Validate, req.Description = pkg.FormatDescription(param.Description)
		req.Validate = strings.Join([]string{"required", req.Validate}, ",")
		if temp := strings.Trim(req.Validate, ","); temp != "" {
			req.Validate = temp
		}
		req.DataType = g.GoTypeMap(param.Schema.Type)
		req.Name = param.Name
		req.BindType = "form"
		st.Rows = append(st.Rows, req)
	}

	return st
}

func (g *Golang) processPostRequest(bindType, name string, schema model.SchemaProperties, structs []goModel.RequestStruct, required []string) []goModel.RequestStruct {
	var st goModel.RequestStruct
	st.Name = name
	in, err := pkg.NewQuickInArray(required)
	if err != nil {
		log.Println(err)
	}
	for k, v := range schema {
		var req goModel.RequestRow
		req.Validate, req.Description = pkg.FormatDescription(v.Description)
		if in.InArray(k) {
			req.Validate = strings.Join([]string{"required", req.Validate}, ",")
		} else {
			if req.Validate != "" {
				req.Validate = strings.Join([]string{"omitempty", req.Validate}, ",")
			}
		}
		if temp := strings.Trim(req.Validate, ","); temp != "" {
			req.Validate = temp
		}
		req.DataType = g.GoTypeMap(v.Type)
		if req.DataType == "object" {
			req.DataType = name + pkg.LineToUpCamel(k)
			structs = g.processPostRequest(bindType, name+pkg.LineToUpCamel(k), v.Properties, structs, v.Required)
		}
		if req.DataType == "array" {
			var items *model.Schema
			req.DataType, items = g.arrayType(v.Items, k, "[]")
			if items != nil {
				structs = g.processPostRequest(bindType, name, items.Properties, structs, items.Required)
			}
		}
		req.Name = k
		req.BindType = bindType
		st.Rows = append(st.Rows, req)
	}
	structs = append(structs, st)
	return structs
}

func (g *Golang) processResponse(name string, schema model.SchemaProperties, structs []goModel.ResponseStruct) []goModel.ResponseStruct {
	var st goModel.ResponseStruct
	st.Name = name
	for k, v := range schema {
		var resp goModel.ResponseRow
		resp.DataType = g.GoTypeMap(v.Type)
		if resp.DataType == "object" {
			resp.DataType = name + pkg.LineToUpCamel(k)
			structs = g.processResponse(name+pkg.LineToUpCamel(k), v.Properties, structs)
		}
		if resp.DataType == "array" {
			var items *model.Schema
			resp.DataType, items = g.arrayType(v.Items, name+pkg.LineToUpCamel(k), "[]")
			if items != nil {
				structs = g.processResponse(name+pkg.LineToUpCamel(k), items.Properties, structs)
			}
		}
		resp.Name = k
		st.Rows = append(st.Rows, resp)
	}
	structs = append(structs, st)
	return structs
}

func (g *Golang) arrayType(t *model.SchemaOrArray, parentName string, prefix string) (string, *model.Schema) {
	if t == nil {
		return "[]interface{}", nil
	}
	if t.Schema == nil {
		return "[]interface{}", nil
	}
	if t.Schema.Type == "object" {
		return prefix + parentName, t.Schema
	}
	if t.Schema.Type == "array" {
		return g.arrayType(t.Schema.Items, parentName, prefix+"[]")
	}
	return prefix + g.GoTypeMap(t.Schema.Type), nil
}

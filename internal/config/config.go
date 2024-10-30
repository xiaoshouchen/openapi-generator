package config

type Config struct {
	OpenAPIFile string           // OpenAPI文档路径
	OutputDir   string           // 输出目录
	PackageName string           // 生成代码的包名
	Generation  GenerationConfig // 生成选项
	Templates   TemplateConfig   // 模板配置
}

type GenerationConfig struct {
	GenerateModel  bool // 是否生成数据模型
	GenerateClient bool // 是否生成HTTP客户端
	GenerateServer bool // 是否生成服务端代码
	GenerateRouter bool // 是否生成路由
	UseValidation  bool // 是否生成验证代码
	UseMiddleware  bool // 是否生成中间件
}

type TemplateConfig struct {
	ModelTemplate  string // 模型模板路径
	ClientTemplate string // 客户端模板路径
	ServerTemplate string // 服务端模板路径
	RouterTemplate string // 路由模板路径
}

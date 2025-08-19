package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/xiaoshouchen/openapi-generator/internal/fetcher"
	"github.com/xiaoshouchen/openapi-generator/internal/generator"
	"github.com/xiaoshouchen/openapi-generator/internal/model"
	"github.com/xiaoshouchen/openapi-generator/internal/parser"
	"github.com/xiaoshouchen/openapi-generator/internal/process"
)

var configPath = flag.String("f", "openapi-conf.json", "配置文件")

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 && args[0] == "tmpl" {
		generator.InitTemplates()
		return
	}

	file, err := os.ReadFile(*configPath)
	if err != nil {
		return
	}
	var conf model.Config
	if err = json.Unmarshal(file, &conf); err != nil {
		log.Fatal("配置解析错误", err)
	}
	// 获取数据
	fetch, err := fetcher.NewFetcher(conf.Fetcher)
	if err != nil {
		log.Fatal(err)
	}
	// 解析数据
	rawData, err := fetch.Bytes()
	if err != nil {
		log.Fatal(err)
	}
	parserTemp := parser.NewParser(rawData, "json")
	schemaData, err := parserTemp.Parse()
	if err != nil {
		log.Fatal("parser error", err.Error())
	}

	// 加工生成数据
	process.NewProcessor(conf).Process(schemaData, generator.NewGenerator(conf))
}

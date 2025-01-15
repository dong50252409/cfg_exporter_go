package render

import (
	"cfg_exporter/entities"
	"fmt"
	"os"
)

type IRender interface {
	ExportDir() string
	Execute() error
	Filename() string
	ConfigName() string
}

var renderRegistry = make(map[string]func(tbl *entities.Table) IRender)

func Register(key string, cls func(tbl *entities.Table) IRender) {
	renderRegistry[key] = cls
}

func ToFile(schemaName string, table *entities.Table) error {
	cls, ok := renderRegistry[schemaName]
	if !ok {
		panic(fmt.Sprintf("配置表：%s 渲染模板：%s 还没有被支持", table.Filename, schemaName))
	}
	r := cls(table)
	dir := r.ExportDir()
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return r.Execute()
}

package render

import (
	"cfg_exporter/entities"
	"cfg_exporter/interfaces"
	"fmt"
	"os"
)

var registry = make(map[string]func(tbl *entities.Table) interfaces.IRender)

func Register(key string, cls func(tbl *entities.Table) interfaces.IRender) {
	registry[key] = cls
}

func ToFile(schemaName string, table *entities.Table) error {
	cls, ok := registry[schemaName]
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

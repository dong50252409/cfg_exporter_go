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
		return fmt.Errorf("配置表：%s 渲染模板：%s 还没有被支持", table.Filename, schemaName)
	}

	fmt.Printf("开始生成配置：%s\n", table.Filename)
	r := cls(table)
	dir := r.ExportDir()

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("导出路径创建失败 %s", err)
	}

	if err := r.Execute(); err != nil {
		return err
	}

	fmt.Printf("配置导出完成：%s\n", table.Filename)
	return nil
}

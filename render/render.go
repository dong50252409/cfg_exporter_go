package render

import (
	"cfg_exporter/constrainer"
	"cfg_exporter/model"
	"cfg_exporter/typesystem"
	"fmt"
	"os"
)

var registry = make(map[string]func(table renderTable) model.Render)

func register(key string, cls func(table renderTable) model.Render) {
	registry[key] = cls
}

type renderTable struct {
	*model.Table
}

func ToFile(schemaName string, table *model.Table) {
	cls, ok := registry[schemaName]
	if !ok {
		panic(fmt.Sprintf("配置表：%s 渲染模板：%s 还没有被支持", table.Filename, schemaName))
	}
	r := cls(renderTable{table})
	dir := r.ExportDir()
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	r.Execute()
}

func (tbl *renderTable) GetPrimaryKeyFields() []*model.Field {
	var fields []*model.Field
	for _, field := range tbl.Fields {
		for _, constraint := range field.Constraints {
			_, ok := constraint.(*constrainer.PrimaryKey)
			if ok {
				continue
			}
			fields = append(fields, field)
		}
	}
	return fields
}

// GetMacroFields 获取宏字段
func (tbl *renderTable) GetMacroFields() []*model.Field {
	var fields []*model.Field
	for _, field := range tbl.Fields {
		_, ok := field.Type.(*typesystem.Macro)
		if ok {
			fields = append(fields, field)
		}
	}
	return fields
}

// GetFieldByName 获取字段
func (tbl *renderTable) GetFieldByName(name string) *model.Field {
	for _, field := range tbl.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

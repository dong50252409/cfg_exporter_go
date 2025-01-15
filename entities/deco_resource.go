package entities

import (
	"cfg_exporter/config"
	"cfg_exporter/util"
	"fmt"
	"os"
	"path/filepath"
)

// Resource 资源引用
type Resource struct {
	Path string
}

func init() {
	decoratorRegister("resource", newResource)
}

func newResource(_ *Table, field *Field, str string) error {
	args := util.SubArgs(str, ",")
	if len(args) == 1 {
		if str != "" {
			wd, _ := os.Getwd()
			path := filepath.Join(wd, args[0])
			_, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("参数路径不存在 完整路径：%s", path)
			}
			field.Decorators["resource"] = &Resource{Path: path}
			return nil
		}
	}
	return fmt.Errorf("参数格式错误 resource(路径)")
}

func (r *Resource) Name() string {
	return "resource"
}

func (r *Resource) RunFieldDecorator(tbl *Table, field *Field) error {
	for rowIndex, row := range tbl.DataSet {
		v := row[field.ColIndex]
		if v == nil || v == "" {
			continue
		}
		_, err := os.Stat(filepath.Join(r.Path, v.(string)))
		if err != nil {
			return fmt.Errorf("第 %d 行 资源不存在 %s", rowIndex+config.Config.BodyStartRow, v)
		}
	}
	return nil
}

package decorator

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Resource 资源引用
type Resource struct {
	path string
}

func init() {
	registry["resource"] = newResource
}

func newResource(_ *entities.Table, field *entities.Field, str string) error {
	if str != "" {
		_, err := os.Stat(str)
		if err != nil {
			return errors.New("参数路径不存在 resource(路径)")
		}
		field.Decorators["resource"] = &Resource{path: str}
		return nil
	}

	return errors.New("参数格式错误 resource(路径)")
}

func (r *Resource) Check() bool {
	if _, err := os.Stat(r.path); err != nil {
		return false
	}
	return true
}

func (r *Resource) RunFieldDecorator(tbl *entities.Table, field *entities.Field) error {
	for rowIndex, row := range tbl.DataSet {
		v := row[field.Column]
		if v == nil || v == "" {
			continue
		}
		_, err := os.Stat(filepath.Join(r.path, v.(string)))
		if err != nil {
			return fmt.Errorf("第 %d 行 资源不存在 %s", rowIndex+config.Config.BodyStartRow, v)
		}
	}
	return nil
}

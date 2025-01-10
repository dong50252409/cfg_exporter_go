package decorator

import (
	"cfg_exporter/entities"
	"errors"
	"os"
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

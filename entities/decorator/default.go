package decorator

import (
	"cfg_exporter/entities"
	"cfg_exporter/entities/base_type"
	"cfg_exporter/util"
	"fmt"
)

type Default struct {
	DefaultValue any
}

func init() {
	register("default", newDefault)
}

func newDefault(_ *entities.Table, field *entities.Field, str string) error {
	if args := util.SubArgs(str, ""); args != nil {
		v, err := base_type.ParseString(args[0])
		if err == nil && field.Type.Equal(v) {
			field.Decorators["default"] = &Default{DefaultValue: v}
		}
	}
	return fmt.Errorf("参数格式错误 default(默认值)")
}

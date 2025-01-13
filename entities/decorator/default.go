package decorator

import (
	"cfg_exporter/entities"
	"cfg_exporter/entities/typesystem"
	"cfg_exporter/interfaces"
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
	args := util.SubArgs(str, ",")
	if len(args) == 1 {
		v, err := typesystem.ParseString(args[0])
		if err == nil {
			field.Decorators["default"] = &Default{v}
			return nil
		}

	}
	return fmt.Errorf("参数格式错误 default(默认值)")
}

func (*Default) Name() string {
	return "default"
}

func (d *Default) RunFieldDecorator(_ *entities.Table, field *entities.Field) error {
	err := field.Type.(interfaces.ITypeSystem).SetDefaultValue(d.DefaultValue)
	if err != nil {
		return err
	}
	return nil
}

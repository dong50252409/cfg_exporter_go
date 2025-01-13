package decorator

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"fmt"
)

// NotNull 非空
type NotNull struct {
}

func init() {
	register("not_null", newNotNull)
}

func newNotNull(_ *entities.Table, field *entities.Field, _ string) error {
	field.Decorators["not_null"] = &NotNull{}
	return nil
}

func (*NotNull) Name() string {
	return "not_null"
}

func (*NotNull) RunFieldDecorator(tbl *entities.Table, field *entities.Field) error {
	_, ok := field.Decorators["not_null"]
	for rowIndex, row := range tbl.DataSet {
		v := row[field.ColIndex]
		if v == nil || v == "" || !ok {
			return fmt.Errorf("第 %d 行 数值不能为空", rowIndex+config.Config.BodyStartRow)
		}
	}
	return nil
}

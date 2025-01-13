package decorator

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"fmt"
)

type Unique struct {
}

func init() {
	register("u_key", newUnique)
}

func newUnique(_ *entities.Table, field *entities.Field, _ string) error {
	field.Decorators["u_key"] = &Unique{}
	return nil
}

func (*Unique) Name() string {
	return "u_key"
}

func (*Unique) RunFieldDecorator(tbl *entities.Table, field *entities.Field) error {
	var set = make(map[any]struct{})
	for rowIndex, row := range tbl.DataSet {
		v := row[field.Column]
		if v == nil {
			continue
		}
		if _, ok := set[v]; ok {
			return fmt.Errorf("第 %d 行 数值重复", rowIndex+config.Config.BodyStartRow)
		}
		set[v] = struct{}{}
	}
	return nil
}

package decorator

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/interfaces"
	"fmt"
)

// PrimaryKey 主键
type PrimaryKey struct {
	Fields []*entities.Field
}

func init() {
	register("p_key", newPrimaryKey)
}

func newPrimaryKey(tbl *entities.Table, field *entities.Field, _ string) error {
	var pk *PrimaryKey
	for _, d := range tbl.Decorators {
		if d1, ok := d.(*PrimaryKey); ok {
			pk = d1
			break
		}
	}
	if pk == nil {
		pk = &PrimaryKey{}
		tbl.Decorators = append(tbl.Decorators, pk)
	}
	pk.Fields = append(pk.Fields, field)
	return nil
}

func (*PrimaryKey) Name() string {
	return "p_key"
}

func (pk *PrimaryKey) RunTableDecorator(tbl *entities.Table) error {
	var set = make(map[interfaces.TupleT]struct{})
	for rowIndex, row := range tbl.DataSet {
		var items interfaces.TupleT
		for index, field := range pk.Fields {
			item := row[field.ColIndex]
			if item == nil {
				return fmt.Errorf("第 %d 行 主键不能为空", rowIndex+config.Config.BodyStartRow)
			}
			items[index] = item
		}
		if _, ok := set[items]; ok {
			return fmt.Errorf("第 %d 行 主键重复 %v", rowIndex+config.Config.BodyStartRow, items)
		} else {
			set[items] = struct{}{}
		}
	}
	return nil
}

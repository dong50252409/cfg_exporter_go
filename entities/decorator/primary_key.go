package decorator

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"fmt"
)

// PrimaryKey 主键
type PrimaryKey struct {
	columnIndies []int
}

func init() {
	register("p_key", newPrimaryKey)
}

func newPrimaryKey(tbl *entities.Table, field *entities.Field, _ string) error {
	if v, ok := tbl.Decorators["p_key"]; ok {
		pk := v.(*PrimaryKey)
		pk.columnIndies = append(pk.columnIndies, field.ColIndex)
	}
	return nil
}

func (*PrimaryKey) RunTableDecorator(tbl *entities.Table) error {
	d, ok := tbl.Decorators["p_key"]
	if !ok {
		return fmt.Errorf("配置表主键不存在")
	}

	var set map[any]struct{}
	pkDecorator := d.(*PrimaryKey)
	for rowIndex, row := range tbl.DataSet {
		var items []any
		for _, colIndex := range pkDecorator.columnIndies {
			item := row[colIndex]
			if item == nil {
				return fmt.Errorf("第 %d 行 主键不能为空", rowIndex+config.Config.BodyStartRow)
			}
			items = append(items, item)
		}
		if _, ok := set[items]; ok {
			return fmt.Errorf("第 %d 行 主键重复 %v", rowIndex+config.Config.BodyStartRow, items)
		}
	}
	return nil
}

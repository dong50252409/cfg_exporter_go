package entities

import (
	"cfg_exporter/config"
	"cfg_exporter/util"
	"fmt"
)

// PrimaryKey 主键
type PrimaryKey struct {
	Fields []*Field
}
type pair struct {
	field *Field
	value interface{}
}

func init() {
	decoratorRegister("p_key", newPrimaryKey)
}

func newPrimaryKey(tbl *Table, field *Field, _ string) error {
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
	field.IsPrimaryKey = true
	pk.Fields = append(pk.Fields, field)
	return nil
}

func (*PrimaryKey) Name() string {
	return "p_key"
}

func (pk *PrimaryKey) RunTableDecorator(tbl *Table) error {
	var set = make(map[TupleT]struct{})
	for rowIndex, row := range tbl.DataSet {
		var items TupleT
		for index, field := range pk.Fields {
			item := row[field.ColIndex]
			if item == nil {
				return fmt.Errorf("单元格：%s 主键不能为空", util.ToCell(rowIndex+config.Config.BodyStartRow, field.Column))
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

package decorator

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/util"
	"fmt"
)

type MacroDetail struct {
	Key     string
	Value   any
	Comment string
}

type Macro struct {
	KeyField     *entities.Field
	ValueField   *entities.Field
	CommentField *entities.Field
	List         []MacroDetail
}

func init() {
	register("macro", newMacro)
}

func newMacro(tbl *entities.Table, field *entities.Field, str string) error {
	args := util.SubArgs(str, ",")
	if len(args) == 1 {
		valueFieldName := args[0]
		valueField := tbl.GetFieldByName(valueFieldName)
		if valueField == nil {
			return fmt.Errorf("%s 字段不存在", valueFieldName)
		}

		tbl.Decorators = append(tbl.Decorators, &Macro{KeyField: field, ValueField: valueField})
		return nil
	}

	if len(args) == 2 {
		valueFieldName, CommentFieldName := args[0], args[1]
		valueField := tbl.GetFieldByName(valueFieldName)
		if valueField == nil {
			return fmt.Errorf("%s 值字段不存在", valueFieldName)
		}

		commentField := tbl.GetFieldByName(CommentFieldName)
		if commentField == nil {
			return fmt.Errorf("%s 描述字段不存在", CommentFieldName)
		}

		tbl.Decorators = append(tbl.Decorators, &Macro{KeyField: field, ValueField: valueField, CommentField: commentField})
		return nil
	}
	return fmt.Errorf("参数格式错误 macro(值字段名[,描述字段名])")
}

func (m *Macro) Name() string {
	return "macro"
}

func (m *Macro) RunTableDecorator(tbl *entities.Table) error {
	for index, row := range tbl.DataSet {
		recordRow := tbl.Records[index+config.Config.BodyStartRow-1]
		if len(recordRow) <= m.KeyField.Column {
			continue
		}
		key := recordRow[m.KeyField.Column]
		value := row[m.ValueField.ColIndex]
		var comment string
		if m.CommentField != nil {
			comment = row[m.CommentField.ColIndex].(string)
		}
		m.List = append(m.List, MacroDetail{
			Key:     key,
			Value:   value,
			Comment: comment,
		})
	}
	return nil
}

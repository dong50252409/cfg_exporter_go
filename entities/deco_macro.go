package entities

import (
	"cfg_exporter/config"
	"cfg_exporter/util"
	"fmt"
)

type MacroDetail struct {
	Key     string
	Value   any
	Comment string
}

type Macro struct {
	MacroName    string
	KeyField     *Field
	ValueField   *Field
	CommentField *Field
	List         []MacroDetail
}

func init() {
	decoratorRegister("macro", newMacro)
}

func newMacro(tbl *Table, field *Field, str string) error {
	args := util.SubArgs(str, ",")
	if len(args) == 2 {
		macroName, valueFieldName := args[0], args[1]
		valueField := tbl.GetFieldByName(valueFieldName)
		if valueField == nil {
			return fmt.Errorf("%s 宏 %s 值字段不存在", macroName, valueFieldName)
		}

		tbl.Decorators = append(tbl.Decorators, &Macro{MacroName: macroName, KeyField: field, ValueField: valueField})
		return nil
	}

	if len(args) == 3 {
		macroName, valueFieldName, commentFieldName := args[0], args[1], args[2]
		valueField := tbl.GetFieldByName(valueFieldName)
		if valueField == nil {
			return fmt.Errorf("%s 宏 %s 值字段不存在", macroName, valueFieldName)
		}

		commentField := tbl.GetFieldByName(commentFieldName)
		if commentField == nil {
			return fmt.Errorf("%s 宏 %s 描述字段不存在", macroName, commentFieldName)
		}

		tbl.Decorators = append(tbl.Decorators, &Macro{MacroName: macroName, KeyField: field, ValueField: valueField, CommentField: commentField})
		return nil
	}
	return fmt.Errorf("参数格式错误 macro(宏名,值字段名[,描述字段名])")
}

func (m *Macro) Name() string {
	return "macro"
}

func (m *Macro) RunTableDecorator(tbl *Table) error {
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

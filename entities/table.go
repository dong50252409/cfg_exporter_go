package entities

type Table struct {
	// 路径
	Path string
	// 文件名
	Filename string
	// 表名
	Name string
	// 字段
	Fields []*Field
	// 装饰器
	Decorators []any
	// 主体数据
	DataSet [][]any
	// 原始数据
	Records [][]string
}

// GetFieldByName 获取字段
func (tbl *Table) GetFieldByName(fieldName string) *Field {
	for _, field := range tbl.Fields {
		if field.Name == fieldName {
			return field
		}
	}
	return nil
}

// GetFieldByColIndex 获取字段
func (tbl *Table) GetFieldByColIndex(colIndex int) *Field {
	for _, field := range tbl.Fields {
		if field.ColIndex == colIndex {
			return field
		}
	}
	return nil
}

// GetPrimaryKeyFields 获取主键字段列表
func (tbl *Table) GetPrimaryKeyFields() []*Field {
	for _, d := range tbl.Decorators {
		d1, ok := d.(*PrimaryKey)
		if ok {
			return d1.Fields
		}
	}
	return []*Field{}
}

// GetPrimaryKeyValues 获取主键值列表
func (tbl *Table) GetPrimaryKeyValues() [][]any {
	fields := tbl.GetPrimaryKeyFields()
	var list = make([][]any, 0, len(tbl.DataSet))
	for _, dataRow := range tbl.DataSet {
		var items []any
		for _, field := range fields {
			items = append(items, dataRow[field.ColIndex])
		}
		list = append(list, items)
	}
	return list
}

// GetPrimaryKeyValuesByString 获取主键值列表,并将值转为字符串
func (tbl *Table) GetPrimaryKeyValuesByString() [][]string {
	fields := tbl.GetPrimaryKeyFields()
	var list = make([][]string, 0, len(tbl.DataSet))
	for _, dataRow := range tbl.DataSet {
		var items []string
		for _, field := range fields {
			v := dataRow[field.ColIndex]
			items = append(items, field.Type.Convert(v))
		}
		list = append(list, items)
	}
	return list
}

// GetMacroDecorators 获取宏装饰器集合列表
func (tbl *Table) GetMacroDecorators() []*Macro {
	var macroList []*Macro
	for _, d := range tbl.Decorators {
		if d1, ok := d.(*Macro); ok {
			macroList = append(macroList, d1)
		}
	}
	return macroList
}

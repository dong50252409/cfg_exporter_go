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

type Field struct {
	// 字段在表中所在列
	Column int
	// 字段索引列
	ColIndex int
	// 字段名
	Name string
	// 字段类型
	Type any
	// 字段描述
	Comment string
	// 装饰器
	Decorators map[string]any
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

//func (tbl *Table) GetPrimaryKeyFields() []*Field {
//	for _, d := range tbl.Decorators {
//		d1, ok := d.(*decorator.PrimaryKey)
//		if ok {
//			return d1.Fields
//		}
//	}
//	return []*Field{}
//}

//// GetMacroFields 获取宏字段
//func (tbl *Table) GetMacroFields() []*decorator.Macro {
//	var macroList []*decorator.Macro
//	for _, d := range tbl.Decorators {
//		if d1, ok := d.(*decorator.Macro); ok {
//			macroList = append(macroList, d1)
//		}
//	}
//	return macroList
//}

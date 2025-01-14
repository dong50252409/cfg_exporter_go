package interfaces

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

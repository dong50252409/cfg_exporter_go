package base_type

type TypeParser interface {
	// ParseFromString 将其他语言中的数据转换为Golang中的数据
	ParseFromString(str string, args ...any) (any, error)
}

type TypeConvert interface {
	// Convert 将Golang中的数据转换为其他语言中的数据
	Convert(val any, args ...any) string
}

type TypeString interface {
	// String 将Golang中的类型转换为其他语言中的类型字符串
	String() string
}

// TupleT 元组 最多支持10个元素
type TupleT [10]interface{}

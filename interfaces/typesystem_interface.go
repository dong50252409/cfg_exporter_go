package interfaces

import "reflect"

type ITypeSystem interface {
	// ParseString 解析字符串为Golang数据
	ParseString(str string) (any, error)

	// Convert 将Golang数据转换为其他语言中的数据字符串
	Convert(val any) string

	// String 将Golang类型转换为其他语言中的类型字符串
	String() string
	// SetDefaultValue 设置类型默认值
	SetDefaultValue(val any) error

	// GetDefaultValue 获取类型默认值字符串
	GetDefaultValue() string

	GetKind() reflect.Kind
}

// TupleT 元组 最多支持10个元素
type TupleT [10]interface{}

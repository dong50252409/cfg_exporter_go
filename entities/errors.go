package entities

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorTypeNotSupported 类型不支持
var ErrorTypeNotSupported = errors.New("%s 不支持的类型")

func errorTypeNotSupported(typeStr string) error {
	return fmt.Errorf(ErrorTypeNotSupported.Error(), typeStr)
}

// ErrorTypeBaseInvalid 类型格式错误
func ErrorTypeBaseInvalid(t ITypeSystem, lst []string) error {
	return fmt.Errorf("类型格式错误 %s|%s(%s)", t.String(), t.String(), strings.Join(lst, "|"))
}

// ErrorTypeListInvalid 类型格式错误
func ErrorTypeListInvalid() error {
	return fmt.Errorf("类型格式错误 list|list(元素类型)")
}

// ErrorTypeTupleInvalid 元组类型格式错误
func ErrorTypeTupleInvalid() error {
	return fmt.Errorf("类型格式错误 tuple|tuple(元素类型)")
}

// ErrorTypeMapInvalid 字典类型格式错误
func ErrorTypeMapInvalid() error {
	return fmt.Errorf("类型格式错误 map|map(键元素类型, 值元素类型)")
}

// ErrorTypeParseFailed 类型解析失败
func ErrorTypeParseFailed(parentType ITypeSystem, str string) error {
	return fmt.Errorf("%s 类型无法解析 %s", parentType.String(), str)
}

// ErrorTypeNotMatch 类型匹配
func ErrorTypeNotMatch(t ITypeSystem, index int, element any) error {
	return fmt.Errorf("第 %d 个元素 %v 与类型 %s 不匹配", index+1, element, formatTypeSign(t))
}

// ErrorTypeMapKeyNotMatch 键类型匹配
func ErrorTypeMapKeyNotMatch(t ITypeSystem, key any) error {
	return fmt.Errorf("键:%v 与类型 %s 不匹配", key, formatTypeSign(t))
}

// ErrorTypeMapValueNotMatch 值类型匹配
func ErrorTypeMapValueNotMatch(t ITypeSystem, val any) error {
	return fmt.Errorf("值:%v 与类型 %s 不匹配", val, formatTypeSign(t))
}

func formatTypeSign(t ITypeSystem) string {
	switch t.(type) {
	case *List:
		return fmt.Sprintf("%s(%s)", t.String(), formatTypeSign(t.(*List).t))
	case *Tuple:
		return fmt.Sprintf("%s(%s)", t.String(), formatTypeSign(t.(*Tuple).t))
	case *Map:
		return fmt.Sprintf("%s(%s,%s)", t.String(), formatTypeSign(t.(*Map).keyT), formatTypeSign(t.(*Map).valueT))
	default:
		return t.String()
	}
}

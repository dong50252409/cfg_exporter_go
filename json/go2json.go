package json

import (
	"cfg_exporter/entities"
	"fmt"
	"strconv"
)

func convert(data any) any {
	switch v := data.(type) {
	case entities.RawT: // 处理原始值，在Erlang中就是atom
		return v
	case []interface{}: // 处理列表
		var elements []any
		for _, item := range v {
			elements = append(elements, convert(item))
		}
		return elements
	case map[any]any: // 处理Map
		elements := make(map[string]any, len(v))
		for key, value := range v {
			elements[toString(key)] = convert(value)
		}
		return elements
	case entities.TupleT: // 处理数组
		var elements []any
		for _, item := range v {
			if item == nil {
				break
			}
			elements = append(elements, toString(item))
		}
		return elements
	default:
		return v
	}
}

func toString(data any) string {
	switch v := data.(type) {
	case int8, int16, int32, int64, int:
		return strconv.FormatInt(v.(int64), 10)
	case float32, float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case string:
		return v
	case entities.RawT:
		return string(v)
	default:
		panic(fmt.Sprintf("%v 无法作为键来使用", v))
	}
}

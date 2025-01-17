package json

import (
	"cfg_exporter/entities"
	"fmt"
	"github.com/stoewer/go-strcase"
	"strconv"
)

func convert(data any) any {
	switch v := data.(type) {
	case string:
		return strcase.LowerCamelCase(v)
	case entities.RawT:
		return strcase.LowerCamelCase(string(v))
	case []interface{}:
		var elements []any
		for _, item := range v {
			elements = append(elements, convert(item))
		}
		return elements
	case map[any]any:
		elements := make(map[string]any, len(v))
		for key, value := range v {
			elements[toString(key)] = convert(value)
		}
		return elements
	case entities.TupleT:
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

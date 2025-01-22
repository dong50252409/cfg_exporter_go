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
	case entities.AnyT:
		return strcase.LowerCamelCase(string(v))
	case entities.TupleT:
		var elements = make([]any, 0, len(v))
		for _, item := range v {
			if item == nil {
				break
			}
			switch item.(type) {
			case entities.TupleT, []any:
				elements = append(elements, map[string]any{"e": convert(item)})
			default:
				elements = append(elements, convert(item))
			}
		}
		return elements
	case []interface{}:
		var elements = make([]any, 0, len(v))
		for _, item := range v {
			switch item.(type) {
			case entities.TupleT, []any:
				elements = append(elements, map[string]any{"e": convert(item)})
			default:
				elements = append(elements, convert(item))
			}
		}
		return elements
	case map[any]any:
		elements := make([]any, 0, len(v))
		for key, value := range v {
			if v1 := convert(value); v1 != nil {
				elements = append(elements, map[string]any{
					"k": toString(key),
					"v": v1,
				})
			}
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
	case entities.AnyT:
		return string(v)
	default:
		panic(fmt.Sprintf("%v 无法作为键来使用", v))
	}
}

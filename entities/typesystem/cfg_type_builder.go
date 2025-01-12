package typesystem

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	registry = make(map[string]func(typeStr string) (any, error))
)

// Register 类型注册器
func Register(key string, cls func(typeStr string) (any, error)) {
	registry[key] = cls
}

func New(typeStr string) (any, error) {
	key, args := getKey(typeStr)
	if cls, ok := registry[key]; ok {
		return cls(args)
	}
	return nil, fmt.Errorf("类型不存在 %s", key)
}

func getKey(typeStr string) (string, string) {
	index := strings.Index(typeStr, "(")
	if index != -1 {
		return typeStr[:index], typeStr[index+1:]
	} else {
		return typeStr, ""
	}
}

func GetTypeKind(t any) reflect.Kind {
	switch t.(type) {
	case *Boolean:
		return reflect.Bool
	case *Integer:
		return reflect.Int
	case *Float:
		return reflect.Float64
	case *Str:
		return reflect.String
	case *List:
		return reflect.Slice
	case *Tuple:
		return reflect.Array
	case *Map:
		return reflect.Map
	case *Lang:
		return reflect.String
	default:
		return reflect.Interface
	}
}

func ConvertToType(v any, t any) {
	switch t.(type) {
	case *Boolean:
		return v.(bool)
	case *Integer:
		return reflect.Int
	case *Float:
		return reflect.Float64
	case *Str:
		return reflect.String
	case *List:
		return reflect.Slice
	case *Tuple:
		return reflect.Array
	case *Map:
		return reflect.Map
	case *Lang:
		return reflect.String
	default:
		return reflect.Interface
	}
}

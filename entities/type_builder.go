package entities

import (
	"cfg_exporter/util"
	"math"
	"reflect"
)

var (
	typeRegistry     = make(map[string]func(typeStr string) (ITypeSystem, error))
	baseTypeRegistry = make(map[string]func(typeStr string) (ITypeSystem, error))
)

// TypeRegister 类型注册器
func TypeRegister(key string, cls func(typeStr string) (ITypeSystem, error)) {
	typeRegistry[key] = cls
	baseTypeRegistry[key] = cls
}

// RecoverBaseTypeRegister 恢复基础类型注册器
func RecoverBaseTypeRegister() {
	for k, v := range baseTypeRegistry {
		typeRegistry[k] = v
	}
}

// MergerTypeRegistry 合并当前的类型注册器
func MergerTypeRegistry(registry map[string]func(typeStr string) (ITypeSystem, error)) {
	for k, v := range registry {
		typeRegistry[k] = v
	}
}

func NewType(typeStr string) (ITypeSystem, error) {
	key, args := util.GetKey(typeStr)
	if cls, ok := typeRegistry[key]; ok {
		return cls(args)
	}
	return nil, errorTypeNotSupported(key)
}

// checkFunc 检查函数
func checkFunc(t ITypeSystem) func(any) bool {
	switch t.GetKind() {
	case reflect.Int8:
		return func(v any) bool { return math.MinInt8 <= v.(int64) && v.(int64) <= math.MaxInt8 }
	case reflect.Int16:
		return func(v any) bool { return math.MinInt16 <= v.(int64) && v.(int64) <= math.MaxInt16 }
	case reflect.Int32:
		return func(v any) bool { return math.MinInt32 <= v.(int64) && v.(int64) <= math.MaxInt32 }
	case reflect.Int64:
		return func(v any) bool { _, ok := v.(int64); return ok }
	case reflect.Float32:
		return func(v any) bool { return math.SmallestNonzeroFloat32 <= v.(float64) && v.(float64) <= math.MaxFloat32 }
	case reflect.Float64:
		return func(v any) bool { _, ok := v.(float64); return ok }
	case reflect.Bool:
		return func(v any) bool { _, ok := v.(bool); return ok }
	case reflect.String:
		return func(v any) bool {
			_, ok := v.(string)
			if !ok {
				_, ok = v.(RawT)
				return ok
			}
			return ok
		}
	case reflect.Slice:
		cf := checkFunc(t.(*List).t)
		return func(v any) bool {
			v1, ok := v.([]any)
			if !ok {
				return false
			}
			for _, e := range v1 {
				if !cf(e) {
					return false
				}
			}
			return true
		}
	case reflect.Array:
		cf := checkFunc(t.(*Tuple).t)
		return func(v any) bool {
			v1, ok := v.(TupleT)
			if !ok {
				return false
			}
			for _, e := range v1 {
				if e == nil {
					continue
				}
				if !cf(e) {
					return false
				}
			}
			return true
		}
	case reflect.Map:
		keyCF := checkFunc(t.(*Map).keyT)
		valueCF := checkFunc(t.(*Map).valueT)
		return func(v any) bool {
			v1, ok := v.(map[any]any)
			if !ok {
				return false
			}
			for key, val := range v1 {
				if !keyCF(key) || !valueCF(val) {
					return false
				}
			}
			return true
		}
	default:
		return func(any) bool { return true }
	}
}

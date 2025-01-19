package entities

import (
	"cfg_exporter/util"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Map struct {
	Field  *Field
	KeyT   ITypeSystem
	ValueT ITypeSystem
}

func init() {
	TypeRegister("map", NewMap)
}

func NewMap(typeStr string, field *Field) (ITypeSystem, error) {
	if param := util.SubParam(typeStr); param == "" {
		return &Map{Field: field}, nil
	} else {
		if l := strings.Split(param, ","); len(l) == 2 {
			kT, err := NewType(l[0], field)
			if errors.Is(err, ErrorTypeNotSupported) {
				return nil, ErrorTypeMapInvalid()
			}

			vT, err := NewType(l[1], field)
			if errors.Is(err, ErrorTypeNotSupported) {
				return nil, ErrorTypeMapInvalid()
			}
			return &Map{Field: field, KeyT: kT, ValueT: vT}, nil
		}
	}
	return nil, ErrorTypeMapInvalid()
}

func (m *Map) ParseString(str string) (any, error) {
	v, err := ParseString(str)
	if err != nil {
		return nil, err
	}
	if m.KeyT != nil && m.ValueT != nil {
		keyCheckFunc := m.KeyT.GetCheckFunc()
		valueCheckFunc := m.ValueT.GetCheckFunc()
		for key, val := range v.(map[any]any) {
			if !keyCheckFunc(key) {
				return nil, ErrorTypeMapKeyNotMatch(m, key)
			}
			if !valueCheckFunc(val) {
				return nil, ErrorTypeMapValueNotMatch(m, val)
			}
		}
	}
	return v, nil
}

func (*Map) Convert(val any) string {
	var strList []string
	for k, v := range val.(map[any]any) {
		strList = append(strList, fmt.Sprintf("%v:%v", k, v))
	}
	return fmt.Sprintf("{%s}", strings.Join(strList, ","))
}

func (m *Map) String() string {
	return "map"
}

func (m *Map) GetDefaultValue() string {
	return "{}"
}

func (m *Map) GetKind() reflect.Kind {
	return reflect.Map
}

func (m *Map) GetCheckFunc() func(any) bool {
	keyCF := m.KeyT.GetCheckFunc()
	valueCF := m.ValueT.GetCheckFunc()
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
}

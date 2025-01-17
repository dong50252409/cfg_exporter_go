package entities

import (
	"cfg_exporter/util"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Map struct {
	keyT           ITypeSystem
	valueT         ITypeSystem
	keyCheckFunc   func(any) bool
	valueCheckFunc func(any) bool
}

func init() {
	TypeRegister("map", NewMap)
}

func NewMap(typeStr string) (ITypeSystem, error) {
	if param := util.SubParam(typeStr); param == "" {
		return &Map{}, nil
	} else {
		if l := strings.Split(param, ","); len(l) == 2 {
			kT, err := NewType(l[0])
			if errors.Is(err, ErrorTypeNotSupported) {
				return nil, ErrorTypeMapInvalid()
			}

			vT, err := NewType(l[1])
			if errors.Is(err, ErrorTypeNotSupported) {
				return nil, ErrorTypeMapInvalid()
			}
			return &Map{keyT: kT, valueT: vT, keyCheckFunc: checkFunc(kT), valueCheckFunc: checkFunc(vT)}, nil
		}
	}
	return nil, ErrorTypeMapInvalid()
}

func (m *Map) ParseString(str string) (any, error) {
	v, err := ParseString(str)
	if err != nil {
		return nil, err
	}
	if m.keyCheckFunc != nil && m.valueCheckFunc != nil {
		for key, val := range v.(map[any]any) {
			if !m.keyCheckFunc(key) {
				return nil, ErrorTypeMapKeyNotMatch(m, key)
			}
			if !m.valueCheckFunc(val) {
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

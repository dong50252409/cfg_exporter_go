package typesystem

import (
	"cfg_exporter/util"
	"fmt"
	"reflect"
	"strings"
)

type Map struct {
	DefaultValue string
	KeyKind      reflect.Kind
	ValueKind    reflect.Kind
}

func NewMap(typeStr string) (*Map, error) {
	args := util.SubArgs(typeStr, ",")
	switch len(args) {
	case 0:
		return &Map{DefaultValue: "{}", KeyKind: reflect.Interface, ValueKind: reflect.Interface}, nil
	case 2:
		kT, err := New(args[0])
		if err != nil {
			return nil, err
		}
		vT, err := New(args[1])
		if err != nil {
			return nil, err
		}
		return &Map{DefaultValue: "{}", KeyKind: GetTypeKind(kT), ValueKind: GetTypeKind(vT)}, nil
	}
	return nil, fmt.Errorf("类型格式错误 map|map(键元素类型, 值元素类型)")
}
func (m *Map) ParseString(str string) (any, error) {
	v, err := ParseString(str)
	if err != nil {
		return nil, err
	}
	if m.KeyKind != reflect.Interface {
		for k, v := range v.(map[any]any) {
			if m.KeyKind != reflect.TypeOf(k).Kind() {
				return nil, fmt.Errorf("键:%v 与泛型不匹配 map(%s, %s)", k, m.KeyKind, m.ValueKind)
			}
			if m.ValueKind != reflect.TypeOf(v).Kind() {
				return nil, fmt.Errorf("值:%v 与泛型不匹配 map(%s, %s)", m, m.KeyKind, m.ValueKind)
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

func (m *Map) SetDefaultValue(val any) error {
	v, ok := val.(map[any]any)
	if ok {
		m.DefaultValue = m.Convert(v)
		return nil
	}
	return fmt.Errorf("类型不匹配 %v", val)
}

func (m *Map) GetDefaultValue() string {
	return m.DefaultValue
}

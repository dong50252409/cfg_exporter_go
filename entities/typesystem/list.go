package typesystem

import (
	"cfg_exporter/interfaces"
	"cfg_exporter/util"
	"fmt"
	"reflect"
	"strings"
)

type List struct {
	DefaultValue string
	ElementKind  reflect.Kind
}

func NewList(typeStr string) (*List, error) {
	args := util.SubArgs(typeStr, ",")
	switch len(args) {
	case 0:
		return &List{DefaultValue: "[]", ElementKind: reflect.Interface}, nil
	case 1:
		t, err := New(args[0])
		if err != nil {
			return nil, err
		}
		return &List{ElementKind: t.(interfaces.ITypeSystem).GetKind()}, nil
	}
	return nil, fmt.Errorf("类型格式错误 list|list(ElementType)")
}
func (l *List) ParseString(str string) (any, error) {
	v, err := ParseString(str)
	if err != nil {
		return nil, err
	}
	if l.ElementKind != reflect.Interface {
		for i, e := range v.([]any) {
			if l.ElementKind != reflect.TypeOf(e).Kind() {
				return nil, fmt.Errorf("第 %d 个元素 %v 与泛型不匹配 list(%s)", i+1, e, l.ElementKind)
			}
		}
	}
	return v, nil
}

func (l *List) Convert(val any) string {
	var strList []string
	for _, e := range val.([]any) {
		strList = append(strList, fmt.Sprintf("%v", e))
	}
	return fmt.Sprintf("[%s]", strings.Join(strList, ","))
}

func (l *List) String() string {
	return "list"
}

func (l *List) SetDefaultValue(val any) error {
	v, ok := val.([]any)
	if ok {
		l.DefaultValue = l.Convert(v)
		return nil
	}
	return fmt.Errorf("类型不匹配 %v", val)
}

func (l *List) GetDefaultValue() string {
	return l.DefaultValue
}

func (l *List) GetKind() reflect.Kind {
	return reflect.Slice
}

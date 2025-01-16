package entities

import (
	"cfg_exporter/util"
	"fmt"
	"reflect"
	"strings"
)

type List struct {
	ElementKind reflect.Kind
}

func init() {
	TypeRegister("list", NewList)
}

func NewList(typeStr string) (ITypeSystem, error) {
	args := util.SubArgs(typeStr, ",")
	switch len(args) {
	case 0:
		return &List{ElementKind: reflect.Interface}, nil
	case 1:
		t, err := NewType(args[0])
		if err != nil {
			return nil, err
		}
		return &List{ElementKind: t.(ITypeSystem).GetKind()}, nil
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

func (l *List) GetDefaultValue() string {
	return "[]"
}

func (l *List) GetKind() reflect.Kind {
	return reflect.Slice
}

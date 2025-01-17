package entities

import (
	"cfg_exporter/util"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type List struct {
	t         ITypeSystem
	checkFunc func(any) bool
}

func init() {
	TypeRegister("list", NewList)
}

func NewList(typeStr string) (ITypeSystem, error) {
	if param := util.SubParam(typeStr); param == "" {
		return &List{}, nil
	} else {
		t, err := NewType(param)
		if errors.Is(err, ErrorTypeNotSupported) {
			return nil, ErrorTypeListInvalid()
		}
		return &List{t: t, checkFunc: checkFunc(t)}, nil
	}
}

func (l *List) ParseString(str string) (any, error) {
	v, err := ParseString(str)
	if err != nil {
		return nil, ErrorTypeParseFailed(l, str)
	}
	if l.checkFunc != nil {
		for i, e := range v.([]any) {
			if !l.checkFunc(e) {
				return nil, ErrorTypeNotMatch(l, i, e)
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

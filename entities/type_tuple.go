package entities

import (
	"cfg_exporter/util"
	"fmt"
	"reflect"
	"strings"
)

type Tuple struct {
	ElementKind reflect.Kind
}

func NewTuple(typeStr string) (*Tuple, error) {
	args := util.SubArgs(typeStr, ",")
	switch len(args) {
	case 0:
		return &Tuple{ElementKind: reflect.Interface}, nil
	case 1:
		t, err := NewType(args[0])
		if err != nil {
			return nil, err
		}
		return &Tuple{ElementKind: t.(ITypeSystem).GetKind()}, nil
	}
	return nil, fmt.Errorf("类型格式错误 tuple|tuple(元素类型)")
}
func (t *Tuple) ParseString(str string) (any, error) {
	v, err := ParseString(str)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return v, nil
	}
	if t.ElementKind != reflect.Interface {
		for i, e := range v.(TupleT) {
			if e != nil {
				if t.ElementKind != reflect.TypeOf(e).Kind() {
					return nil, fmt.Errorf("第 %d 个元素 %v 与泛型不匹配 tuple(%s)", i+1, e, t.ElementKind)
				}
			} else {
				break
			}
		}
	}
	return v, nil
}

func (*Tuple) Convert(val any) string {
	var strList []string
	for _, e := range val.(TupleT) {
		if e != nil {
			strList = append(strList, fmt.Sprintf("%v", e))
		} else {
			break
		}
	}
	return fmt.Sprintf("(%v)", strings.Join(strList, ","))
}

func (t *Tuple) String() string {
	return "array"
}

func (t *Tuple) GetDefaultValue() string {
	return "[]"
}

func (*Tuple) GetKind() reflect.Kind {
	return reflect.Array
}

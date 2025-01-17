package entities

import (
	"cfg_exporter/util"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Tuple struct {
	t         ITypeSystem
	checkFunc func(any) bool
}

func init() {
	TypeRegister("tuple", NewTuple)
}

func NewTuple(typeStr string) (ITypeSystem, error) {
	if param := util.SubParam(typeStr); param == "" {
		return &Tuple{}, nil
	} else {
		t, err := NewType(param)
		if errors.Is(err, ErrorTypeNotSupported) {
			return nil, ErrorTypeTupleInvalid()
		}
		return &Tuple{t: t}, nil
	}
}

func (t *Tuple) ParseString(str string) (any, error) {
	v, err := ParseString(str)
	if err != nil {
		return nil, ErrorTypeParseFailed(t, str)
	}
	if v == nil {
		return v, nil
	}
	if t.t != nil {
		checkFunc := t.t.GetCheckFunc()
		for i, e := range v.(TupleT) {
			if e != nil {
				if !checkFunc(e) {
					return nil, ErrorTypeNotMatch(t, i, e)
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
	return "tuple"
}

func (t *Tuple) GetDefaultValue() string {
	return "[]"
}

func (t *Tuple) GetKind() reflect.Kind {
	return reflect.Array
}

func (t *Tuple) GetCheckFunc() func(any) bool {
	cf := t.t.GetCheckFunc()
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
}

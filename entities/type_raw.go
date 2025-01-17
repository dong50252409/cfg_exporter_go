package entities

import (
	"fmt"
	"reflect"
)

type Raw struct {
}

func init() {
	TypeRegister("raw", NewRaw)
}

func NewRaw(_ string) (ITypeSystem, error) {
	return &Raw{}, nil
}

func (r *Raw) ParseString(str string) (any, error) {
	return str, nil
}

func (*Raw) Convert(val any) string {
	return fmt.Sprintf("%v", val)
}

func (r *Raw) String() string {
	return "raw"
}

func (r *Raw) GetDefaultValue() string {
	return "nil"
}

func (r *Raw) GetKind() reflect.Kind {
	return reflect.String
}

func (r *Raw) GetCheckFunc() func(any) bool {
	return func(v any) bool {
		_, ok := v.(string)
		if !ok {
			_, ok = v.(RawT)
			return ok
		}
		return ok
	}
}

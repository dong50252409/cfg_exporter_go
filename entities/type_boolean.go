package entities

import (
	"reflect"
	"strconv"
)

type Boolean struct {
}

func init() {
	TypeRegister("bool", NewBoolean)
}

func NewBoolean(_ string) (ITypeSystem, error) {
	t := &Boolean{}
	return t, nil
}

func (*Boolean) ParseString(str string) (any, error) {
	return strconv.ParseBool(str)
}

func (*Boolean) Convert(val any) string {
	return strconv.FormatBool(val.(bool))
}

func (b *Boolean) String() string {
	return "bool"
}

func (b *Boolean) GetDefaultValue() string {
	return "false"
}

func (b *Boolean) GetKind() reflect.Kind {
	return reflect.Bool
}

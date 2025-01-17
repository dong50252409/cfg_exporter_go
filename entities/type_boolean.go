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

func (b *Boolean) ParseString(str string) (any, error) {
	parseBool, err := strconv.ParseBool(str)
	if err != nil {
		return nil, ErrorTypeParseFailed(b, str)
	}
	return parseBool, nil
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

package typesystem

import (
	"fmt"
	"strconv"
)

type Boolean struct {
	DefaultValue string
}

func NewBoolean(_ string) (*Boolean, error) {
	t := &Boolean{DefaultValue: "false"}
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

func (b *Boolean) SetDefaultValue(val any) error {
	v, ok := val.(bool)
	if ok {
		b.DefaultValue = strconv.FormatBool(v)
		return nil
	}
	return fmt.Errorf("类型不匹配 %v", val)
}

func (b *Boolean) GetDefaultValue() string {
	return b.DefaultValue
}

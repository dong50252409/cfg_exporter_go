package base_type

import (
	"fmt"
	"strconv"
)

type Boolean struct {
	Value bool
}

func (b *Boolean) ParseFromString(str string, _ ...any) (any, error) {
	return strconv.ParseBool(str)
}

func (*Boolean) Convert(val any, _ ...any) string {
	return strconv.FormatBool(val.(bool))
}

func (b *Boolean) String() string {
	return fmt.Sprintf("%T", b)
}

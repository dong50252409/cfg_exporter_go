package base_type

import (
	"fmt"
	"strconv"
)

type Integer struct {
	BitSize int
}

func (i *Integer) ParseFromString(str string, _ ...any) (any, error) {
	return strconv.ParseInt(str, 10, i.BitSize)
}

func (i *Integer) Convert(val any, _ ...any) string {
	return strconv.Itoa(val.(int))
}

func (i *Integer) String() string {
	return fmt.Sprintf("%T", i)
}

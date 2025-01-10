package base_type

import (
	"fmt"
	"strconv"
)

type Float struct {
	BitSize int
}

func (f *Float) ParseFromString(str string, _ ...any) (any, error) {
	return strconv.ParseFloat(str, f.BitSize)
}

func (f *Float) Convert(val any, _ ...any) string {
	v := val.(float64)
	return strconv.FormatFloat(v, 'f', -1, f.BitSize)
}

func (f *Float) String() string {
	return fmt.Sprintf("%T", f)
}

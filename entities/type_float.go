package entities

import (
	"cfg_exporter/util"
	"fmt"
	"maps"
	"reflect"
	"strconv"
	"strings"
)

type Float struct {
	BitSize int
}

var floatByteSizes = map[string]int{
	"32": 32,
	"64": 64,
}

func NewFloat(typeStr string) (*Float, error) {
	args := util.SubArgs(typeStr, ",")
	bit := "64"
	if len(args) == 1 {
		bit = args[0]
	}
	if bytes, ok := floatByteSizes[bit]; ok {
		return &Float{BitSize: bytes}, nil
	}

	l := make([]string, 0, len(floatByteSizes))
	for k := range maps.Keys(floatByteSizes) {
		l = append(l, k)
	}
	return nil, fmt.Errorf("类型格式错误 float|float(%s)", strings.Join(l, "|"))
}

func (f *Float) ParseString(str string) (any, error) {
	return strconv.ParseFloat(str, f.BitSize)
}

func (f *Float) Convert(val any) string {
	return strconv.FormatFloat(val.(float64), 'f', -1, f.BitSize)
}

func (f *Float) String() string {
	if f.BitSize == 32 {
		return "float32"
	}
	return "float64"
}

func (f *Float) GetDefaultValue() string {
	return "0.0"
}

func (f *Float) GetKind() reflect.Kind {
	if f.BitSize == 32 {
		return reflect.Float32
	}
	return reflect.Float64
}

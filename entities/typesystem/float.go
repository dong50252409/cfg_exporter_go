package typesystem

import (
	"cfg_exporter/util"
	"fmt"
	"maps"
	"strconv"
	"strings"
)

type Float struct {
	DefaultValue string
	ByteSize     int
}

var floatByteSizes = map[string]int{
	"32": 4,
	"64": 8,
}

func NewFloat(typeStr string, DefaultValue string) (*Float, error) {
	args := util.SubArgs(typeStr, "")
	bit := "64"
	if len(args) == 1 {
		bit = args[0]
	}
	if bytes, ok := floatByteSizes[bit]; ok {
		return &Float{DefaultValue: DefaultValue, ByteSize: bytes}, nil
	}

	l := make([]string, 0, len(floatByteSizes))
	for k := range maps.Keys(floatByteSizes) {
		l = append(l, k)
	}
	return nil, fmt.Errorf("类型格式错误 float|float(%s)", strings.Join(l, "|"))
}

func (f *Float) ParseString(str string) (any, error) {
	return strconv.ParseFloat(str, f.ByteSize)
}

func (f *Float) Convert(val any) string {
	return strconv.FormatFloat(val.(float64), 'f', -1, f.ByteSize)
}

func (f *Float) String() string {
	if f.ByteSize == 32 {
		return "float32"
	}
	return "float64"
}

func (f *Float) SetDefaultValue(val any) error {
	v, ok := val.(float64)
	if ok {
		f.DefaultValue = strconv.FormatFloat(v, 'f', -1, f.ByteSize)
		return nil
	}
	return fmt.Errorf("类型不匹配 %v", val)
}

func (f *Float) GetDefaultValue() string {
	return f.DefaultValue
}

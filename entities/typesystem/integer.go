package typesystem

import (
	"cfg_exporter/util"
	"fmt"
	"maps"
	"strconv"
	"strings"
)

type Integer struct {
	DefaultValue string
	ByteSize     int
}

var intByteSizes = map[string]int{
	"8":  1,
	"16": 2,
	"32": 3,
	"64": 4,
}

func NewInteger(typeStr string, DefaultValue string) (*Integer, error) {
	args := util.SubArgs(typeStr, "")
	bit := "64"
	if len(args) == 1 {
		bit = args[0]
	}
	if bytes, ok := intByteSizes[bit]; ok {
		return &Integer{DefaultValue: DefaultValue, ByteSize: bytes}, nil
	}

	l := make([]string, 0, len(intByteSizes))
	for k := range maps.Keys(intByteSizes) {
		l = append(l, k)
	}
	return nil, fmt.Errorf("类型格式错误 int|int(%s)", strings.Join(l, "|"))
}
func (i *Integer) ParseString(str string) (any, error) {
	return strconv.ParseInt(str, 10, i.ByteSize)
}

func (i *Integer) Convert(val any) string {
	return strconv.Itoa(val.(int))
}

func (i *Integer) String() string {
	switch i.ByteSize {
	case 8:
		return "int8"
	case 16:
		return "int16"
	case 32:
		return "int32"
	case 64:
		return "int64"
	default:
		return "int"
	}
}

func (i *Integer) SetDefaultValue(val any) error {
	v, ok := val.(int)
	if ok {
		i.DefaultValue = strconv.Itoa(v)
		return nil
	}
	return fmt.Errorf("类型不匹配 %v", val)
}

func (i *Integer) GetDefaultValue() string {
	return i.DefaultValue
}

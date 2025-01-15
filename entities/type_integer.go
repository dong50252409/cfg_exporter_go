package entities

import (
	"cfg_exporter/util"
	"fmt"
	"maps"
	"reflect"
	"strconv"
	"strings"
)

type Integer struct {
	BitSize int
}

var intByteSizes = map[string]int{
	"8":  8,
	"16": 16,
	"32": 32,
	"64": 64,
}

func NewInteger(typeStr string) (*Integer, error) {
	args := util.SubArgs(typeStr, ",")
	bit := "64"
	if len(args) == 1 {
		bit = args[0]
	}
	if bytes, ok := intByteSizes[bit]; ok {
		return &Integer{BitSize: bytes}, nil
	}

	l := make([]string, 0, len(intByteSizes))
	for k := range maps.Keys(intByteSizes) {
		l = append(l, k)
	}
	return nil, fmt.Errorf("类型格式错误 int|int(%s)", strings.Join(l, "|"))
}
func (i *Integer) ParseString(str string) (any, error) {
	return strconv.ParseInt(str, 10, i.BitSize)
}

func (i *Integer) Convert(val any) string {
	return strconv.FormatInt(val.(int64), 10)
}

func (i *Integer) String() string {
	switch i.BitSize {
	case 8:
		return "int8"
	case 16:
		return "int16"
	case 32:
		return "int32"
	case 64:
		return "int64"
	default:
		return "int64"
	}
}

func (i *Integer) GetDefaultValue() string {
	return "0"
}

func (i *Integer) GetKind() reflect.Kind {
	switch i.BitSize {
	case 8:
		return reflect.Int8
	case 16:
		return reflect.Int16
	case 32:
		return reflect.Int32
	case 64:
		return reflect.Int64
	default:
		return reflect.Int64
	}
}

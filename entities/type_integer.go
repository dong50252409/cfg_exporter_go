package entities

import (
	"cfg_exporter/util"
	"maps"
	"reflect"
	"strconv"
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

func init() {
	TypeRegister("int", NewInteger)
}

func NewInteger(typeStr string) (ITypeSystem, error) {
	bit := "64"
	if param := util.SubParam(typeStr); param != "" {
		bit = param
	}
	if bytes, ok := intByteSizes[bit]; ok {
		return &Integer{BitSize: bytes}, nil
	}

	l := make([]string, 0, len(intByteSizes))
	for k := range maps.Keys(intByteSizes) {
		l = append(l, k)
	}
	return nil, ErrorTypeBaseInvalid(&Integer{}, l)
}
func (i *Integer) ParseString(str string) (any, error) {
	parseInt, err := strconv.ParseInt(str, 10, i.BitSize)
	if err != nil {
		return nil, ErrorTypeParseFailed(i, str)
	}
	return parseInt, nil
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

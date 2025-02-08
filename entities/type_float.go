package entities

import (
	"cfg_exporter/util"
	"maps"
	"math"
	"math/big"
	"reflect"
	"strconv"
)

type Float struct {
	Field   *Field
	BitSize int
}

var (
	floatDefaultBitSizes = "32"
	floatByteSizes       = map[string]int{
		"32": 32,
		"64": 64,
	}
)

func init() {
	TypeRegister("float", NewFloat)
}

func NewFloat(typeStr string, field *Field) (ITypeSystem, error) {
	bit := floatDefaultBitSizes
	if param := util.SubParam(typeStr); param != "" {
		bit = param
	}
	if bytes, ok := floatByteSizes[bit]; ok {
		return &Float{Field: field, BitSize: bytes}, nil
	}

	l := make([]string, 0, len(floatByteSizes))
	for k := range maps.Keys(floatByteSizes) {
		l = append(l, k)
	}
	return nil, ErrorTypeBaseInvalid(&Float{}, l)
}

func (f *Float) ParseString(str string) (any, error) {
	parseFloat, _, err := big.ParseFloat(str, 10, 0, big.ToNearestEven)
	if err != nil {
		return nil, ErrorTypeParseFailed(f, str)
	}
	float, _ := parseFloat.Float64()
	return float, nil
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

func (f *Float) DefaultValue() string {
	return "0.0"
}

func (f *Float) Kind() reflect.Kind {
	if f.BitSize == 32 {
		return reflect.Float32
	}
	return reflect.Float64
}

func (f *Float) CheckFunc() func(any) bool {
	if f.BitSize == 32 {
		return func(v any) bool {
			v1, ok := v.(float64)
			return ok && math.SmallestNonzeroFloat32 <= v1 && v1 <= math.MaxFloat32
		}
	}
	return func(v any) bool { _, ok := v.(float64); return ok }
}

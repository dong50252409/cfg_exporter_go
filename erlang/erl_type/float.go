package erl_type

import (
	"cfg_exporter/entities/base_type"
	"fmt"
	"maps"
	"strings"
)

var floatBitSizes = map[string]int{
	"32": 32,
	"64": 64,
}

type ErlFloat struct {
	*base_type.Float
}

func init() {
	base_type.Register("float", newFloat)
}

func newFloat(args []string) (any, error) {

	var bit string
	if args == nil {
		bit = "64"
	} else if len(args) == 1 {
		bit = args[0]
	}

	if bitSize, ok := floatBitSizes[bit]; ok {
		return &ErlFloat{Float: &base_type.Float{BitSize: bitSize}}, nil
	}

	l := make([]string, 0, len(floatBitSizes))
	for k := range maps.Keys(floatBitSizes) {
		l = append(l, k)
	}
	return nil, fmt.Errorf("类型格式错误 float|float(%s)", strings.Join(l, "|"))
}

func (f *ErlFloat) String() string {
	return "float()"
}

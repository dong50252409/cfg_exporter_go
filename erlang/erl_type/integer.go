package erl_type

import (
	"cfg_exporter/entities/base_type"
	"fmt"
	"maps"
	"strings"
)

var intBitSizes = map[string]int{
	"8":  8,
	"16": 16,
	"32": 32,
	"64": 64,
}

type ErlInteger struct {
	*base_type.Integer
}

func init() {
	base_type.Register("int", newInt)
}

func newInt(args []string) (any, error) {
	var bit string
	if args == nil {
		bit = "64"
	} else if len(args) == 1 {
		bit = args[0]
	}

	if bitSize, ok := intBitSizes[bit]; ok {
		return &ErlInteger{Integer: &base_type.Integer{BitSize: bitSize}}, nil
	}

	l := make([]string, 0, len(intBitSizes))
	for k := range maps.Keys(intBitSizes) {
		l = append(l, k)
	}
	return nil, fmt.Errorf("类型格式错误 int|int(%s)", strings.Join(l, "|"))
}

func (i *ErlInteger) String() string {
	return "integer()"
}

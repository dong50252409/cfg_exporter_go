package erl_type

import (
	"cfg_exporter/entities/base_type"
	"fmt"
)

type ErlTuple struct {
	*base_type.Tuple
}

func init() {
	base_type.Register("tuple", newTuple)
}

func newTuple(args []string) (any, error) {
	if args == nil {
		return &ErlTuple{Tuple: &base_type.Tuple{}}, nil
	}

	if len(args) == 1 {
		t, err := base_type.New(args[0])
		if err != nil {
			return nil, err
		}
		return &ErlTuple{Tuple: &base_type.Tuple{ElementType: t}}, nil
	}

	return nil, fmt.Errorf("类型格式错误 tuple|tuple(元素类型)")
}

func (*ErlTuple) Convert(val ...any) string {
	return toString(val[0])
}

func (*ErlTuple) String() string {
	return "tuple()"
}

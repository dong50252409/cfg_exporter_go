package erl_type

import (
	"cfg_exporter/entities/base_type"
	"fmt"
)

type ErlMap struct {
	*base_type.Map
}

func init() {
	base_type.Register("map", newMap)
}

func newMap(args []string) (any, error) {
	if args == nil {
		return &ErlMap{}, nil
	}

	if len(args) == 2 {
		kT, err := base_type.New(args[0])
		if err != nil {
			return nil, err
		}
		vT, err := base_type.New(args[1])
		if err != nil {
			return nil, err
		}
		return &ErlMap{Map: &base_type.Map{KeyT: kT, ValueT: vT}}, nil
	}

	return nil, fmt.Errorf("类型格式错误 map|map(键元素类型, 值元素类型)")
}

func (*ErlMap) Convert(val ...any) string {
	return toString(val[0])
}

func (m *ErlMap) String() string {
	return "map()"
}

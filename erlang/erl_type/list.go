package erl_type

import (
	"cfg_exporter/entities/base_type"
	"fmt"
)

type ErlList struct {
	*base_type.List
}

func init() {
	base_type.Register("list", newList)
}

func newList(args []string) (any, error) {

	if args == nil {
		return &ErlList{List: &base_type.List{}}, nil
	}

	if len(args) == 1 {
		t, err := base_type.New(args[0])
		if err != nil {
			return nil, err
		}
		return &ErlList{List: &base_type.List{ElementType: t}}, nil
	}

	return nil, fmt.Errorf("类型格式错误 list|list(ElementType)")

}

func (l *ErlList) Convert(val any) string {
	return toString(val)
}

func (l *ErlList) String() string {
	return "list()"
}

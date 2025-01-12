package typesystem

import (
	"cfg_exporter/entities/typesystem"
)

type ErlList struct {
	*typesystem.List
}

func init() {
	typesystem.Register("list", newList)
}

func newList(typeStr string) (any, error) {
	list, err := typesystem.NewList(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlList{List: list}, nil
}

func (l *ErlList) Convert(val any) string {
	return toString(val)
}

func (l *ErlList) String() string {
	return "list()"
}

package typesystem

import (
	"cfg_exporter/entities"
	"cfg_exporter/erlang"
)

type ErlList struct {
	*entities.List
}

func init() {
	entities.TypeRegister("list", newList)
}

func newList(typeStr string) (entities.ITypeSystem, error) {
	list, err := entities.NewList(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlList{List: list}, nil
}

func (l *ErlList) Convert(val any) string {
	return erlang.toString(val)
}

func (l *ErlList) String() string {
	return "list()"
}

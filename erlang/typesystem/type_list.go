package typesystem

import (
	"cfg_exporter/entities"
)

type ErlList struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("list", newList)
}

func newList(typeStr string) (entities.ITypeSystem, error) {
	list, err := entities.NewList(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlList{ITypeSystem: list}, nil
}

func (l *ErlList) Convert(val any) string {
	return toString(val)
}

func (l *ErlList) String() string {
	return "list()"
}

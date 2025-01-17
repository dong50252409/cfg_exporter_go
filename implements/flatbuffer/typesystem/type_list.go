package typesystem

import (
	"cfg_exporter/entities"
)

type FBList struct {
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
	return &FBList{ITypeSystem: list}, nil
}

func (l *FBList) String() string {
	return "[]"
}

func (*FBList) GetDefaultValue() string {
	return "[]"
}

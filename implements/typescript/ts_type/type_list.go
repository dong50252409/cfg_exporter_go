package ts_type

import (
	"cfg_exporter/entities"
)

type TSList struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("list", newList)
}

func newList(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	list, err := entities.NewList(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &TSList{ITypeSystem: list}, nil
}

func (l *TSList) Convert(val any) string {
	return toString(val)
}

func (l *TSList) String() string {
	return "any[]"
}

func (*TSList) GetDefaultValue() string {
	return "[]"
}

package typesystem

import (
	"cfg_exporter/entities"
)

type JSONList struct {
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
	return &JSONList{List: list}, nil
}

func (l *JSONList) Convert(val any) string {
	return toString(val)
}

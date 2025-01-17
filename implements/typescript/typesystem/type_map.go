package typesystem

import (
	"cfg_exporter/entities"
)

type TSMap struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("map", newMap)
}

func newMap(typeStr string) (entities.ITypeSystem, error) {
	mapType, err := entities.NewMap(typeStr)
	if err != nil {
		return nil, err
	}
	return &TSMap{ITypeSystem: mapType}, nil
}

func (*TSMap) Convert(val any) string {
	return toString(val)
}

func (m *TSMap) String() string {
	return "Map<any, any>"
}

func (*TSMap) GetDefaultValue() string {
	return "new Map()"
}

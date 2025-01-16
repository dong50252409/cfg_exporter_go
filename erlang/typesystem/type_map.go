package typesystem

import (
	"cfg_exporter/entities"
)

type ErlMap struct {
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
	return &ErlMap{ITypeSystem: mapType}, nil
}

func (*ErlMap) Convert(val any) string {
	return toString(val)
}

func (m *ErlMap) String() string {
	return "map()"
}

func (*ErlMap) GetDefaultValue() string {
	return "#{}"
}

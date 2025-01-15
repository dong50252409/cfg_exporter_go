package typesystem

import (
	"cfg_exporter/entities"
	"cfg_exporter/erlang"
)

type ErlMap struct {
	*entities.Map
}

func init() {
	entities.TypeRegister("map", newMap)
}

func newMap(typeStr string) (entities.ITypeSystem, error) {
	mapType, err := entities.NewMap(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlMap{Map: mapType}, nil
}

func (*ErlMap) Convert(val any) string {
	return erlang.toString(val)
}

func (m *ErlMap) String() string {
	return "map()"
}

func (*ErlMap) GetDefaultValue() string {
	return "#{}"
}

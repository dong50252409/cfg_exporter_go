package typesystem

import (
	"cfg_exporter/entities"
)

type JSONMap struct {
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
	return &JSONMap{Map: mapType}, nil
}

func (*JSONMap) Convert(val any) string {
	return toString(val)
}

package typesystem

import (
	"cfg_exporter/entities"
)

type FBMap struct {
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
	return &FBMap{ITypeSystem: mapType}, nil
}

func (m *FBMap) String() string {
	return "[]"
}

func (*FBMap) GetDefaultValue() string {
	return "[]"
}

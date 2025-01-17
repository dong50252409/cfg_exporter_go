package typesystem

import (
	"cfg_exporter/entities"
)

type FBRaw struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("str", newRaw)
}

func newRaw(typeStr string) (entities.ITypeSystem, error) {
	raw, err := entities.NewRaw(typeStr)
	if err != nil {
		return nil, err
	}
	return &FBRaw{ITypeSystem: raw}, nil
}

func (s *FBRaw) String() string {
	return "string"
}

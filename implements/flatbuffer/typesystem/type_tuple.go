package typesystem

import (
	"cfg_exporter/entities"
)

type FBTuple struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("tuple", newTuple)
}

func newTuple(typeStr string) (entities.ITypeSystem, error) {
	tuple, err := entities.NewTuple(typeStr)
	if err != nil {
		return nil, err
	}
	return &FBTuple{ITypeSystem: tuple}, nil
}

func (*FBTuple) String() string {
	return "[]"
}

func (*FBTuple) GetDefaultValue() string {
	return "[]"
}

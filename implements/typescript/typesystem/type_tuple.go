package typesystem

import (
	"cfg_exporter/entities"
)

type TSTuple struct {
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
	return &TSTuple{ITypeSystem: tuple}, nil
}

func (*TSTuple) Convert(val any) string {
	return toString(val)
}

func (*TSTuple) String() string {
	return "any[]"
}

func (*TSTuple) GetDefaultValue() string {
	return "[]"
}

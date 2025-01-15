package typesystem

import (
	"cfg_exporter/entities"
)

type ErlTuple struct {
	*entities.Tuple
}

func init() {
	entities.TypeRegister("tuple", newTuple)
}

func newTuple(typeStr string) (entities.ITypeSystem, error) {
	tuple, err := entities.NewTuple(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlTuple{Tuple: tuple}, nil
}

func (*ErlTuple) Convert(val any) string {
	return toString(val)
}

func (*ErlTuple) String() string {
	return "tuple()"
}

func (*ErlTuple) GetDefaultValue() string {
	return "{}"
}

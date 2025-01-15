package typesystem

import (
	"cfg_exporter/entities"
)

type JSONTuple struct {
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
	return &JSONTuple{Tuple: tuple}, nil
}

func (*JSONTuple) Convert(val any) string {
	return toString(val)
}

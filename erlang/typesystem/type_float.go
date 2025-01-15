package typesystem

import (
	"cfg_exporter/entities"
)

type ErlFloat struct {
	*entities.Float
}

func init() {
	entities.TypeRegister("float", newFloat)
}

func newFloat(typeStr string) (entities.ITypeSystem, error) {
	float, err := entities.NewFloat(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlFloat{Float: float}, nil
}

func (f *ErlFloat) String() string {
	return "float()"
}

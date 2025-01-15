package typesystem

import (
	"cfg_exporter/entities"
)

type JSONFloat struct {
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
	return &JSONFloat{Float: float}, nil
}

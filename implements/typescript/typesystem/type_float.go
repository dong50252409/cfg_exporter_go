package typesystem

import (
	"cfg_exporter/entities"
)

type TSFloat struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("float", newFloat)
}

func newFloat(typeStr string) (entities.ITypeSystem, error) {
	float, err := entities.NewFloat(typeStr)
	if err != nil {
		return nil, err
	}
	return &TSFloat{ITypeSystem: float}, nil
}

func (f *TSFloat) String() string {
	return "number"
}

package typesystem

import (
	"cfg_exporter/entities"
)

type FBFloat struct {
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
	return &FBFloat{ITypeSystem: float}, nil
}

func (f *FBFloat) String() string {
	return "float64"
}

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

func newFloat(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	float, err := entities.NewFloat(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &FBFloat{ITypeSystem: float}, nil
}

func (f *FBFloat) String() string {
	switch f.ITypeSystem.(*entities.Float).BitSize {
	case 32:
		return "float32"
	case 64:
		return "float64"
	default:
		return "float64"
	}
}

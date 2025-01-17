package typesystem

import (
	"cfg_exporter/entities"
)

type FBInteger struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("int", newInteger)
}

func newInteger(typeStr string) (entities.ITypeSystem, error) {
	integer, err := entities.NewInteger(typeStr)
	if err != nil {
		return nil, err
	}
	return &FBInteger{ITypeSystem: integer}, nil
}

func (i *FBInteger) String() string {
	return "float64"
}

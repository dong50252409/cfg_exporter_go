package typesystem

import (
	"cfg_exporter/entities"
)

type TSInteger struct {
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
	return &TSInteger{ITypeSystem: integer}, nil
}

func (i *TSInteger) String() string {
	return "number"
}

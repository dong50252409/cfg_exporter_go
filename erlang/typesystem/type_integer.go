package typesystem

import (
	"cfg_exporter/entities"
)

type ErlInteger struct {
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
	return &ErlInteger{ITypeSystem: integer}, nil
}

func (i *ErlInteger) String() string {
	return "integer()"
}

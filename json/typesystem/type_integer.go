package typesystem

import (
	"cfg_exporter/entities"
)

type JSONInteger struct {
	*entities.Integer
}

func init() {
	entities.TypeRegister("int", newInteger)
}

func newInteger(typeStr string) (entities.ITypeSystem, error) {
	integer, err := entities.NewInteger(typeStr)
	if err != nil {
		return nil, err
	}
	return &JSONInteger{Integer: integer}, nil
}

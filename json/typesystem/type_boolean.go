package typesystem

import (
	"cfg_exporter/entities"
)

type JSONBoolean struct {
	*entities.Boolean
}

func init() {
	entities.TypeRegister("bool", newBoolean)
}

func newBoolean(typeStr string) (entities.ITypeSystem, error) {
	boolean, err := entities.NewBoolean(typeStr)
	if err != nil {
		return nil, err
	}
	return &JSONBoolean{Boolean: boolean}, nil
}

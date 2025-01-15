package typesystem

import (
	"cfg_exporter/entities"
)

type JSONRaw struct {
	*entities.Raw
}

func init() {
	entities.TypeRegister("str", newRaw)
}

func newRaw(typeStr string) (entities.ITypeSystem, error) {
	newRaw, err := entities.NewRaw(typeStr)
	if err != nil {
		return nil, err
	}
	return &JSONRaw{Raw: newRaw}, nil
}

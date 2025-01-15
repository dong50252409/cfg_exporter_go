package typesystem

import (
	"cfg_exporter/entities"
)

type JSONStr struct {
	*entities.Str
}

func init() {
	entities.TypeRegister("str", newStr)
}

func newStr(typeStr string) (entities.ITypeSystem, error) {
	newStr, err := entities.NewStr(typeStr)
	if err != nil {
		return nil, err
	}
	return &JSONStr{Str: newStr}, nil
}

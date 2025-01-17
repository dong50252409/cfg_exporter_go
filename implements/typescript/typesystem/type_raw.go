package typesystem

import (
	"cfg_exporter/entities"
	"fmt"
)

type TSRaw struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("str", newRaw)
}

func newRaw(typeStr string) (entities.ITypeSystem, error) {
	raw, err := entities.NewRaw(typeStr)
	if err != nil {
		return nil, err
	}
	return &TSRaw{ITypeSystem: raw}, nil
}

func (s *TSRaw) Convert(val any) string {
	return fmt.Sprintf("'%s'", val)
}

func (s *TSRaw) String() string {
	return "string"
}

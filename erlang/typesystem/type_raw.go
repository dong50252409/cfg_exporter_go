package typesystem

import (
	"cfg_exporter/entities"
	"fmt"
)

type ErlRaw struct {
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
	return &ErlRaw{ITypeSystem: raw}, nil
}

func (s *ErlRaw) Convert(val any) string {
	return fmt.Sprintf("'%s'", val)
}

func (s *ErlRaw) String() string {
	return "atom()"
}

func (s *ErlRaw) GetDefaultValue() string {
	return "undefined"
}

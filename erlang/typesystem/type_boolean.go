package typesystem

import (
	"cfg_exporter/entities"
)

type ErlBoolean struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("bool", newBoolean)
}

func newBoolean(typeStr string) (entities.ITypeSystem, error) {
	boolean, err := entities.NewBoolean(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlBoolean{ITypeSystem: boolean}, nil
}

func (b *ErlBoolean) String() string {
	return "boolean()"
}

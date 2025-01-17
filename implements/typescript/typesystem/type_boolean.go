package typesystem

import (
	"cfg_exporter/entities"
)

type TSBoolean struct {
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
	return &TSBoolean{ITypeSystem: boolean}, nil
}

func (b *TSBoolean) String() string {
	return "boolean"
}

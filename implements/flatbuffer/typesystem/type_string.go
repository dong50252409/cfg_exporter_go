package typesystem

import (
	"cfg_exporter/entities"
)

type FBStr struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("str", newStr)
}

func newStr(typeStr string) (entities.ITypeSystem, error) {
	s, err := entities.NewStr(typeStr)
	if err != nil {
		return nil, err
	}
	return &FBStr{ITypeSystem: s}, nil
}

func (s *FBStr) String() string {
	return "string"
}

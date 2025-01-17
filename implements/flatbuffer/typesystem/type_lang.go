package typesystem

import (
	"cfg_exporter/entities"
)

type FBLang struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("lang", newLang)
}

func newLang(typeStr string) (entities.ITypeSystem, error) {
	lang, err := entities.NewLang(typeStr)
	if err != nil {
		return nil, err
	}
	return &FBLang{ITypeSystem: lang}, nil
}

func (l *FBLang) String() string {
	return "string"
}

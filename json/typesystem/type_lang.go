package typesystem

import (
	"cfg_exporter/entities"
	"fmt"
)

type JSONLang struct {
	*entities.Lang
}

func init() {
	entities.TypeRegister("lang", newLang)
}

func newLang(typeStr string) (entities.ITypeSystem, error) {
	lang, err := entities.NewLang(typeStr)
	if err != nil {
		return nil, err
	}
	return &JSONLang{Lang: lang}, nil
}

func (l *JSONLang) Convert(val any) string {
	return fmt.Sprintf("<<\"%s\"/utf8>>", val)
}

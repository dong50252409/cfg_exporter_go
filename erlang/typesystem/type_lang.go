package typesystem

import (
	"cfg_exporter/entities"
	"fmt"
)

type ErlLang struct {
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
	return &ErlLang{Lang: lang}, nil
}

func (l *ErlLang) Convert(val any) string {
	return fmt.Sprintf("<<\"%s\"/utf8>>", val)
}

func (l *ErlLang) String() string {
	return "binary()"
}

func (l *ErlLang) GetDefaultValue() string {
	return "<<>>"
}

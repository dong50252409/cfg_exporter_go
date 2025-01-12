package typesystem

import (
	"cfg_exporter/entities/typesystem"
)

type ErlLang struct {
	*typesystem.Lang
}

func init() {
	typesystem.Register("lang", newLang)
}

func newLang(typeStr string) (any, error) {
	lang, err := typesystem.NewLang(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlLang{Lang: lang}, nil
}

func (l *ErlLang) String() string {
	return "binary()"
}

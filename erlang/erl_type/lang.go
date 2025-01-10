package erl_type

import (
	"cfg_exporter/entities/base_type"
)

type ErlLang struct {
	*base_type.Lang
}

func init() {
	base_type.Register("lang", newLang)
}

func newLang(_ []string) (any, error) {
	return &ErlLang{Lang: &base_type.Lang{}}, nil
}

func (l *ErlLang) String() string {
	return "binary()"
}

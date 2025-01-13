package typesystem

import (
	"cfg_exporter/entities/typesystem"
	"fmt"
)

type ErlStr struct {
	*typesystem.Str
}

func init() {
	typesystem.Register("str", newString)
}

func newString(typeStr string) (any, error) {
	newStr, err := typesystem.NewStr(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlStr{Str: newStr}, nil
}

func (s *ErlStr) ParseString(typeStr string) (any, error) {
	return typeStr, nil
}

func (s *ErlStr) Convert(val any) string {
	return fmt.Sprintf("<<\"%s\"/utf8>>", val)
}

func (s *ErlStr) String() string {
	return "binary()"
}

func (s *ErlStr) GetDefaultValue() string {
	return "<<>>"
}

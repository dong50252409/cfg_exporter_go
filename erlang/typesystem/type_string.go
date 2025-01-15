package typesystem

import (
	"cfg_exporter/entities"
	"fmt"
)

type ErlStr struct {
	*entities.Str
}

func init() {
	entities.TypeRegister("str", newStr)
}

func newStr(typeStr string) (entities.ITypeSystem, error) {
	newStr, err := entities.NewStr(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlStr{Str: newStr}, nil
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

package entities

import (
	"fmt"
	"reflect"
)

type Str struct {
}

func init() {
	TypeRegister("str", NewStr)
}

func NewStr(_ string) (ITypeSystem, error) {
	return &Str{}, nil
}

func (s *Str) ParseString(str string) (any, error) {
	return str, nil
}

func (*Str) Convert(val any) string {
	return fmt.Sprintf(`"%v"`, val)
}

func (s *Str) String() string {
	return "str"
}

func (s *Str) GetDefaultValue() string {
	return `""`
}

func (s *Str) GetKind() reflect.Kind {
	return reflect.String
}

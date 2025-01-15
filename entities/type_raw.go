package entities

import (
	"fmt"
	"reflect"
)

type Raw struct {
}

func NewRaw(_ string) (*Raw, error) {
	return &Raw{}, nil
}

func (s *Raw) ParseString(str string) (any, error) {
	return str, nil
}

func (*Raw) Convert(val any) string {
	return fmt.Sprintf("%v", val)
}

func (s *Raw) String() string {
	return "raw"
}

func (s *Raw) GetDefaultValue() string {
	return "nil"
}

func (s *Raw) GetKind() reflect.Kind {
	return reflect.String
}

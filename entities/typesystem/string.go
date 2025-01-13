package typesystem

import (
	"fmt"
	"reflect"
)

type Str struct {
	DefaultValue string
}

func NewStr(_ string) (*Str, error) {
	return &Str{DefaultValue: ""}, nil
}

func (s *Str) ParseString(str string) (any, error) {
	return str, nil
}

func (*Str) Convert(val any) string {
	return fmt.Sprintf("%v", val)
}

func (s *Str) String() string {
	return "str"
}

func (s *Str) SetDefaultValue(val any) error {
	v, ok := val.(string)
	if ok {
		s.DefaultValue = v
		return nil
	}
	return fmt.Errorf("类型不匹配 %v", val)
}
func (s *Str) GetDefaultValue() string {
	return s.DefaultValue
}

func (s *Str) GetKind() reflect.Kind {
	return reflect.String
}

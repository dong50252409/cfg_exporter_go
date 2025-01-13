package typesystem

import (
	"fmt"
	"reflect"
)

type Lang struct {
	DefaultValue string
}

func NewLang(_ string) (*Lang, error) {
	return &Lang{DefaultValue: ""}, nil
}

func (l *Lang) ParseString(str string) (any, error) {
	return str, nil
}

func (l *Lang) Convert(val any) string {
	return val.(string)
}

func (l *Lang) String() string {
	return "lang"
}

func (l *Lang) SetDefaultValue(val any) error {
	v, ok := val.(string)
	if ok {
		l.DefaultValue = v
		return nil
	}
	return fmt.Errorf("类型不匹配 %v", val)
}

func (l *Lang) GetDefaultValue() string {
	return l.DefaultValue
}

func (l *Lang) GetKind() reflect.Kind {
	return reflect.String
}

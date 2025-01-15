package entities

import (
	"fmt"
	"reflect"
)

type Lang struct {
}

func NewLang(_ string) (*Lang, error) {
	return &Lang{}, nil
}

func (l *Lang) ParseString(str string) (any, error) {
	return str, nil
}

func (l *Lang) Convert(val any) string {
	return fmt.Sprintf(`"%v"`, val)
}

func (l *Lang) String() string {
	return "lang"
}

func (l *Lang) GetDefaultValue() string {
	return `""`
}

func (l *Lang) GetKind() reflect.Kind {
	return reflect.String
}

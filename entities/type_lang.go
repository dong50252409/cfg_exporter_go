package entities

import (
	"fmt"
	"reflect"
)

type Lang struct {
	Field *Field
}

func init() {
	TypeRegister("lang", NewLang)
}

func NewLang(_ string, field *Field) (ITypeSystem, error) {
	return &Lang{Field: field}, nil
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

func (l *Lang) GetCheckFunc() func(any) bool {
	return func(v any) bool {
		_, ok := v.(string)
		if !ok {
			_, ok = v.(AnyT)
			return ok
		}
		return ok
	}
}

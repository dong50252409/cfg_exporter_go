package base_type

import (
	"fmt"
)

type Lang struct {
}

func (l *Lang) ParseFromString(str string, _ ...any) (any, error) {
	return str, nil
}

func (l *Lang) Convert(val any, _ ...any) string {
	return val.(string)
}

func (l *Lang) String() string {
	return fmt.Sprintf("%T", l)
}

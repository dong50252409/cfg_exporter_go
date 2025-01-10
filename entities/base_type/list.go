package base_type

import (
	"fmt"
)

type List struct {
	ElementType any
}

func (l *List) ParseFromString(str string, _ ...any) (any, error) {
	return ParseString(str)
}

func (l *List) Convert(val any, _ ...any) string {
	return fmt.Sprintf("%v", val)
}

func (l *List) String() string {
	return fmt.Sprintf("%T", l)
}

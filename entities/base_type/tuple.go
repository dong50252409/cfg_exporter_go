package base_type

import (
	"fmt"
)

type Tuple struct {
	ElementType any
}

func (t *Tuple) ParseFromString(str string, _ ...any) (any, error) {
	return ParseString(str)
}

func (*Tuple) Convert(val ...any) string {
	return fmt.Sprintf("%v", val[0])
}

func (t *Tuple) String() string {
	return fmt.Sprintf("%T", t)
}

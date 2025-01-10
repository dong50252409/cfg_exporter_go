package base_type

import (
	"fmt"
)

type Str struct {
}

func (s *Str) ParseFromString(str string, _ ...any) (any, error) {
	return str, nil
}

func (*Str) Convert(val ...any) string {
	return fmt.Sprintf("%v", val[0])
}

func (s *Str) String() string {
	return fmt.Sprintf("%T", s)
}

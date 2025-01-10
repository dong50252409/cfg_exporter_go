package base_type

import (
	"fmt"
)

type Map struct {
	KeyT   any
	ValueT any
}

func (m *Map) ParseFromString(str string, _ ...any) (any, error) {
	return ParseString(str)
}

func (*Map) Convert(val ...any) string {
	return fmt.Sprintf("%v", val[0])
}

func (m *Map) String() string {
	return fmt.Sprintf("%T", m)
}

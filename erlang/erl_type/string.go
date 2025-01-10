package erl_type

import (
	"cfg_exporter/entities/base_type"
	"fmt"
)

type ErlStr struct {
	*base_type.Str
}

func init() {
	base_type.Register("str", newString)
}

func newString(_ []string) (any, error) {
	return &ErlStr{Str: &base_type.Str{}}, nil
}

func (s *ErlStr) Convert(val ...any) string {
	return fmt.Sprintf("<<\"%s\"/utf8>>", val[0])
}

func (s *ErlStr) String() string {
	return "binary()"
}

package erl_type

import (
	"cfg_exporter/entities/base_type"
)

type ErlBoolean struct {
	*base_type.Boolean
}

func init() {
	base_type.Register("bool", newBoolean)
}

func newBoolean(_ []string) (any, error) {
	t := &base_type.Boolean{}
	return &ErlBoolean{Boolean: t}, nil
}

func (b *ErlBoolean) String() string {
	return "boolean()"
}

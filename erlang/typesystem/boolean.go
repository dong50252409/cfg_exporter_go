package typesystem

import (
	"cfg_exporter/entities/typesystem"
)

type ErlBoolean struct {
	*typesystem.Boolean
}

func init() {
	typesystem.Register("bool", newBoolean)
}

func newBoolean(typeStr string) (any, error) {
	boolean, err := typesystem.NewBoolean(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlBoolean{Boolean: boolean}, nil
}

func (b *ErlBoolean) String() string {
	return "boolean()"
}

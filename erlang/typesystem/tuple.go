package typesystem

import (
	"cfg_exporter/entities/typesystem"
)

type ErlTuple struct {
	*typesystem.Tuple
}

func init() {
	typesystem.Register("tuple", newTuple)
}

func newTuple(typeStr string) (any, error) {
	tuple, err := typesystem.NewTuple(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlTuple{Tuple: tuple}, nil
}

func (*ErlTuple) Convert(val any) string {
	return toString(val)
}

func (*ErlTuple) String() string {
	return "tuple()"
}

func (*ErlTuple) GetDefaultValue() string {
	return "{}"
}

package typesystem

import (
	"cfg_exporter/entities/typesystem"
)

type ErlInteger struct {
	*typesystem.Integer
}

func init() {
	typesystem.Register("int", newInt)
}

func newInt(typeStr string) (any, error) {
	integer, err := typesystem.NewInteger(typeStr, "0")
	if err != nil {
		return nil, err
	}
	return &ErlInteger{Integer: integer}, nil
}

func (i *ErlInteger) String() string {
	return "integer()"
}

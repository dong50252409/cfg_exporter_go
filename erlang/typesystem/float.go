package typesystem

import (
	"cfg_exporter/entities/typesystem"
)

type ErlFloat struct {
	*typesystem.Float
}

func init() {
	typesystem.Register("float", newFloat)
}

func newFloat(typeStr string) (any, error) {
	float, err := typesystem.NewFloat(typeStr, "0.0")
	if err != nil {
		return nil, err
	}
	return &ErlFloat{Float: float}, nil
}

func (f *ErlFloat) String() string {
	return "float()"
}

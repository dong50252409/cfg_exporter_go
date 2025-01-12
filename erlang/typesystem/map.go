package typesystem

import (
	"cfg_exporter/entities/typesystem"
)

type ErlMap struct {
	*typesystem.Map
}

func init() {
	typesystem.Register("map", newMap)
}

func newMap(typeStr string) (any, error) {
	mapType, err := typesystem.NewMap(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlMap{Map: mapType}, nil
}

func (*ErlMap) Convert(val any) string {
	return toString(val)
}

func (m *ErlMap) String() string {
	return "map()"
}

func (*ErlMap) GetDefaultValue() string {
	return "#{}"
}

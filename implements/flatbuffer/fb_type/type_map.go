package fb_type

import (
	"cfg_exporter/entities"
	"fmt"
	"github.com/stoewer/go-strcase"
)

type FBMap struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("map", newMap)
}

func newMap(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	mapType, err := entities.NewMap(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &FBMap{ITypeSystem: mapType}, nil
}

func (m *FBMap) String() string {
	if m.ITypeSystem.(*entities.Map).ValueT == nil {
		return "string"
	}
	return fmt.Sprintf("[%sEntry]", strcase.UpperCamelCase(m.ITypeSystem.(*entities.Map).Field.Name))
}

func (*FBMap) GetDefaultValue() string {
	return "[]"
}

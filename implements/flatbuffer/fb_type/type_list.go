package fb_type

import (
	"cfg_exporter/entities"
	"fmt"
	"github.com/stoewer/go-strcase"
)

type FBList struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("list", newList)
}

func newList(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	list, err := entities.NewList(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &FBList{ITypeSystem: list}, nil
}

func (l *FBList) String() string {
	t := l.ITypeSystem.(*entities.List).T
	switch t.(type) {
	case *FBInteger:
		return fmt.Sprintf("[%s]", t.(*FBInteger).String())
	case *FBFloat:
		return fmt.Sprintf("[%s]", t.(*FBFloat).String())
	case *FBBoolean:
		return fmt.Sprintf("[%s]", t.(*FBBoolean).String())
	case *FBStr:
		return fmt.Sprintf("[%s]", t.(*FBStr).String())
	case *FBLang:
		return fmt.Sprintf("[%s]", t.(*FBLang).String())
	case *FBRaw:
		return fmt.Sprintf("[%s]", t.(*FBRaw).String())
	case *FBList, *FBTuple, *FBMap:
		return fmt.Sprintf("[%sEntry]", strcase.UpperCamelCase(l.ITypeSystem.(*entities.List).Field.Name))
	default:
		return "string"
	}
}

func (*FBList) GetDefaultValue() string {
	return "[]"
}

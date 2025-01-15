package entities

import (
	"cfg_exporter/util"
	"fmt"
)

var (
	typeRegistry = make(map[string]func(typeStr string) (ITypeSystem, error))
)

// TypeRegister 类型注册器
func TypeRegister(key string, cls func(typeStr string) (ITypeSystem, error)) {
	typeRegistry[key] = cls
}

func NewType(typeStr string) (ITypeSystem, error) {
	key, args := util.GetKey(typeStr)
	if cls, ok := typeRegistry[key]; ok {
		return cls(args)
	}
	return nil, fmt.Errorf("类型不存在 %s", key)
}

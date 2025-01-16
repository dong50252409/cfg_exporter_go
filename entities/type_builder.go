package entities

import (
	"cfg_exporter/util"
	"fmt"
)

var (
	typeRegistry     = make(map[string]func(typeStr string) (ITypeSystem, error))
	baseTypeRegistry = make(map[string]func(typeStr string) (ITypeSystem, error))
)

// TypeRegister 类型注册器
func TypeRegister(key string, cls func(typeStr string) (ITypeSystem, error)) {
	typeRegistry[key] = cls
	baseTypeRegistry[key] = cls
}

// RecoverBaseTypeRegister 恢复基础类型注册器
func RecoverBaseTypeRegister() {
	for k, v := range baseTypeRegistry {
		typeRegistry[k] = v
	}
}

// MergerTypeRegistry 合并当前的类型注册器
func MergerTypeRegistry(registry map[string]func(typeStr string) (ITypeSystem, error)) {
	for k, v := range registry {
		typeRegistry[k] = v
	}
}

func NewType(typeStr string) (ITypeSystem, error) {
	key, args := util.GetKey(typeStr)
	if cls, ok := typeRegistry[key]; ok {
		return cls(args)
	}
	return nil, fmt.Errorf("类型不存在 %s", key)
}

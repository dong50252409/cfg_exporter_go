package fb_type

import "cfg_exporter/entities"

var (
	typeRegistry = make(map[string]func(typeStr string, field *entities.Field) (entities.ITypeSystem, error))
)

// flatbuffer类型注册器
func typeRegister(key string, cls func(typeStr string, field *entities.Field) (entities.ITypeSystem, error)) {
	typeRegistry[key] = cls
}

// GetTypeRegister 获取类型注册器
func GetTypeRegister() map[string]func(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	return typeRegistry
}

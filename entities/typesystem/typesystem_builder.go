package typesystem

import (
	"fmt"
	"strings"
)

var (
	registry = make(map[string]func(typeStr string) (any, error))
)

// Register 类型注册器
func Register(key string, cls func(typeStr string) (any, error)) {
	registry[key] = cls
}

func New(typeStr string) (any, error) {
	key, args := getKey(typeStr)
	if cls, ok := registry[key]; ok {
		return cls(args)
	}
	return nil, fmt.Errorf("类型不存在 %s", key)
}

func getKey(typeStr string) (string, string) {
	index := strings.Index(typeStr, "(")
	if index != -1 {
		return typeStr[:index], typeStr[index:]
	} else {
		return typeStr, ""
	}
}

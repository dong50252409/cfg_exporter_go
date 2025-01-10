package base_type

import (
	"fmt"
	"strings"
)

var (
	registry = make(map[string]func(args []string) (any, error))
)

// Register 类型注册器
func Register(key string, cls func(args []string) (any, error)) {
	registry[key] = cls
}

func New(val string) (any, error) {
	var key string
	var args []string
	LIndex := strings.Index(val, "(")
	if LIndex == -1 {
		key = val
	} else {
		key = val[:LIndex]
		rIndex := strings.LastIndex(val, ")")
		args = strings.Split(val[LIndex+1:rIndex], ",")
	}

	if cls, ok := registry[key]; ok {
		return cls(args)
	}
	return nil, fmt.Errorf("类型不存在 %s", key)

}

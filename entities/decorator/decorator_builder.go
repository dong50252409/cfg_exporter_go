package decorator

import (
	"cfg_exporter/entities"
	"fmt"
	"strings"
)

var (
	registry = make(map[string]func(tbl *entities.Table, field *entities.Field, str string) error)
)

func register(key string, cls func(tbl *entities.Table, field *entities.Field, str string) error) {
	registry[key] = cls
}

func New(tbl *entities.Table, field *entities.Field, str string) error {
	key, args := getKey(str)
	cons, ok := registry[key]
	if !ok {
		return fmt.Errorf("装饰器不存在")
	}

	return cons(tbl, field, args)
}

func getKey(str string) (string, string) {
	index := strings.Index(str, "(")
	if index != -1 {
		return str[:index], str[index:]
	} else {
		return str, ""
	}
}

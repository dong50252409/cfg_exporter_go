package reader

import (
	"path/filepath"
)

var registry = make(map[string]IReader)

type IReader interface {
	CheckSupport(path string) bool
	Read(path string) ([][]string, error)
}

func Register(key string, cls IReader) {
	registry[key] = cls
}

func CheckSupport(path string) bool {
	ext := filepath.Ext(path)[1:]
	reader, ok := registry[ext]
	if !ok {
		return false
	}
	return reader.CheckSupport(path)
}

func Read(path string) ([][]string, error) {
	ext := filepath.Ext(path)[1:]
	return registry[ext].Read(path)
}

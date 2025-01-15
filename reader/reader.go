package reader

import (
	"path/filepath"
)

type IReader interface {
	CheckSupport(path string) bool
	Read(path string) ([][]string, error)
}

var registry = make(map[string]IReader)

// Register 注册文件读取器
func Register(key string, cls IReader) {
	registry[key] = cls
}

// CheckSupport 检查文件是否支持
func CheckSupport(path string) bool {
	ext := filepath.Ext(path)[1:]
	reader, ok := registry[ext]
	if !ok {
		return false
	}
	return reader.CheckSupport(path)
}

// Read 读取文件
func Read(path string) ([][]string, error) {
	ext := filepath.Ext(path)[1:]
	return registry[ext].Read(path)
}

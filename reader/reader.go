package reader

import (
	"fmt"
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
	fmt.Printf("读取配置文件：%s\n", path)
	ext := filepath.Ext(path)[1:]
	records, err := registry[ext].Read(path)
	if err != nil {
		return nil, err
	}

	if records == nil || len(records) == 0 {
		return nil, fmt.Errorf("没有发现可读取的sheet页签，请检查页签名是否正确！")
	}

	return records, nil
}

package reader

import (
	"cfg_exporter/config"
	"fmt"
	"path/filepath"
)

type IReader interface {
	Read() ([][]string, error)
}

type Reader struct {
	cls  IReader
	Path string
}

var registry = make(map[string]func(reader *Reader) IReader)

// Register 注册文件读取器
func Register(key string, cls func(reader *Reader) IReader) {
	registry[key] = cls
}

func NewReader(path string) (*Reader, error) {
	ext := filepath.Ext(path)[1:]
	cls, ok := registry[ext]
	if !ok {
		return nil, fmt.Errorf("配置表不支持！%s", path)
	}
	r := &Reader{Path: path}
	r.cls = cls(r)
	return r, nil
}

// Read 读取文件
func (r *Reader) Read() ([][]string, error) {
	fmt.Printf("读取配置文件：%s\n", r.Path)

	records, err := r.cls.Read()
	if err != nil {
		return nil, fmt.Errorf("配置表读取失败！%s %s", r.Path, err)
	}

	if records == nil || len(records) == 0 {
		return nil, fmt.Errorf("没有发现可读取的sheet页签，请检查页签名是否正确！%s", r.Path)
	}

	// 删除空行
	for index := config.Config.BodyStartRow - 1; index < len(records); {
		if records[index] == nil {
			records = append(records[:index], records[index+1:]...)
		} else {
			index++
		}
	}

	return records, nil
}

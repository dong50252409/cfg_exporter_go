package reader

import (
	"errors"
	"fmt"
)

var (
	ErrorTableTempFile     error
	ErrorTableNotSupported error
	ErrorTableReadFailed   error
	ErrorTableNotSheet     error
)

func errorTableTempFile(path string) error {
	ErrorTableTempFile = errors.New(fmt.Sprintf("无法读取临时文件！%s", path))
	return ErrorTableTempFile
}

func errorTableNotSupported(path string) error {
	ErrorTableNotSupported = errors.New(fmt.Sprintf("配置表不支持！%s", path))
	return ErrorTableNotSupported
}

func errorTableReadFailed(path string, err error) error {
	ErrorTableReadFailed = errors.New(fmt.Sprintf("配置表读取失败！%s 错误:%s", path, err))
	return ErrorTableReadFailed
}

func errorTableNotSheet(path string) error {
	ErrorTableNotSheet = errors.New(fmt.Sprintf("没有发现可读取的数据，请检查页签名命名是否正确！%s", path))
	return ErrorTableNotSheet
}

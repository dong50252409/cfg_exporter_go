package util

import (
	"fmt"
	"strings"
)

// SubTableName 获取表名
func SubTableName(filename string) (string, error) {
	index := strings.Index(filename, "(")
	lastIndex := strings.LastIndex(filename, ")")
	if index == -1 || lastIndex == -1 {
		return "", fmt.Errorf("文件名格式错误 配表描述(表名).ext")
	}
	return filename[strings.Index(filename, "(")+1 : strings.LastIndex(filename, ")")], nil
}

// GetKey 获取key以及未解析的参数
func GetKey(typeStr string) (string, string) {
	index := strings.Index(typeStr, "(")
	if index != -1 {
		return typeStr[:index], typeStr[index:]
	} else {
		return typeStr, ""
	}
}

func SubArgs(str string, sep string) []string {
	index := strings.Index(str, "(")
	if index != -1 {
		end := strings.LastIndex(str, ")")
		if end == -1 {
			return []string{}
		}
		return strings.Split(str[index+1:end], sep)
	}
	return []string{}
}

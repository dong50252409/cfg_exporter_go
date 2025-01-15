package util

import "strings"

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

package util

import "strings"

func SubArgs(str string, sep string) []string {
	index := strings.Index(str, "(")
	if index != -1 {
		end := strings.LastIndex(str, ")")
		if end == -1 {
			return nil
		}
		return strings.Split(str[index+1:end], sep)
	}
	return nil
}

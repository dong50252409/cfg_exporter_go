package entities

import (
	"strings"
	"text/template"
)

var FuncMap = template.FuncMap{
	"toUpper":     strings.ToUpper,
	"toLower":     strings.ToLower,
	"add":         Add,
	"seq":         Seq,
	"joinByComma": JoinByComma,
}

// Add 两个数相加
func Add(a, b int) int {
	return a + b
}

// Seq 给定长度生成一个初始化后的切片
func Seq(n int) []int {
	return make([]int, n)
}

// JoinByComma 将字符串切片使用逗号链接
func JoinByComma(items []string) string {
	return strings.Join(items, ", ")
}

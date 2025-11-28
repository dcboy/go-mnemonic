package main

import (
	_ "embed"
	"strings"
)

//go:embed wordlists/chinese_simplified.txt
var zhCNWordlistText string

var zhCNWordlist []string

func init() {
	// 将内嵌的词表文本按行拆分为2048个词
	zhCNWordlist = strings.Split(strings.TrimSpace(zhCNWordlistText), "\n")
}

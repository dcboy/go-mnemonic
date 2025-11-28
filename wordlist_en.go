package main

import (
    _ "embed"
    "strings"
)

//go:embed wordlists/english.txt
var enWordlistText string

var enWordlist []string

func init() {
    enWordlist = strings.Split(strings.TrimSpace(enWordlistText), "\n")
}


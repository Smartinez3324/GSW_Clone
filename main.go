package main

import (
	"fmt"
	"github.com/AarC10/GSW-V2/proc"
	"strings"
)

type Field struct {
	Type   string `json:"type"`
	Endian string `json:"endian"`
}

type Port struct {
	Fields map[string]Field `json:"fields"`
}

func printWithTabbing(level int, text string) {
	fmt.Print(strings.Repeat("\t", level))
	fmt.Println(text)
}

func main() {
	proc.ParseConfiguration("data/config/test.json")
}

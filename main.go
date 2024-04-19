package main

import (
	"LoongPanel/Panel/Files"
	"fmt"
	"strings"
)

func main() {
	s := fmt.Sprint(Files.Dir("/"))
	s = strings.Replace(s, "}", "}\n", -1)
	fmt.Println(s)
}

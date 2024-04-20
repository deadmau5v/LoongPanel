package main

import (
	"LoongPanel/Panel/Files"
	"fmt"
)

func main() {
	dir, _ := Files.Dir("C:/")
	for _, file := range dir {
		fmt.Println(
			"Name:", file.Name,
			"Ext:", file.Ext,
		)

	}
}

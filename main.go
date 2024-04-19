package main

import (
	"LoongPanel/Panel/Files"
	"fmt"
)

func main() {
	dir, _ := Files.Dir("/")
	for _, v := range dir {
		fmt.Println(
			"Name:", v.Name,
			"Size:", v.Size,
			"Time:", v.Time,
			"IsDir:", v.IsDir,
			"Mod:", v.Mod,
			"Path:", v.Path,
			"Ext:", v.Ext,
			"IsHidden:", v.IsHidden,
		)
	}
}

package main

import (
	"LoongPanel/Panel/Files"
	"fmt"
)

func main() {
	dir, _ := Files.Dir("/root/LoongPanel/test")
	for _, file := range dir {
		fmt.Println(
			"Name:", file.Name,
			"Size:", file.Size,
			"Time:", file.Time,
			"IsDir:", file.IsDir,
			"Mode:", file.Mode,
			"Path:", file.Path,
			"Ext:", file.Ext,
			"IsHidden:", file.IsHidden,
		)
		if !file.IsDir {
			newFile := Files.NewFile()
			newFile.Path = file.Path + " back"
			err := file.Copy(newFile)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}

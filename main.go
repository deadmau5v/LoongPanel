package main

import (
	"LoongPanel/Panel/Files"
	"fmt"
)

func main() {
	dir, _ := Files.Dir("C:\\Users\\d5v\\Desktop\\test")
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
		newFile := Files.NewFile()
		newFile.Path = file.Path + ".copy.txt"
		err := file.Copy(newFile)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

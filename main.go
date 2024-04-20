package main

import (
	"LoongPanel/Panel/Files"
	"fmt"
)

func main() {
	dir, _ := Files.Dir("D:\\code\\golang\\LoongPanel\\test")
	//	dir, _ := Files.Dir("/root/LoongPanel/test")
	for _, file := range dir {
		fmt.Println(
			"文件:", file.Name,
		)
		toFile := file.CopyObj()

		f := toFile.SetPath("D:\\code\\golang\\LoongPanel\\test\\test2")
		err := file.Move(f)
		if err != nil {
			fmt.Println("错误", err.Error())
			return
		}

	}
}

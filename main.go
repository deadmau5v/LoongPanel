/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：入口文件
 */

package main

import (
	"LoongPanel/Panel/API"
	"fmt"
)

func main() {
	fmt.Println("http://127.0.0.1:8080")
	err := API.App.Run()
	if err != nil {
		return
	}
}

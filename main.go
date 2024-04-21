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

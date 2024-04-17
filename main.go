package main

import (
	"LoongPanel/Panel/Status"
	"fmt"
	"time"
)

func main() {
	for {
		time.Sleep(1 * time.Second)
		fmt.Println(
			"负荷: ", Status.LoadAverage1m(),
			fmt.Sprintf(" CPU: %.2f", Status.CPUPercent()),
			fmt.Sprintf(" Memroy: %.2f", Status.MemroyPercent()),
		)

	}

}

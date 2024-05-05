package System

import (
	"fmt"
)

func init() {
	var err = error(nil)
	Data, err = GetOSData()
	temp, err := getPublicIP()
	if err != nil {
		fmt.Println("GetPublicIP() Error: ", err.Error())
	}
	PublicIP = temp

	// 开启线程 实时监控CPU占用 防止调用时阻塞
	go func() {
		for {
			CPUPercent = getCPUPercent()
		}
	}()
}

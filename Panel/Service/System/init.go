/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：初始化系统信息
 */

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

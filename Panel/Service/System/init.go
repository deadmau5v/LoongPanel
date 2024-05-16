/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：初始化系统信息
 */

package System

import (
	"LoongPanel/Panel/Service/Log"
	"os"
)

func init() {
	var err = error(nil)
	Data, err = GetOSData()
	temp, err := getPublicIP()
	if err != nil {
		Log.ERROR("GetPublicIP() Error: ", err.Error())
	}
	PublicIP = temp
	WORKDIR, err = os.Getwd()
	// 开启线程 实时监控CPU占用 防止调用时阻塞
	go func() {
		for {
			CPUPercent = getCPUPercent()
		}
	}()
}

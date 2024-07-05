/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：初始化系统信息
 */

package System

import (
	"LoongPanel/Panel/Service/PanelLog"
	"os"
	"sync"
)

func init() {
	var err error
	// 获取系统信息
	Data, err = GetOSData()
	if err != nil {
		PanelLog.ERROR("[系统管理] GetOSData() Error: ", err.Error())
	}
	// 获取公网IP
	PublicIP, err = getPublicIP()
	if err != nil {
		PanelLog.ERROR("[系统管理] GetPublicIP() Error: ", err.Error())
	}
	// 获取工作路径
	WORKDIR, err = os.Getwd()
	if err != nil {
		PanelLog.ERROR("[系统管理] Getwd() Error: ", err.Error())
	}
	// 多线程任务
	go func() {
		for {
			var wg sync.WaitGroup
			wg.Add(4)
			// 获取CPU使用率
			go func() { defer wg.Done(); CPUPercent = getCPUPercent() }()
			// 获取磁盘读取IO
			go func() { defer wg.Done(); diskReadIO() }()
			// 获取磁盘写入IO
			go func() { defer wg.Done(); diskWriteIO() }()
			// 获取网络IO
			go func() {
				defer wg.Done()
				err := networkIO()
				if err != nil {
					PanelLog.ERROR("[系统管理] NetworkIO() Error: ", err.Error())
					return
				}
			}()

			wg.Wait()
		}
	}()
}

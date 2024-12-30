/*
 * 创建人： deadmau5v
 * 创建时间： 2024-7-1
 * 文件作用：
 */

package Process

import (
	Log2 "LoongPanel/Panel/Service/PanelLog"

	"github.com/shirou/gopsutil/process"
)

// GetProcessesList 获取进程列表
func GetProcessesList() []*process.Process {
	processes, err := process.Processes()
	if err != nil {
		Log2.ERROR("获取进程列表失败", err)
		return nil
	}
	Log2.INFO("[进程管理] 获取进程列表")
	return processes
}

// KillProcess 结束进程
func KillProcess(pid int32) bool {
	processes := GetProcessesList()
	for _, p := range processes {
		if p.Pid == pid {
			proceessName, err := p.Name()
			if err != nil {
				Log2.ERROR("[进程管理] 获取进程名错误", err.Error())
			}
			Log2.INFO("[进程管理] 结束进程", proceessName, pid)
			err = p.Kill()
			if err != nil {
				Log2.ERROR("[进程管理] 结束进程错误", err.Error())
				return false
			}

			return true
		}
	}
	return false
}

// ProcessCount 获取进程数量
func ProcessCount() int {
	return len(GetProcessesList())
}

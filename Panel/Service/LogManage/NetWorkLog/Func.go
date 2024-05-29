/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-30
 * 文件作用：网络日志
 */

package NetWorkLog

import (
	Log2 "LoongPanel/Panel/Service/Log"
	"LoongPanel/Panel/Service/LogManage"
	"os/exec"
	"strconv"
)

func GetNetWorkLog() *LogManage.Log_ {
	log := LogManage.Log_{
		Ok:   true,
		Name: "网络日志",
	}
	log.GetLog = func(line int) []byte {
		output, err := exec.Command("journalctl", "-u", "NetworkManager", "-n", strconv.Itoa(line)).Output()
		if err != nil {
			log.Ok = false
			Log2.ERROR("[日志管理] 获取网络日志失败：" + err.Error())
			return nil
		}
		return output
	}

	log.ClearLog = func() {
		_, err := exec.Command("journalctl", "-u", "NetworkManager", "--rotate").Output()
		// Todo 未测试
		if err != nil {
			log.Ok = false
			Log2.ERROR("[日志管理] 清空网络日志失败：" + err.Error())
			return
		}
		return
	}

	if !log.Ok {
		return nil
	}

	return &log
}

func init() {
	LogManage.AddLog("网络日志", *GetNetWorkLog())
}

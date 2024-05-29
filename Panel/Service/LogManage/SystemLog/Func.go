/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-29
 * 文件作用：获取日志内容
 */

package SystemLog

import (
	Log2 "LoongPanel/Panel/Service/Log"
	"LoongPanel/Panel/Service/LogManage"
	"io"
	"os"
	"os/exec"
	"strings"
)

// GetLog 获取日志
func GetLog(log *LogManage.Log_, line int) []byte {
	if !log.Ok {
		return nil
	}

	Log2.DEBUG("获取日志", log.Name, log.Path)
	file, err := os.Open(log.Path)
	if err != nil {
		Log2.ERROR("打开日志文件失败", log.Name, log.Path)
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			Log2.ERROR("关闭日志文件失败", log.Name, log.Path)
		}
	}(file)
	all, err := io.ReadAll(file)
	if err != nil {
		Log2.ERROR("读取日志文件失败", log.Name, log.Path)
		return nil
	}

	if line == 0 {
		return all
	} else {
		allStr := string(all)
		allStrSplit := strings.Split(allStr, "\n")
		if len(allStrSplit) > line {
			allStrSplit = allStrSplit[len(allStrSplit)-line:]
			all = []byte(strings.Join(allStrSplit, "\n"))
		}
		return all
	}
}

// ClearLog 清空日志
func ClearLog(log *LogManage.Log_) {
	if !log.Ok {
		return
	}

	file, err := os.Open(log.Path)
	if err != nil {
		Log2.ERROR("[系统日志] 打开日志文件错误", err.Error())
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			Log2.ERROR("[系统日志] 关闭日志文件错误", err.Error())
		}
	}(file)
	// 截断文件
	err = file.Truncate(0)
	if err != nil {
		Log2.ERROR("[系统日志] 清空日志文件错误", err.Error())
	}
}

// creatLog 创建日志对象 简化流程
func createLog(path, name string, customGetLog func(log *LogManage.Log_, line int) []byte) *LogManage.Log_ {
	log := &LogManage.Log_{
		Path: path,
		Name: name,
	}
	log.Ok = log.CheckLogExist()

	if !log.Ok {
		return nil
	}

	log.GetLog = func(line int) []byte {
		if customGetLog != nil {
			return customGetLog(log, line)
		}
		return GetLog(log, line)
	}
	log.ClearLog = func() {
		ClearLog(log)
	}

	return log
}

// GetBootLog 获取启动日志
func GetBootLog() *LogManage.Log_ {
	return createLog("/var/log/boot.log", "系统启动日志", nil)
}

// GetKDumpLog 获取KDump日志
func GetKDumpLog() *LogManage.Log_ {
	return createLog("/var/log/kdump.log", "内核崩溃日志", nil)
}

// GetCronLog 获取定时任务日志
func GetCronLog() *LogManage.Log_ {
	return createLog("/var/log/cron.log", "定时任务日志", nil)
}

// GetFirewalldLog 获取Firewalld日志
func GetFirewalldLog() *LogManage.Log_ {
	return createLog("/var/log/firewalld", "防火墙日志", nil)
}

// GetMessagesLog 获取系统消息日志
func GetMessagesLog() *LogManage.Log_ {
	return createLog("/var/log/messages", "系统消息日志", nil)
}

// GetSecureLog 获取安全日志
func GetSecureLog() *LogManage.Log_ {
	return createLog("/var/log/secure", "安全日志", nil)
}

// GetWtmpLog 获取登录日志
func GetWtmpLog() *LogManage.Log_ {
	return createLog("/var/log/wtmp", "登录日志", func(log *LogManage.Log_, line int) []byte {
		output, err := exec.Command("utmpdump", log.Path).Output()
		if err != nil {
			Log2.ERROR("执行 utmpdump 命令失败", log.Name, log.Path)
			return nil
		}

		// 将输出转换为字符串并按行分割
		outputStr := string(output)
		outputStrSplit := strings.Split(outputStr, "\n")

		if len(outputStrSplit) > line {
			outputStrSplit = outputStrSplit[len(outputStrSplit)-line:]
		}
		return []byte(strings.Join(outputStrSplit, "\n"))
	})
}

// GetKernelLog 获取内核日志
func GetKernelLog() *LogManage.Log_ {
	log := createLog("", "内核日志", func(log *LogManage.Log_, line int) []byte {
		output, err := exec.Command("journalctl", "-k").Output()
		if err != nil {
			Log2.ERROR("执行 journalctl 命令失败", log.Name)
			return nil
		}

		// 将输出转换为字符串并按行分割
		outputStr := string(output)
		outputStrSplit := strings.Split(outputStr, "\n")

		if len(outputStrSplit) > line {
			outputStrSplit = outputStrSplit[len(outputStrSplit)-line:]
		}
		return []byte(strings.Join(outputStrSplit, "\n"))
	})
	log.ClearLog = func() {
		err := exec.Command("journalctl", "--rotate").Run()
		if err != nil {
			Log2.ERROR("执行 journalctl --rotate 命令失败", log.Name)
		}
	}
	return log
}

func init() {
	LogManage.AddLog("系统启动日志", *GetBootLog())
	LogManage.AddLog("内核崩溃日志", *GetKDumpLog())
	LogManage.AddLog("定时任务日志", *GetCronLog())
	LogManage.AddLog("防火墙日志", *GetFirewalldLog())
	LogManage.AddLog("系统消息日志", *GetMessagesLog())
	LogManage.AddLog("安全日志", *GetSecureLog())
	LogManage.AddLog("登录日志", *GetWtmpLog())
	LogManage.AddLog("内核日志", *GetKernelLog())
}

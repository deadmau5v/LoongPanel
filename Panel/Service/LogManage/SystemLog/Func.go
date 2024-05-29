/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-29
 * 文件作用：获取日志内容
 */

package SystemLog

import (
	Log2 "LoongPanel/Panel/Service/Log"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Log 日志结构体
type Log struct {
	Path     string // 日志文件路径
	Name     string // 日志名称
	Ok       bool   // 是否通过检查
	GetLog   func(line int) []byte
	ClearLog func()
}

// CheckLogExist 检查日志是否存在
func (log *Log) CheckLogExist() bool {
	file, err := os.Stat(log.Path)

	if err != nil {
		Log2.ERROR("获取日志文件信息失败", log.Path)
		return false
	}
	if file != nil && file.IsDir() {
		Log2.DEBUG("日志文件不存在或者是一个目录")
		return false
	}
	return true
}

// GetLog 获取日志
func GetLog(log *Log, line int) []byte {
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
func ClearLog(log *Log) {
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
func createLog(path, name string, customGetLog func(log *Log, line int) []byte) *Log {
	log := &Log{
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
func GetBootLog() *Log {
	return createLog("/var/log/boot.log", "系统启动日志", nil)
}

// GetKDumpLog 获取KDump日志
func GetKDumpLog() *Log {
	return createLog("/var/log/kdump.log", "内核崩溃日志", nil)
}

// GetCronLog 获取定时任务日志
func GetCronLog() *Log {
	return createLog("/var/log/cron.log", "定时任务日志", nil)
}

// GetFirewalldLog 获取Firewalld日志
func GetFirewalldLog() *Log {
	return createLog("/var/log/firewalld", "防火墙日志", nil)
}

// GetMessagesLog 获取系统消息日志
func GetMessagesLog() *Log {
	return createLog("/var/log/messages", "系统消息日志", nil)
}

// GetSecureLog 获取安全日志
func GetSecureLog() *Log {
	return createLog("/var/log/secure", "安全日志", nil)
}

// GetWtmpLog 获取登录日志
func GetWtmpLog() *Log {
	return createLog("/var/log/wtmp", "登录日志", func(log *Log, line int) []byte {
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

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

type Log struct {
	Path     string // 日志文件路径
	Name     string // 日志名称
	Ok       bool   // 是否通过检查
	GetLog   func(line int) []byte
	ClearLog func()
}

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

// GetBootLog 获取启动日志
func GetBootLog() *Log {
	log := &Log{
		Path: "/var/log/boot.log",
		Name: "系统启动日志",
	}
	log.Ok = log.CheckLogExist()

	if !log.Ok {
		return nil
	}

	log.GetLog = func(line int) []byte {
		return GetLog(log, line)
	}

	return log
}

// GetKDumpLog 获取KDump日志
func GetKDumpLog() *Log {
	log := &Log{
		Path: "/var/log/kdump.log",
		Name: "内核崩溃日志",
	}
	log.Ok = log.CheckLogExist()
	if !log.Ok {
		return nil
	}

	log.GetLog = func(line int) []byte {
		return GetLog(log, line)
	}
	log.ClearLog = func() {
		ClearLog(log)
	}

	return log
}

// GetCronLog 获取定时任务日志
func GetCronLog() *Log {
	log := &Log{
		Path: "/var/log/cron.log",
		Name: "定时任务日志",
	}
	log.Ok = log.CheckLogExist()
	if !log.Ok {
		return nil
	}

	log.GetLog = func(line int) []byte {
		return GetLog(log, line)
	}
	log.ClearLog = func() {
		ClearLog(log)
	}

	return log
}

// GetFirewalldLog 获取Firewalld日志
func GetFirewalldLog() *Log {
	log := &Log{
		Path: "/var/log/firewalld",
		Name: "防火墙日志",
	}
	log.Ok = log.CheckLogExist()
	if !log.Ok {
		return nil
	}

	log.GetLog = func(line int) []byte {
		return GetLog(log, line)
	}
	log.ClearLog = func() {
		ClearLog(log)
	}

	return log
}

// GetMessagesLog 获取系统消息日志
func GetMessagesLog() *Log {
	log := &Log{
		Path: "/var/log/messages",
		Name: "系统消息日志",
	}
	log.Ok = log.CheckLogExist()
	if !log.Ok {
		return nil
	}

	log.GetLog = func(line int) []byte {
		return GetLog(log, line)
	}
	log.ClearLog = func() {
		ClearLog(log)
	}

	return log
}

// GetSecureLog 获取安全日志
func GetSecureLog() *Log {
	log := &Log{
		Path: "/var/log/secure",
		Name: "安全日志",
	}
	log.Ok = log.CheckLogExist()
	if !log.Ok {
		return nil
	}

	log.GetLog = func(line int) []byte {
		return GetLog(log, line)
	}
	log.ClearLog = func() {
		ClearLog(log)
	}

	return log
}

// GetWtmpLog 获取登录日志
func GetWtmpLog() *Log {
	log := &Log{
		Path: "/var/log/wtmp",
		Name: "登录日志",
	}
	log.Ok = log.CheckLogExist()
	if !log.Ok {
		return nil
	}

	log.GetLog = func(line int) []byte {
		output, err := exec.Command("utmpdump", log.Path).Output()

		if err != nil {
			return nil
		}

		// 切分出最后几行
		outputStr := string(output)
		outputStrSplite := strings.Split(outputStr, "\n")
		if len(outputStrSplite) > line {
			outputStrSplite = outputStrSplite[len(outputStrSplite)-line:]
			output = []byte(strings.Join(outputStrSplite, "\n"))
		}

		return output
	}
	log.ClearLog = func() {
		ClearLog(log)
	}

	return log
}

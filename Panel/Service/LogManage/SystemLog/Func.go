/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-29
 * 文件作用：获取日志内容
 */

package SystemLog

import (
	Log2 "LoongPanel/Panel/Service/Log"
	"errors"
	"io"
	"os"
	"os/exec"
)

type Log struct {
	Path   string // 日志文件路径
	Name   string // 日志名称
	Err    error  // 错误信息
	GetLog func() []byte
}

func (log Log) CheckLogExist() Log {
	if log.Err != nil {
		return log
	}

	file, err := os.Stat(log.Path)
	log.Err = err
	if file == nil || file.IsDir() {
		Log2.DEBUG("日志文件不存在或者是一个目录")
		log.Err = errors.New("日志文件不存在或者是一个目录")
	}
	return log
}

func GetLog(log *Log) []byte {
	if log.Err != nil {
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
	return all
}

// GetBootLog 获取启动日志
func GetBootLog() *Log {
	log := &Log{
		Path: "/var/log/boot.log",
		Name: "系统启动日志",
	}
	log.CheckLogExist()
	if log.Err != nil {
		return nil
	}

	log.GetLog = func() []byte {
		return GetLog(log)
	}

	return log
}

// GetKDumpLog 获取KDump日志
func GetKDumpLog() *Log {
	log := &Log{
		Path: "/var/log/kdump.log",
		Name: "内核崩溃日志",
	}
	log.CheckLogExist()
	if log.Err != nil {
		return nil
	}

	log.GetLog = func() []byte {
		return GetLog(log)
	}

	return log
}

// GetCronLog 获取定时任务日志
func GetCronLog() *Log {
	log := &Log{
		Path: "/var/log/cron.log",
		Name: "定时任务日志",
	}
	log.CheckLogExist()
	if log.Err != nil {
		return nil
	}

	log.GetLog = func() []byte {
		return GetLog(log)
	}

	return log
}

// GetFirewalldLog 获取Firewalld日志
func GetFirewalldLog() *Log {
	log := &Log{
		Path: "/var/log/firewalld",
		Name: "防火墙日志",
	}
	log.CheckLogExist()
	if log.Err != nil {
		return nil
	}

	log.GetLog = func() []byte {
		return GetLog(log)
	}

	return log
}

// GetMessagesLog 获取系统消息日志
func GetMessagesLog() *Log {
	log := &Log{
		Path: "/var/log/messages",
		Name: "系统消息日志",
	}
	log.CheckLogExist()
	if log.Err != nil {
		return nil
	}

	log.GetLog = func() []byte {
		return GetLog(log)
	}

	return log
}

// GetSecureLog 获取安全日志
func GetSecureLog() *Log {
	log := &Log{
		Path: "/var/log/secure",
		Name: "安全日志",
	}
	log.CheckLogExist()
	if log.Err != nil {
		return nil
	}

	log.GetLog = func() []byte {
		return GetLog(log)
	}

	return log
}

// GetWtmpLog 获取登录日志
func GetWtmpLog() *Log {
	log := &Log{
		Path: "/var/log/wtmp",
		Name: "登录日志",
	}
	log.CheckLogExist()
	if log.Err != nil {
		return nil
	}

	log.GetLog = func() []byte {
		output, err := exec.Command("utmpdump", log.Path).Output()
		if err != nil {
			return nil
		}
		return output
	}

	return log
}

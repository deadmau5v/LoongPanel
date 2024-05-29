/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-29
 * 文件作用：包管理工具日志
 */

package PkgLog

import (
	Log2 "LoongPanel/Panel/Service/Log"
	"io"
	"os"
	"strings"
)

type Log struct {
	Path     string // 日志文件路径
	Name     string // 包管理名称
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

func createLog(path, name string, line int) *Log {
	log := Log{
		Path: path,
		Name: name,
		Ok:   true,
	}

	if !log.CheckLogExist() {
		return nil
	}

	log.GetLog = func(line int) []byte {
		if !log.Ok {
			Log2.DEBUG("[包管理日志] 错误跳过读取")
			return nil
		}
		logFile, err := os.Open(log.Path)
		if err != nil {
			Log2.ERROR("[包管理日志] 打开日志文件失败: ", err.Error())
			log.Ok = false
			return nil
		}
		defer logFile.Close()

		all, err := io.ReadAll(logFile)
		if err != nil {
			Log2.ERROR("[包管理日志] 读取日志文件失败: ", err.Error())
			log.Ok = false
			return nil
		}

		allStr := string(all)
		allStrSplite := strings.Split(allStr, "\n")

		if len(allStrSplite) > line {
			allStrSplite = allStrSplite[len(allStrSplite)-line:]
			all = []byte(strings.Join(allStrSplite, "\n"))
		}

		Log2.INFO("[包管理日志] 读取日志成功")
		return all
	}

	log.ClearLog = func() {
		if !log.Ok {
			Log2.DEBUG("[包管理日志] 错误跳过清除")
			return
		}
		file, err := os.Open(log.Path)
		if err != nil {
			Log2.ERROR("[包管理日志] 打开日志文件失败: ", err.Error())
			return
		}
		defer file.Close()

		err = file.Truncate(0)
		if err != nil {
			Log2.ERROR("[包管理日志] 截断日志文件失败: ", err.Error())
			return
		}
		Log2.INFO("[包管理日志] 清除日志成功")
	}

	return &log
}

// GetYumLog 获取yum日志
func GetYumLog(line int) *Log {
	return createLog("/var/log/yum.log", "yum", line)
}

// GetDnfLog 获取dnf日志
func GetDnfLog(line int) *Log {
	return createLog("/var/log/dnf.log", "dnf", line)
}

// GetAptLog 获取apt日志
func GetAptLog(line int) *Log {
	return createLog("/var/log/apt/history.log", "apt", line)
}

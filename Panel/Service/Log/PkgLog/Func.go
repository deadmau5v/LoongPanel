/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-29
 * 文件作用：包管理工具日志
 */

package PkgLog

import (
	"LoongPanel/Panel/Service/Log"
	Log2 "LoongPanel/Panel/Service/PanelLog"
	"io"
	"os"
	"strings"
)

func createLog(path, name string) *Log.Log_ {
	log := Log.Log_{
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
		file, err := os.Open(log.Path)
		if err != nil {
			Log2.ERROR("[包管理日志] 打开日志文件失败: ", err.Error())
			log.Ok = false
			return nil
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				Log2.ERROR("[包管理日志] 关闭日志文件失败: ", err.Error())
			}
		}(file)

		all, err := io.ReadAll(file)
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
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				Log2.ERROR("[包管理日志] 关闭日志文件失败: ", err.Error())
			}
		}(file)

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
func GetYumLog() *Log.Log_ {
	return createLog("/var/log/yum.log", "yum")
}

// GetDnfLog 获取dnf日志
func GetDnfLog() *Log.Log_ {
	return createLog("/var/log/dnf.log", "dnf")
}

// GetAptLog 获取apt日志
func GetAptLog() *Log.Log_ {
	return createLog("/var/log/apt/history.log", "apt")
}

func init() {

}

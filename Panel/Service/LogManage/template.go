/*
 * 创建人： deadmau5v
 * 创建时间： 2024-0-0
 * 文件作用：
 */

package LogManage

import (
	Log2 "LoongPanel/Panel/Service/Log"
	"os"
)

type Log interface {
	// GetLog 获取日志
	GetLog()
	// ClearLog 清空日志
	ClearLog()
}

type Log_ struct {
	Path     string // 日志文件路径
	Name     string // 日志名称
	Ok       bool   // 是否通过检查
	GetLog   func(line int) []byte
	ClearLog func()
}

// CheckLog_Exist 检查日志是否存在
func (Log_ *Log_) CheckLog_Exist() bool {
	file, err := os.Stat(Log_.Path)

	if err != nil {
		Log2.ERROR("获取日志文件信息失败", Log_.Path)
		return false
	}
	if file != nil && file.IsDir() {
		Log2.DEBUG("日志文件不存在或者是一个目录")
		return false
	}
	return true
}

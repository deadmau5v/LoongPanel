/*
 * 创建人： deadmau5v
 * 创建时间： 2024-0-0
 * 文件作用：
 */

package LogManage

type Log interface {
	// GetLog 获取日志
	GetLog()
	// ClearLog 清空日志
	ClearLog()
}

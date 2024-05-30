/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-30
 * 文件作用：日志管理函数
 */

package LogManage

import Log2 "LoongPanel/Panel/Service/Log"

func AddLog(Name string, fn func() *Log_) {
	log := fn()
	if log != nil {
		Log2.DEBUG("[日志管理] 添加日志支持", Name, log.Path)
		// Name仅作为答应提高可读性 使用log.Name作为键值 确保一致性
		AllLog[log.Name] = *log
	}
}

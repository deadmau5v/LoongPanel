/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-29
 * 文件作用：面板日志
 */

package PanelLog

import (
	"LoongPanel/Panel/Service/LogManage"
)

func GetPanelLog() *LogManage.Log_ {

	return &LogManage.Log_{
		Path: "panel.log",
		Name: "面板日志",
		Ok:   true,
	}
}

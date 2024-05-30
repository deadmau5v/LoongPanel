/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-29
 * 文件作用：面板日志
 */

package PanelLog

import (
	Log2 "LoongPanel/Panel/Service/Log"
	"LoongPanel/Panel/Service/LogManage"
	"io"
	"os"
	"strings"
)

func GetPanelLog() *LogManage.Log_ {

	log := &LogManage.Log_{
		Path: "panel.log",
		Name: "面板日志",
		Ok:   true,
	}

	log.Ok = log.CheckLogExist()
	if !log.Ok {
		return nil
	}

	log.GetLog = func(line int) []byte {
		file, err := os.Open(log.Path)
		if err != nil {
			Log2.ERROR("[面板日志] 打开日志文件失败: ", err.Error())
			return nil
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				Log2.ERROR("[面板日志] 关闭日志文件失败: ", err.Error())
			}
		}(file)

		all, err := io.ReadAll(file)
		if err != nil {
			Log2.ERROR("[面板日志] 读取日志文件失败: ", err.Error())
			return nil
		}

		allStr := string(all)
		allStrSplite := strings.Split(allStr, "\n")
		if len(allStrSplite) > line {
			allStrSplite = allStrSplite[len(allStrSplite)-line:]
			all = []byte(strings.Join(allStrSplite, "\n"))
		}

		return all
	}

	log.ClearLog = func() {
		if !log.Ok {
			return
		}
		err := os.Truncate(log.Path, 0)
		if err != nil {
			Log2.ERROR("[面板日志] 清空日志文件失败: ", err.Error())
		}
	}

	return log
}

func init() {

}

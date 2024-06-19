/*
 * 创建人： deadmau5v
 * 创建时间： 2024-6-18
 * 文件作用：应用程序 frp server 支持
 */

package FrpServer

import (
	"LoongPanel/Panel/Service/AppStore"
)

var App = AppStore.App{}

// getVersion 获取版本
func getVersion() (string, error) {

	return "", nil
}

// isInstall 是否安装
func isInstall() bool {

	return true
}

// isRunning 是否运行
func isRunning() bool {

	return true
}

// Install 安装
func Install() (bool, error) {

	return true, nil
}

func Init() {
	App.Name = "frp Server"
	App.Tags = []string{"网络工具", "端口转发"}
	App.Icon = "frps.png"
	App.Path = "/usr/bin/frps"

	App.Version = getVersion
	App.IsInstall = isInstall
	App.IsRunning = isRunning
	App.Install = Install

	AppStore.Apps = append(AppStore.Apps, App)
}

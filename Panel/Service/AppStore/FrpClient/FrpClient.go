/*
 * 创建人： deadmau5v
 * 创建时间： 2024-6-19
 * 文件作用：FRP Client 管理
 */

package FRPClient

import (
	"LoongPanel/Panel/Service/AppStore"
	"errors"
	"os"
	"os/exec"
	"strings"
)

var App = AppStore.App{}

// getVersion 获取版本
func getVersion() (string, error) {
	if isInstall := isInstall(); !isInstall {
		return "", nil
	}

	output, err := exec.Command("frpc", "-v").Output()
	if err != nil {
		return "", err
	}
	output = []byte(strings.TrimSpace(string(output)))
	return string(output), nil
}

// isInstall 是否安装
func isInstall() bool {
	f, err := os.Stat("/usr/bin/frpc")
	if err != nil || os.IsNotExist(err) || f.IsDir() {
		return false
	}
	return true
}

// isRunning 是否运行
func isRunning() bool {
	if isInstall := isInstall(); !isInstall {
		return false
	}

	output, err := exec.Command("pgrep", "-f", "frpc").Output()
	if err != nil || len(output) == 0 {
		return false
	}
	return true
}

// install 安装
func install() (bool, error) {
	if isInstall := isInstall(); isInstall {
		return false, errors.New("FRP client已安装")
	}
	err := AppStore.RunScript("install_frpc.sh")
	if err != nil {
		return false, err
	}

	return true, nil
}

// uninstall 卸载
func uninstall() (bool, error) {
	if isInstall := isInstall(); !isInstall {
		return false, errors.New("FRP client未安装")
	}
	err := AppStore.RunScript("remove_frpc.sh")
	if err != nil {
		return false, err
	}
	return true, nil
}

// start 启动
func start() (bool, error) {
	if isInstall := isInstall(); !isInstall {
		return false, errors.New("FRP client未安装")
	}
	if isRunning := isRunning(); isRunning {
		return false, errors.New("FRP client已启动")
	}
	cmd := exec.Command("frpc", "-c", "/etc/frp/frpc.ini")
	if err := cmd.Start(); err != nil {
		return false, err
	}
	return true, nil
}

// stop 停止
func stop() (bool, error) {
	if isInstall := isInstall(); !isInstall {
		return false, errors.New("FRP client未安装")
	}
	if isRunning := isRunning(); !isRunning {
		return false, errors.New("FRP client未启动")
	}
	cmd := exec.Command("pkill", "-f", "frpc")
	if err := cmd.Run(); err != nil {
		return false, err
	}
	return true, nil
}

func Init() {
	App.Name = "FRP Client"
	App.Tags = []string{"网络工具"}
	App.Icon = "frp.png"
	App.Path = "/usr/bin/frpc"

	App.Version = getVersion
	App.IsInstall = isInstall
	App.IsRunning = isRunning
	App.Install = install
	App.Uninstall = uninstall
	App.Start = start
	App.Stop = stop

	AppStore.Apps = append(AppStore.Apps, App)
}

/*
 * 创建人： deadmau5v
 * 创建时间： 2024-6-17
 * 文件作用：应用程序 Docker支持
 */

package Docker

import (
	"LoongPanel/Panel/Service/AppStore"
	"errors"
	"os/exec"
	"strings"
)

var App = AppStore.App{}

// getVersion 获取版本
func getVersion() (string, error) {
	if isInstall := isInstall(); !isInstall {
		return "", nil
	}

	output, err := exec.Command("docker", "version", "--format", "'{{.Server.Version}}'").Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// isInstall 是否安装
func isInstall() bool {
	_, err := exec.LookPath("docker")
	if err != nil {
		return false
	}
	return true
}

// isRunning 是否运行
func isRunning() bool {
	if isInstall := isInstall(); !isInstall {
		return false
	}

	output, err := exec.Command("docker", "info").Output()
	if err != nil || !strings.Contains(string(output), "Server Version") {
		return false
	}
	return true
}

// Install 安装
func Install() (bool, error) {
	if isInstall := isInstall(); isInstall {
		return false, errors.New("docker已安装")
	}
	err := AppStore.RunScript("install_docker.sh")
	if err != nil {
		return false, err
	}

	return true, nil
}

// Uninstall 卸载
func Uninstall() (bool, error) {
	if isInstall := isInstall(); !isInstall {
		return false, errors.New("docker未安装")
	}
	err := AppStore.RunScript("remove_docker.sh")
	if err != nil {
		return false, err
	}
	return true, nil
}

// Start 启动
func Start() (bool, error) {
	if isInstall := isInstall(); !isInstall {
		return false, errors.New("docker未安装")
	}
	if isRunning := isRunning(); isRunning {
		return false, errors.New("docker已启动")
	}
	cmd := exec.Command("systemctl", "start", "docker")
	if err := cmd.Run(); err != nil {
		return false, err
	}
	return true, nil
}

// Stop 停止
func Stop() (bool, error) {
	if isInstall := isInstall(); !isInstall {
		return false, errors.New("docker未安装")
	}
	if isRunning := isRunning(); !isRunning {
		return false, errors.New("docker未启动")
	}
	cmd := exec.Command("systemctl", "stop", "docker")
	if err := cmd.Run(); err != nil {
		return false, err
	}
	return true, nil
}

func Init() {
	App.Name = "Docker"
	App.Tags = []string{"运行环境"}
	App.Icon = "Docker.png"
	App.Path = "/usr/bin/docker"

	App.Version = getVersion
	App.IsInstall = isInstall
	App.IsRunning = isRunning
	App.Install = Install
	App.Uninstall = Uninstall
	App.Start = Start
	App.Stop = Stop

	AppStore.Apps = append(AppStore.Apps, App)
}

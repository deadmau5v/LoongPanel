/*
 * 创建人： deadmau5v
 * 创建时间： 2024-6-17
 * 文件作用：应用程序 Docker支持
 */

package Apps

import (
	"LoongPanel/Panel/Service/AppStore"
	"errors"
	"os/exec"
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

	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// Install 安装
func Install() (bool, error) {
	if isInstall := isInstall(); isInstall {
		return false, errors.New("Docker已安装")
	}

	return true, nil
}

func init() {
	App.Name = "Docker"
	App.Tags = []string{"运行环境"}
	App.Icon = "Docker.png"
	App.Path = "/usr/bin/docker"

	App.Version = getVersion
	App.IsInstall = isInstall
	App.IsRunning = isRunning
	App.Install = Install

	AppStore.Apps = append(AppStore.Apps, App)
}

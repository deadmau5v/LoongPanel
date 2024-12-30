/*
 * 创建人： deadmau5v
 * 创建时间： 2024-6-17
 * 文件作用：应用程序 Docker支持
 */

package Docker

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

	output, err := exec.Command("docker", "version", "--format", "'{{.Server.Version}}'").Output()
	if err != nil {
		return "", err
	}
	output = []byte(strings.Replace(string(output), "'", "", -1))
	return string(output), nil
}

// isInstall 是否安装
func isInstall() bool {
	f, err := os.Stat("/usr/bin/docker")
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

	output, err := exec.Command("docker", "info").Output()
	if err != nil || !strings.Contains(string(output), "Server Version") {
		return false
	}
	return true
}

// install 安装
func install() (bool, error) {
	if isInstall := isInstall(); isInstall {
		return false, errors.New("docker已安装")
	}
	err := AppStore.RunScript("install_docker.sh")
	if err != nil {
		return false, err
	}

	return true, nil
}

// uninstall 卸载
func uninstall() (bool, error) {
	if isInstall := isInstall(); !isInstall {
		return false, errors.New("docker未安装")
	}
	err := AppStore.RunScript("remove_docker.sh")
	if err != nil {
		return false, err
	}
	return true, nil
}

// start 启动
func start() (bool, error) {
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

// stop 停止
func stop() (bool, error) {
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
	App.Install = install
	App.Uninstall = uninstall
	App.Start = start
	App.Stop = stop

	AppStore.Apps = append(AppStore.Apps, App)
}

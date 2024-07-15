/*
 * 创建人： deadmau5v
 * 创建时间： 2024-0-0
 * 文件作用：应用商店 自动安装应用
 */

package AppStore

import (
	"LoongPanel/Panel/Service/PanelLog"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
)

func DownloadScripts() (bool, error) {
	url := "https://cdn1.d5v.cc/CDN/Project/LoongPanel/bin/scripts.zip"
	PanelLog.INFO("[应用商店] 开始下载脚本文件")
	const DistPath = "./scripts.zip"
	const DistDir = "./script"

	resp, err := http.Get(url)
	if err != nil {
		PanelLog.ERROR("[应用商店]", err.Error())
		return false, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			PanelLog.ERROR("[应用商店]", err.Error())
		}
	}(resp.Body)
	//创建文件
	distFile, err := os.Create(DistPath)
	if err != nil {
		PanelLog.ERROR("[应用商店]", err.Error())
		return false, err
	}

	defer func(File *os.File) {
		err := File.Close()
		if err != nil {
			PanelLog.ERROR("[应用商店]", err.Error())
		}
	}(distFile)
	//写入文件
	_, err = io.Copy(distFile, resp.Body)
	if err != nil {
		PanelLog.ERROR("[应用商店]", err.Error())
		return false, err
	}
	//解压文件
	err = exec.Command("unzip", DistPath, "-d", DistDir).Run()
	if err != nil {
		PanelLog.ERROR("[应用商店]", err.Error())
		return false, err
	}

	return true, nil
}

func RunScript(scriptName string) error {
	// 这个函数会造成阻塞
	cmd := exec.Command("bash", path.Join("./script", scriptName))
	err := cmd.Run()
	if err != nil {
		return err
	}
	return err
}

func init() {
	stat, err := os.Stat("./script")
	if err != nil || !stat.IsDir() {
		_, err := DownloadScripts()
		if err != nil {
			PanelLog.ERROR("[应用商店]", "下载脚本文件失败")
		}
	}
}

func FindApp(name string) *App {
	for _, app := range Apps {
		if name == app.Name {
			return &app
		}
	}
	return nil
}

// 统计已安装应用数量
func AppCount() (int, int) {
	count := 0
	for _, app := range Apps {
		if app.IsInstall() {
			count++
		}
	}
	return len(Apps), count
}

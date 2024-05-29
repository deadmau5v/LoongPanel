/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：入口文件
 */

package main

import (
	"LoongPanel/Panel/API"
	"LoongPanel/Panel/Service/Log"
	"io"
	"net/http"
	"os"
	"os/exec"
)

//region 下载前端文件

// downloadDist 下载前端文件
func downloadDist() {
	Log.INFO("开始下载前端文件")
	// 下载地址：https://cdn1.d5v.cc/CDN/Project/LoongPanel/bin/dist.zip
	const DistURL = "https://cdn1.d5v.cc/CDN/Project/LoongPanel/bin/dist.zip"
	const DistPath = "./dist.zip"
	const DistDir = "./dist"

	resp, err := http.Get(DistURL)
	if err != nil {
		Log.ERROR(err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Log.ERROR(err.Error())
		}
	}(resp.Body)
	//创建文件
	distFile, err := os.Create(DistPath)
	if err != nil {
		Log.ERROR(err.Error())
		return
	}
	defer func(File *os.File) {
		err := File.Close()
		if err != nil {
			Log.ERROR(err.Error())
		}
	}(distFile)
	//写入文件
	_, err = io.Copy(distFile, resp.Body)
	if err != nil {
		Log.ERROR(err.Error())
		return
	}
	//解压文件
	err = exec.Command("unzip", DistPath, "-d", DistDir).Run()
	if err != nil {
		Log.ERROR(err.Error())
		return
	}
}

//endregion

func main() {
	var run = true

	//region 入口

	if run {
		// 如果 "./dist" 不存在，则执行下载Dist函数
		if _, err := os.Stat("./dist"); os.IsNotExist(err) {
			downloadDist()
		}

		Log.INFO("[入口] http://127.0.0.1:8080")
		err := API.App.Run(":8080")
		if err != nil {
			return
		}
	}

	//endregion

	//region 测试

	//endregion
}

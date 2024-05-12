/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：入口文件
 */

package main

import (
	"LoongPanel/Panel/API"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
)

// downloadDist 下载前端文件
func downloadDist() {
	slog.Info("开始下载前端文件")
	// 下载地址：https://cdn1.d5v.cc/CDN/Project/LoongPanel/bin/dist.zip
	const DistURL = "https://cdn1.d5v.cc/CDN/Project/LoongPanel/bin/dist.zip"
	const DistPath = "./dist.zip"
	const DistDir = "./dist"

	resp, err := http.Get(DistURL)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(resp.Body)
	//创建文件
	distFile, err := os.Create(DistPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func(File *os.File) {
		err := File.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(distFile)
	//写入文件
	_, err = io.Copy(distFile, resp.Body)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	//解压文件
	err = exec.Command("unzip", DistPath, "-d", DistDir).Run()
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func main() {
	if _, err := os.Stat("./dist"); os.IsNotExist(err) {
		downloadDist()
	}

	fmt.Println("http://127.0.0.1:8080")
	err := API.App.Run()
	if err != nil {
		return
	}
}

/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：入口文件
 */

package main

import (
	"LoongPanel/Panel/API"
	"LoongPanel/Panel/Service/Log"
	"LoongPanel/Panel/Service/LogManage"
	"LoongPanel/Panel/Service/LogManage/NetWorkLog"
	"LoongPanel/Panel/Service/LogManage/PanelLog"
	"LoongPanel/Panel/Service/LogManage/PkgLog"
	"LoongPanel/Panel/Service/LogManage/SystemLog"
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
	var run = false

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

	LogManage.AllLog = make(map[string]LogManage.Log_)
	LogManage.AddLog("系统启动日志", SystemLog.GetBootLog)
	LogManage.AddLog("内核崩溃日志", SystemLog.GetKDumpLog)
	LogManage.AddLog("定时任务日志", SystemLog.GetCronLog)
	LogManage.AddLog("防火墙日志", SystemLog.GetFirewalldLog)
	LogManage.AddLog("系统消息日志", SystemLog.GetMessagesLog)
	LogManage.AddLog("安全日志", SystemLog.GetSecureLog)
	LogManage.AddLog("登录日志", SystemLog.GetWtmpLog)
	LogManage.AddLog("内核日志", SystemLog.GetKernelLog)
	LogManage.AddLog("yum包管理工具日志", PkgLog.GetYumLog)
	LogManage.AddLog("dnf包管理工具日志", PkgLog.GetDnfLog)
	LogManage.AddLog("apt包管理工具日志", PkgLog.GetAptLog)
	LogManage.AddLog("yum包管理工具日志", PanelLog.GetPanelLog)
	LogManage.AddLog("网络日志", NetWorkLog.GetNetWorkLog)

	for _, log := range LogManage.AllLog {
		Log.DEBUG(log.Name, log.Ok)
	}

	//endregion
}

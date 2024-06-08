/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：入口文件
 */

package main

import (
	"LoongPanel/Panel/API"
	"LoongPanel/Panel/Service/Log"
	"LoongPanel/Panel/Service/Log/DataBaseLog"
	"LoongPanel/Panel/Service/Log/NetWorkLog"
	PanelLog2 "LoongPanel/Panel/Service/Log/PanelLog"
	"LoongPanel/Panel/Service/Log/PkgLog"
	"LoongPanel/Panel/Service/Log/SystemLog"
	"LoongPanel/Panel/Service/PanelLog"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

//region 下载前端文件

// downloadDist 下载前端文件
func downloadDist() {
	PanelLog.INFO("开始下载前端文件")
	const DistURL = "https://cdn1.d5v.cc/CDN/Project/LoongPanel/bin/dist.zip"
	const DistPath = "./dist.zip"
	const DistDir = "./dist"

	resp, err := http.Get(DistURL)
	if err != nil {
		PanelLog.ERROR(err.Error())
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			PanelLog.ERROR(err.Error())
		}
	}(resp.Body)
	//创建文件
	distFile, err := os.Create(DistPath)
	if err != nil {
		PanelLog.ERROR(err.Error())
		return
	}

	defer func(File *os.File) {
		err := File.Close()
		if err != nil {
			PanelLog.ERROR(err.Error())
		}
	}(distFile)
	//写入文件
	_, err = io.Copy(distFile, resp.Body)
	if err != nil {
		PanelLog.ERROR(err.Error())
		return
	}
	//解压文件
	err = exec.Command("unzip", DistPath, "-d", DistDir).Run()
	if err != nil {
		PanelLog.ERROR(err.Error())
		return
	}
}

//endregion

func main() {
	//region 初始化日志

	Log.AllLog = make(map[string]Log.Log_)
	Log.Add("系统启动日志", SystemLog.GetBootLog)
	Log.Add("内核崩溃日志", SystemLog.GetKDumpLog)
	Log.Add("定时任务日志", SystemLog.GetCronLog)
	Log.Add("防火墙日志", SystemLog.GetFirewalldLog)
	Log.Add("系统消息日志", SystemLog.GetMessagesLog)
	Log.Add("安全日志", SystemLog.GetSecureLog)
	Log.Add("登录日志", SystemLog.GetWtmpLog)
	Log.Add("内核日志", SystemLog.GetKernelLog)
	Log.Add("yum包管理工具日志", PkgLog.GetYumLog)
	Log.Add("dnf包管理工具日志", PkgLog.GetDnfLog)
	Log.Add("apt包管理工具日志", PkgLog.GetAptLog)
	Log.Add("面板日志", PanelLog2.GetPanelLog)
	Log.Add("网络日志", NetWorkLog.GetNetWorkLog)
	Log.Add("数据库日志", DataBaseLog.GetDataBaseLog)
	for _, log := range Log.AllLog {
		status := "初始化失败"
		if log.Ok {
			status = "初始化成功"
		}
		PanelLog.DEBUG("[日志管理]", log.Name, status)
	}

	//endregion

	//region 入口

	// 如果 "./dist" 不存在，则执行下载Dist函数
	if _, err := os.Stat("./dist"); os.IsNotExist(err) {
		downloadDist()
	}

	port := flag.String("port", "8080", "端口")
	host := flag.String("host", "0.0.0.0", "监控主机地址, 0.0.0.0监控全部访问 127.0.0.1 监控本机访问")
	flag.Parse()
	PanelLog.INFO(fmt.Sprintf("[LoongPanel] http://%s:%s", *host, *port))
	err := API.App.Run(*host + ":" + *port)
	if err != nil {
		return
	}

	//endregion

	defer func() {
		PanelLog.INFO("[LoongPanel] 程序退出")
	}()
}

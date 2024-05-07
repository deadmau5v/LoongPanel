/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供首页的API接口 包含系统信息、磁盘信息、系统状态、系统负载等
 */

package home

import (
	System2 "LoongPanel/Panel/Service/System"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SystemInfo(ctx *gin.Context) {
	data := map[string]interface{}{
		"system_arch":          "",
		"system_public_ip":     "",
		"system_cpu_name":      "",
		"system_linux_version": "",
		"system_run_time":      "",
	}

	data["system_arch"] = System2.Data.OSArch
	data["system_public_ip"] = System2.PublicIP
	data["system_cpu_name"] = System2.Data.CPUName
	data["system_linux_version"] = System2.Data.LinuxVersion
	data["system_run_time"] = System2.GetRunTime()

	ctx.JSON(http.StatusOK, data)
}

func Disks(ctx *gin.Context) {
	data := map[string]interface{}{}
	data["disks"] = System2.Data.Disks
	ctx.JSON(http.StatusOK, data)
}

func SystemStatus(ctx *gin.Context) {
	var (
		DiskUsage   float32
		AverageLoad float32
		MemoryUsage float32
		CpuUsage    float32
	)

	DiskUsage = System2.GetDiskUsage()
	AverageLoad, err := System2.LoadAverage1m()
	if err != nil {
		AverageLoad = 0
	}
	MemoryUsage, err = System2.MemoryUsage()
	if err != nil {
		MemoryUsage = 0
	}
	CpuUsage = System2.GetCpuUsage()
	res := map[string]interface{}{
		"disk_usage":   DiskUsage,
		"average_load": AverageLoad,
		"memory_usage": MemoryUsage,
		"cpu_usage":    CpuUsage,
	}
	ctx.JSON(http.StatusOK, res)
}

func Shutdown(ctx *gin.Context) {
	System2.Shutdown()
	// 没必要返回数据 都关机了 返回个屁
}

func Reboot(ctx *gin.Context) {
	System2.Reboot()
	// 没必要返回数据 都重启了 返回个屁
}

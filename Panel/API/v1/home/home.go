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
	ctx.JSON(http.StatusOK, gin.H{
		"disks": System2.Data.Disks,
	})
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

	// 负载
	loads, err := System2.LoadAverage()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	load1m, load5m, load15m := loads[0], loads[1], loads[2]
	memyUsed, memyFree := System2.GetRAMUsedAndFree()

	ctx.JSON(http.StatusOK, gin.H{
		"disk_usage":   DiskUsage,
		"average_load": AverageLoad,
		"memory_usage": MemoryUsage,
		"cpu_usage":    CpuUsage,

		"load1m":  load1m,
		"load5m":  load5m,
		"load15m": load15m,

		"cpu_number": System2.Data.CPUNumber,
		"cpu_cores":  System2.Data.CPUCores,
		"cpu_mhz":    System2.Data.CPUMHz,
		"cpu_arch":   System2.Data.OSArch,

		"ram_total":     System2.Data.RAM,
		"ram_used_free": [2]uint64{memyUsed, memyFree},
		"ram_mhz":       System2.Data.RAMMHz,
		"ram_swap":      System2.Data.Swap,

		"disk_total": System2.Data.DiskTotal,
		"disks":      System2.Data.Disks,
	})
}

func Shutdown(ctx *gin.Context) {
	System2.Shutdown()
	ctx.Abort()
	// 没必要返回数据 都关机了 无需返回
}

func Reboot(ctx *gin.Context) {
	System2.Reboot()
	ctx.Abort()
	// 没必要返回数据 都重启了 无需返回
}

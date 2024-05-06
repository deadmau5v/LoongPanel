package home

import (
	"LoongPanel/Panel/System"
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

	data["system_arch"] = System.Data.OSArch
	data["system_public_ip"] = System.PublicIP
	data["system_cpu_name"] = System.Data.CPUName
	data["system_linux_version"] = System.Data.LinuxVersion
	data["system_run_time"] = System.GetRunTime()

	ctx.JSON(http.StatusOK, data)
}

func Disks(ctx *gin.Context) {
	data := map[string]interface{}{}
	data["disks"] = System.Data.Disks
	ctx.JSON(http.StatusOK, data)
}

func SystemStatus(ctx *gin.Context) {
	var (
		DiskUsage   float32
		AverageLoad float32
		MemoryUsage float32
		CpuUsage    float32
	)

	DiskUsage = System.GetDiskUsage()
	AverageLoad, err := System.LoadAverage1m()
	if err != nil {
		AverageLoad = 0
	}
	MemoryUsage, err = System.MemoryUsage()
	if err != nil {
		MemoryUsage = 0
	}
	CpuUsage = System.GetCpuUsage()
	res := map[string]interface{}{
		"disk_usage":   DiskUsage,
		"average_load": AverageLoad,
		"memory_usage": MemoryUsage,
		"cpu_usage":    CpuUsage,
	}
	ctx.JSON(http.StatusOK, res)
}

package API

import (
	FileService "LoongPanel/Panel/Files"
	"LoongPanel/Panel/System"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

func CPUPercent(ctx *gin.Context) {
	data := map[string]interface{}{}
	res, err := System.CPU()
	if err != nil {
		data["msg"] = "获取CPU使用率失败"
		data["status"] = -1
		data["percent"] = 0
		ctx.JSON(http.StatusInternalServerError, data)
		return
	}
	res = math.Round(res*100) / 100
	data["msg"] = ""
	data["status"] = 0
	data["percent"] = res
	ctx.JSON(http.StatusOK, data)
}

func MemoryPercent(ctx *gin.Context) {
	data := map[string]interface{}{}
	res, err := System.Memory()
	if err != nil {
		data["msg"] = "获取内存使用率失败"
		data["status"] = -1
		data["percent"] = 0
		ctx.JSON(http.StatusInternalServerError, data)
		return
	}
	data["msg"] = ""
	data["status"] = 0
	data["percent"] = res
	data["max"] = System.Data.RAM
	ctx.JSON(http.StatusOK, data)
}

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

func FileDir(ctx *gin.Context) {
	data := map[string]interface{}{}
	var err error
	data["files"], err = FileService.Dir("/")
	if err != nil {
		data["status"] = -1
		data["msg"] = err.Error()
	} else {
		data["status"] = 0
		data["msg"] = ""
	}

	ctx.JSON(http.StatusOK, data)
}

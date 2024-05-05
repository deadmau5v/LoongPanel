package API

import (
	FileService "LoongPanel/Panel/Files"
	"LoongPanel/Panel/System"
	"LoongPanel/Panel/Terminal"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func FileDir(ctx *gin.Context) {
	path := ctx.Query("path")
	if path == "" {
		path = "/"
	}
	data := map[string]interface{}{}
	var err error
	data["files"], err = FileService.Dir(path)
	if err != nil {
		data["status"] = -1
		data["msg"] = err.Error()
	} else {
		data["status"] = 0
		data["msg"] = ""
	}

	ctx.JSON(http.StatusOK, data)
}

func screenInput(ctx *gin.Context) {
	id := getIntQuery(ctx, "id")
	cmd := getQuery(ctx, "cmd")
	screen := Terminal.MainScreenManager.GetScreen(id)
	screen.Input(cmd)
	data := map[string]interface{}{
		"msg":    "ok",
		"status": 0,
	}
	ctx.JSON(200, data)
}

func screenCreate(ctx *gin.Context) {
	fmt.Println("INFO screenCreate")
	id := getIntQuery(ctx, "id")
	name := getQuery(ctx, "name")
	err := Terminal.MainScreenManager.Create(name, uint32(id))
	data := map[string]interface{}{}
	if err != nil {
		data["status"] = -1
		data["msg"] = err.Error()
		ctx.JSON(200, data)
		fmt.Println("创建Screen错误")
		return
	}
	data["status"] = 0
	data["msg"] = "ok"
	ctx.JSON(200, data)
}

func screenClose(ctx *gin.Context) {
	id := getIntQuery(ctx, "id")
	Terminal.MainScreenManager.Close(id)
	ctx.JSON(200, "ok")
}

func getQuery(ctx *gin.Context, key string) string {
	value := ctx.Query(key)
	if value == "" {
		data := map[string]interface{}{
			"msg":    "无法获取参数: " + key,
			"status": -1,
		}
		ctx.JSON(200, data)

	}
	return value
}

func getIntQuery(ctx *gin.Context, key string) int {
	value, err := strconv.Atoi(ctx.Query(key))
	if err != nil {
		data := map[string]interface{}{
			"msg":    "参数无效 需要Int: " + key,
			"status": -1,
		}
		ctx.JSON(http.StatusInternalServerError, data)
	}
	return value
}

func getScreens(ctx *gin.Context) {
	data := make([]map[string]interface{}, 0)
	for _, v := range Terminal.MainScreenManager.Screens {
		data1 := map[string]interface{}{}
		data1["id"] = v.Id
		data1["name"] = v.Name

		data = append(data, data1)
	}
	ctx.JSON(200, data)
}

func screenOutput(ctx *gin.Context) {
	id := getIntQuery(ctx, "id")
	idx := getIntQuery(ctx, "idx")
	screen := Terminal.MainScreenManager.GetScreen(id)
	if screen == nil {
		data := map[string]interface{}{
			"msg":    "无法查询到ID",
			"status": -1,
		}
		ctx.JSON(200, data)
		return
	}

	data := map[string]interface{}{
		"msg":    "ok",
		"status": 0,
		"data":   screen.GetOutput()[idx:],
	}
	ctx.JSON(200, data)
}

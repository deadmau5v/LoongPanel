package API

import (
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
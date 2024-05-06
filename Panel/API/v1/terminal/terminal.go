package terminal

import (
	"LoongPanel/Panel/Terminal"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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

func GetScreens(ctx *gin.Context) {
	data := make([]map[string]interface{}, 0)
	for _, v := range Terminal.MainScreenManager.Screens {
		data1 := map[string]interface{}{}
		data1["id"] = v.Id
		data1["name"] = v.Name

		data = append(data, data1)
	}
	ctx.JSON(200, data)
}

func ScreenOutput(ctx *gin.Context) {
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

func ScreenClose(ctx *gin.Context) {
	id := getIntQuery(ctx, "id")
	Terminal.MainScreenManager.Close(id)
	ctx.JSON(200, "ok")
}

func ScreenCreate(ctx *gin.Context) {
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

func ScreenInput(ctx *gin.Context) {
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

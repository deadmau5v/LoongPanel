/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：终端相关API 主要实现在 Service/Terminal 包中
 */

package terminal

import (
	"LoongPanel/Panel/Service/Log"
	"LoongPanel/Panel/Service/Terminal"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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

func ScreenClose(ctx *gin.Context) {
	id := getIntQuery(ctx, "id")
	Terminal.MainScreenManager.Close(id)
	ctx.JSON(200, "ok")
}

func ScreenCreate(ctx *gin.Context) {
	Log.INFO("screenCreate")
	idStr := ctx.Query("id")
	var id int
	if idStr == "" {
		id = Terminal.GetNextId() // 时间戳
	} else {
		_id, err := strconv.Atoi(idStr)
		if err != nil {
			data := map[string]interface{}{
				"msg":    "无法获取参数: id",
				"status": -1,
			}
			ctx.JSON(200, data)
			return
		}
		id = _id
	}

	name := strconv.Itoa(id)
	err := Terminal.MainScreenManager.Create(name, uint32(id))
	data := map[string]interface{}{}
	if err != nil {
		data["status"] = -1
		data["msg"] = err.Error()
		ctx.JSON(200, data)
		Log.ERROR("创建Screen错误")
		return
	}
	data["status"] = 0
	data["msg"] = "ok"
	ctx.JSON(200, data)
}

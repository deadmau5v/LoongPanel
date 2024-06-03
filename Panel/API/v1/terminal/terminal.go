/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：终端相关API 主要实现在 Service/Terminal 包中
 */

package terminal

import (
	"LoongPanel/Panel/Service/PanelLog"
	"LoongPanel/Panel/Service/Terminal"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getIntQuery(ctx *gin.Context, key string) int {
	value, err := strconv.Atoi(ctx.Query(key))
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "参数无效 需要Int: " + key,
			"status": -1,
		})
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
	PanelLog.INFO("[网页终端] 创建终端")
	idStr := ctx.Query("id")
	var id int
	if idStr == "" {
		id = Terminal.GetNextId() // 时间戳
	} else {
		_id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(200, gin.H{
				"msg":    "无法获取参数: id",
				"status": -1,
			})
			return
		}
		id = _id
	}

	name := strconv.Itoa(id)
	err := Terminal.MainScreenManager.Create(name, uint32(id))
	if err != nil {

		ctx.JSON(200, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		PanelLog.ERROR("创建Screen错误")
		return
	}

	ctx.JSON(200, gin.H{
		"status": 0,
		"msg":    "ok",
	})
}

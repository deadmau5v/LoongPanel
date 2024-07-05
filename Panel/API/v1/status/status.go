/*
 * 创建人： deadmau5v
 * 创建时间： 2024-7-4
 * 文件作用：状态监控 API
 */

package status

import (
	"LoongPanel/Panel/Service/PanelLog"
	"LoongPanel/Panel/Service/Status"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// SetStatusStepTime 设置状态监控时间间隔
func SetStatusStepTime(ctx *gin.Context) {
	stepTime := ctx.Query("time")

	if stepTime == "" {
		ctx.JSON(200, gin.H{
			"status": 1,
			"msg":    "time 不能为空",
		})
		return
	}

	number, err := strconv.Atoi(stepTime)
	if err != nil && number < 0 {
		ctx.JSON(200, gin.H{
			"status": 1,
			"msg":    "time 必须是大于等于0的数字",
		})
		return
	}

	// 修改时间间隔
	Status.SetStepTime(time.Duration(number) * time.Second)
}

// SetSaveTime 设置状态监控保存时间
func SetSaveTime(ctx *gin.Context) {
	saveTime := ctx.Query("time")

	if saveTime == "" {
		ctx.JSON(200, gin.H{
			"status": 1,
			"msg":    "time 不能为空",
		})
		return
	}

	number, err := strconv.Atoi(saveTime)
	if err != nil && number < 0 {
		ctx.JSON(200, gin.H{
			"status": 1,
			"msg":    "time 必须是大于等于0的数字",
		})
		return
	}

	// 修改保存时间
	Status.SetSaveTime(time.Duration(number) * time.Hour)
}

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// GetStatus 获取状态监控信息 是用WebSocket
func GetStatus(ctx *gin.Context) {
	PanelLog.INFO("[状态监控]", "WebSocket 连接成功")
	defer func() {
		PanelLog.INFO("[状态监控]", "WebSocket 连接关闭")
	}()

	conn, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		PanelLog.ERROR("[状态监控]", "WebSocket 连接失败", err.Error())
		ctx.JSON(500, gin.H{
			"status": 1,
			"msg":    "WebSocket 连接失败",
		})
		return
	}
	defer conn.Close()

	t := uint64(0)

	for {
		data := Status.GetStatus(t)
		if len(data) == 0 {
			time.Sleep(time.Second * 5)
			continue
		}
		conn.WriteJSON(data)
		t = uint64(time.Now().Unix())
		time.Sleep(time.Second * 5)
	}
}

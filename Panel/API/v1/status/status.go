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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// SetStatusConfig 设置状态监控时间间隔和保存时间
func SetStatusConfig(ctx *gin.Context) {
	var config struct {
		StepTimeValue int    `json:"stepTimeValue"`
		StepTimeUnit  string `json:"stepTimeUnit"`
		SaveTimeValue int    `json:"saveTimeValue"`
		SaveTimeUnit  string `json:"saveTimeUnit"`
	}

	if err := ctx.BindJSON(&config); err != nil {
		ctx.JSON(200, gin.H{
			"status": 1,
			"msg":    "参数绑定失败",
		})
		return
	}

	if config.StepTimeValue < 0 {
		ctx.JSON(200, gin.H{
			"status": 1,
			"msg":    "stepTimeValue 必须是大于等于0的数字",
		})
		return
	}

	if config.SaveTimeValue < 0 {
		ctx.JSON(200, gin.H{
			"status": 1,
			"msg":    "saveTimeValue 必须是大于等于0的数字",
		})
		return
	}

	// 设置时间间隔
	var stepTimeDuration time.Duration
	switch config.StepTimeUnit {
	case "second":
		stepTimeDuration = time.Duration(config.StepTimeValue) * time.Second
	case "minute":
		stepTimeDuration = time.Duration(config.StepTimeValue) * time.Minute
	default:
		ctx.JSON(200, gin.H{
			"status": 1,
			"msg":    "无效的 stepTimeUnit",
		})
		return
	}
	Status.SetStepTime(stepTimeDuration)

	// 设置保存时间
	var saveTimeDuration time.Duration
	switch config.SaveTimeUnit {
	case "second":
		saveTimeDuration = time.Duration(config.SaveTimeValue) * time.Second
	case "minute":
		saveTimeDuration = time.Duration(config.SaveTimeValue) * time.Minute
	case "hour":
		saveTimeDuration = time.Duration(config.SaveTimeValue) * time.Hour
	default:
		ctx.JSON(200, gin.H{
			"status": 1,
			"msg":    "无效的 saveTimeUnit",
		})
		return
	}
	Status.SetSaveTime(saveTimeDuration)

	ctx.JSON(200, gin.H{
		"status": 0,
		"msg":    "设置成功",
	})
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

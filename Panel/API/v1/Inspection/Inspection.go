package inspection

import (
	inspectionService "LoongPanel/Panel/Service/Inspection"
	"LoongPanel/Panel/Service/PanelLog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Check(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	PanelLog.INFO("[巡检] 连接成功")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		PanelLog.ERROR("[巡检] 连接失败", err)
		return
	}
	defer conn.Close()

	channel := inspectionService.Check()

	for msg := range channel {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			PanelLog.ERROR("[巡检] 写入失败", err)
			break
		}
	}
}

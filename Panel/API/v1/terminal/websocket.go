package terminal

import (
	"LoongPanel/Panel/Service/PanelLog"
	"LoongPanel/Panel/Service/Terminal"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ScreenWs(c *gin.Context) {
	PanelLog.INFO("[网页终端]", "ScreenWs WebSocket 连接")
	w := c.Writer
	r := c.Request
	PanelLog.DEBUG("[网页终端]", "升级连接")
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		PanelLog.ERROR("[网页终端] 无法打开 WebSocket 连接: ", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer func(conn *websocket.Conn) {
		PanelLog.DEBUG("[网页终端]", "WebSocket 连接关闭")
		err := conn.Close()
		if err != nil {
			PanelLog.ERROR("[网页终端]", "Ws 链接异常关闭: ", err)
		}
	}(conn)

	PanelLog.DEBUG("[网页终端]", "获取参数")
	id := getIntQuery(c, "id")
	screen := Terminal.MainScreenManager.GetScreen(uint32(id))
	if screen.WS != nil {
		err := screen.WS.Close()
		if err != nil {
			PanelLog.ERROR("[网页终端] 关闭原有连接失败: ", err)
			return
		}
	}
	screen.WS = conn
	PanelLog.DEBUG("[网页终端]", "创建新连接")
	input := make(chan []byte, 1024)
	output := screen.Subscribe()
	defer close(input)

	go func() {
		PanelLog.DEBUG("[网页终端]", "开始发送数据")
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				screen.Publish()
			}
		}
	}()

	go func() {
		PanelLog.DEBUG("[网页终端]", "开始读取数据")
		for {
			if conn != nil {
				_, message, err := conn.ReadMessage()
				if err != nil {
					PanelLog.ERROR("[网页终端] 读取消息失败: ", err)
					// 关闭连接
					conn.Close()
					return
				}
				input <- message
			}
		}
	}()

	for {
		if conn != nil {
			select {
			case o := <-output:
				err := conn.WriteMessage(websocket.TextMessage, o)
				if err != nil {
					PanelLog.ERROR("[网页终端] 发送数据失败: ", err)
					// 关闭连接
					conn.Close()
					return
				}
			case i := <-input:
				screen.InputByte(i)
			}
		}
	}
}

/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供 terminal 的流式传输
 */

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
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		PanelLog.ERROR("[网页终端] 无法打开websocket连接")
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer func(conn *websocket.Conn) {
		PanelLog.DEBUG("[网页终端]", "websocket 连接关闭")
		err := conn.Close()
		if err != nil {
			PanelLog.ERROR("[网页终端]", "Ws链接异常关闭")
		}
	}(conn)

	id := getIntQuery(c, "id")
	screen := Terminal.MainScreenManager.GetScreen(uint32(id))
	if screen.WS != nil {
		err := screen.WS.Close()
		if err != nil {
			PanelLog.DEBUG("[网页终端] 关闭原有连接失败")
			return
		}
		screen.WS = conn
	}
	input := make(chan []byte, 1024)
	output := screen.Subscribe()
	defer close(input)

	close_ := func() {
		PanelLog.DEBUG("[网页终端] websocket 关闭中...")
		id_, name_ := screen.Id, screen.Name
		if screen.WS != nil && screen.WS == conn {
			err := conn.Close()
			if err != nil {
				PanelLog.DEBUG("[网页终端] 关闭连接失败")
				return
			}
			screen.Close()
			conn = nil
			screen.WS = nil
		}

		Terminal.MainScreenManager.Close(id)
		// 重建screen

		_ = Terminal.MainScreenManager.Create(name_, id_)

		c.Abort()
	}

	go func() {
		for {
			screen.Publish()
			time.Sleep(100)
		}
	}()

	go func() {
		for {
			if conn != nil {
				_, message, err := conn.ReadMessage()
				if err != nil {
					close_()
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
				err := conn.WriteMessage(1, o)
				if err != nil {
					close_()
					return
				}
			case i := <-input:
				screen.InputByte(i)
			}
		}
	}
}

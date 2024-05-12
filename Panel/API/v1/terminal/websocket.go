/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供 terminal 的流式传输
 */

package terminal

import (
	"LoongPanel/Panel/Service/Terminal"
	"log/slog"
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
	slog.Info("ScreenWs创建")
	w := c.Writer
	r := c.Request
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("无法打开websocket连接")
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			slog.Error("Ws链接异常关闭")
		}
	}(conn)

	id := getIntQuery(c, "id")
	screen := Terminal.MainScreenManager.GetScreen(uint32(id))
	if screen.WS != nil {
		screen.WS.Close()
		screen.WS = conn
	}

	input := make(chan []byte, 1024)
	output := screen.Subscribe()
	defer close(input)

	close_ := func() {
		id_, name_ := screen.Id, screen.Name
		if screen.WS != nil && screen.WS == conn {
			conn.Close()
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

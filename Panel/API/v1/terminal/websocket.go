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

	screen := Terminal.DefaultScreen
	if screen.Connected {
		// 关闭之前的连接
		screen.Close()
		// 新建一个连接
		_ = Terminal.MainScreenManager.Create(screen.Name, screen.Id)
		screen = Terminal.MainScreenManager.GetScreen(screen.Id)
	}
	output := screen.Subscribe()

	go func() {
		for {
			screen.Publish()
			time.Sleep(100)
		}
	}()

	input := make(chan []byte, 1024)
	defer close(input)

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			input <- message
		}
	}()
	for {
		select {
		case o := <-output:
			err := conn.WriteMessage(1, o)
			if err != nil {
				break
			}
		case i := <-input:
			screen.InputByte(i)
		}

	}
}

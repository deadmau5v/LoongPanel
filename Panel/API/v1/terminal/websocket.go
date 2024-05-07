/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供 terminal 的流式传输
 */

package terminal

import (
	"LoongPanel/Panel/Service/Terminal"
	"fmt"
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
	w := c.Writer
	r := c.Request
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Ws链接异常关闭")
		}
	}(conn)

	id := getIntQuery(c, "id")
	screen := Terminal.MainScreenManager.GetScreen(id)
	if screen == nil {
		c.JSON(200, map[string]interface{}{
			"status": 500,
			"msg":    "无法获取screen",
		})
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

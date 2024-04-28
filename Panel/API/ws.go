package API

import (
	"fmt"
	"net/http"

	"LoongPanel/Panel/Terminal"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func screenWs(c *gin.Context) {
	w := c.Writer
	r := c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
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

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		screen.InputByte(message)
	}
}

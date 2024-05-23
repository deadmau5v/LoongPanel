/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供 terminal 的流式传输
 */

package terminal

import (
	"LoongPanel/Panel/Service/Log"
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
	Log.INFO("ScreenWs WebSocket 连接")
	w := c.Writer
	r := c.Request
	Log.DEBUG("升级为websocket连接")
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		Log.ERROR("无法打开websocket连接")
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer func(conn *websocket.Conn) {
		Log.DEBUG("websocket 连接关闭")
		err := conn.Close()
		if err != nil {
			Log.ERROR("Ws链接异常关闭")
		}
	}(conn)

	id := getIntQuery(c, "id")
	Log.DEBUG("获取到id为", id)
	screen := Terminal.MainScreenManager.GetScreen(uint32(id))
	Log.DEBUG(fmt.Sprintf("获取到的screen为%v", screen))
	if screen.WS != nil {
		screen.WS.Close()
		screen.WS = conn
	}
	Log.DEBUG("创建输入管道...")
	input := make(chan []byte, 1024)
	Log.DEBUG("创建输出管道...")
	output := screen.Subscribe()
	defer close(input)

	close_ := func() {
		Log.DEBUG("close_() websocket 关闭中...")
		id_, name_ := screen.Id, screen.Name
		Log.DEBUG(id_, name_)
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

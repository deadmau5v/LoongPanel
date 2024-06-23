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
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ScreenWs(c *gin.Context) {
	PanelLog.INFO("[网页终端]", "terminal WebSocket 连接")
	w := c.Writer
	r := c.Request
	host := c.Query("host")
	port := c.Query("port")
	user := c.Query("user")
	pwd := c.Query("pwd")
	PanelLog.DEBUG("[调试]", host, port, user, pwd)
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		PanelLog.ERROR("[网页终端]", "无法打开websocket连接")
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()
	PanelLog.DEBUG("[调试]", "链接中")
	err = Terminal.Shell(conn, host, port, user, pwd)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			PanelLog.ERROR("[网页终端]", "连接SSH失败")
		} else {
			PanelLog.ERROR("[网页终端]", err.Error())
		}

	}
}

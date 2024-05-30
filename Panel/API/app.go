/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：API入口 启动请在main.go中调用
 */

package API

import (
	"LoongPanel/Panel/Service/Auth"
	"LoongPanel/Panel/Service/PanelLog"
	"github.com/gin-gonic/gin"
	"net/http"
)

var App *gin.Engine

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	App = gin.New()
	App.Use(Cors())

	App.Use(gin.Logger())
	App.Use(PanelLog.GinLogToFile())
	App.Use(gin.Recovery())
	App.Use(Auth.UserAuth())

	initRoute(App)
}

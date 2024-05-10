/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：API入口 启动请在main.go中调用
 */

package API

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var App *gin.Engine
var WORKDIR string
var AppName = "LoongPanel"

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
	var err error
	WORKDIR, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	App = gin.Default()
	App.Use(Cors())

	gin.SetMode(gin.DebugMode)
	initRoute(App)
}

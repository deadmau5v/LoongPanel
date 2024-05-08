/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：API入口 启动请在main.go中调用
 */

package API

import (
	"github.com/gin-gonic/gin"
	"os"
)

var App *gin.Engine
var WORKDIR string
var AppName = "LoongPanel"

func init() {
	var err error
	WORKDIR, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	App = gin.Default()
	gin.SetMode(gin.DebugMode)
	initRoute(App)
}

/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：程序路由
 */

package API

import (
	"LoongPanel/Panel/API/v1/auth"
	"LoongPanel/Panel/API/v1/clean"
	"LoongPanel/Panel/API/v1/files"
	"LoongPanel/Panel/API/v1/home"
	"LoongPanel/Panel/API/v1/terminal"

	"github.com/gin-gonic/gin"
)

func initRoute(app *gin.Engine) {
	// 其他
	app.Static("/assets", WORKDIR+"/dist/assets")
	//  home 首页
	// -- home -> status 状态监控(实时)
	app.GET("/api/v1/status/system_status", home.SystemStatus)
	app.GET("/api/v1/status/system_info", home.SystemInfo)
	app.GET("/api/v1/status/disks", home.Disks)
	// -- home -> clean 清理垃圾
	app.GET("/api/v1/clean/pkg_auto_clean", clean.PkgAutoClean)
	// -- home -> power 电源操作
	app.GET("/api/v1/power/shutdown", home.Reboot)
	app.GET("/api/v1/power/reboot", home.Shutdown)

	//  files 文件
	app.GET("/api/v1/files/dir", files.FileDir)

	//  terminal 终端
	app.GET("/api/v1/screen/input", terminal.ScreenInput)
	app.GET("/api/v1/screen/create", terminal.ScreenCreate)
	app.GET("/api/v1/screen/close", terminal.ScreenClose)
	app.GET("/api/v1/screen/output", terminal.ScreenOutput)
	app.GET("/api/v1/screen/get_screens", terminal.GetScreens)
	// -- terminal -> WebSocket
	app.GET("/api/ws/screen", terminal.ScreenWs)

	// ping
	app.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "pong",
		})
	})

	// 登录
	app.POST("/api/v1/login", auth.Login)

	// 静态页面
	app.NoRoute(func(c *gin.Context) {
		c.File(WORKDIR + "/dist/index.html")
	})

}

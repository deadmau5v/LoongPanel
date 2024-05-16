/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：程序路由
 */

package API

import (
	AuthAPI "LoongPanel/Panel/API/v1/auth"
	"LoongPanel/Panel/API/v1/clean"
	"LoongPanel/Panel/API/v1/files"
	"LoongPanel/Panel/API/v1/home"
	"LoongPanel/Panel/API/v1/terminal"
	"LoongPanel/Panel/Service/System"
	"github.com/gin-gonic/gin"
)

func initRoute(app *gin.Engine) {
	// 其他
	app.Static("/assets", System.WORKDIR+"/dist/assets")

	// api v1
	v1 := app.Group("/api/v1")
	ws := app.Group("/api/ws")
	//  home 首页
	// -- home -> status 状态监控(实时)
	v1.GET("/status/system_status", home.SystemStatus)
	v1.GET("/status/system_info", home.SystemInfo)
	v1.GET("/status/disks", home.Disks)
	// -- home -> clean 清理垃圾
	v1.GET("/clean/pkg_auto_clean", clean.PkgAutoClean)
	// -- home -> power 电源操作
	v1.GET("/power/shutdown", home.Reboot)
	v1.GET("/power/reboot", home.Shutdown)

	//  files 文件
	v1.GET("/files/dir", files.FileDir)

	//  terminal 终端
	v1.GET("/screen/input", terminal.ScreenInput)
	v1.GET("/screen/create", terminal.ScreenCreate)
	v1.GET("/screen/close", terminal.ScreenClose)
	v1.GET("/screen/output", terminal.ScreenOutput)
	v1.GET("/screen/get_screens", terminal.GetScreens)
	// -- terminal -> WebSocket
	ws.GET("/api/ws/screen", terminal.ScreenWs)

	// ping
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "pong",
		})
	})

	// 登录
	app.POST("/api/v1/login", AuthAPI.Login)

	// 静态页面
	app.NoRoute(func(c *gin.Context) {
		c.File(System.WORKDIR + "/dist/index.html")
	})

}

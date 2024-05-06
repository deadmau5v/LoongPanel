package API

import (
	"LoongPanel/Panel/API/v1/files"
	"LoongPanel/Panel/API/v1/home"
	"LoongPanel/Panel/API/v1/terminal"
	"github.com/gin-gonic/gin"
)

func initRoute(app *gin.Engine) {
	// 静态页面
	app.NoRoute(func(c *gin.Context) {
		c.File(WORKDIR + "/Panel/Front/LoongPanel/dist/index.html")
	})

	// 其他
	app.Static("/assets", WORKDIR+"/Panel/Front/LoongPanel/dist/assets")

	// home 页面
	app.GET("/api/v1/status/system_status", home.SystemStatus)
	app.GET("/api/v1/status/system_info", home.SystemInfo)
	app.GET("/api/v1/status/disks", home.Disks)
	// files 页面
	app.GET("/api/v1/files/dir", files.FileDir)
	// terminal 页面
	app.GET("/api/v1/screen/input", terminal.ScreenInput)
	app.GET("/api/v1/screen/create", terminal.ScreenCreate)
	app.GET("/api/v1/screen/close", terminal.ScreenClose)
	app.GET("/api/v1/screen/output", terminal.ScreenOutput)
	app.GET("/api/v1/screen/get_screens", terminal.GetScreens)
	// WebSocket
	app.GET("/api/ws/screen", func(ctx *gin.Context) {
		terminal.ScreenWs(ctx)
	})
}

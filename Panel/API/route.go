package API

import "github.com/gin-gonic/gin"

func initRoute(app *gin.Engine) {
	// 静态页面
	app.NoRoute(func(c *gin.Context) {
		c.File(WORKDIR + "/Panel/Front/LoongPanel/dist/index.html")
	})

	// 其他
	app.Static("/assets", WORKDIR+"/Panel/Front/LoongPanel/dist/assets")

	// API
	app.GET("/api/v1/status/system_status", SystemStatus)
	app.GET("/api/v1/status/system_info", SystemInfo)
	app.GET("/api/v1/status/disks", Disks)
	app.GET("/api/v1/files/dir", FileDir)

	app.GET("/api/v1/screen/input", screenInput)
	app.GET("/api/v1/screen/create", screenCreate)
	app.GET("/api/v1/screen/close", screenClose)
	app.GET("/api/v1/screen/output", screenOutput)
	app.GET("/api/v1/screen/get_screens", getScreens)

	// WebSocket
	app.GET("/api/ws/screen", func(ctx *gin.Context) {
		screenWs(ctx)
	})
}

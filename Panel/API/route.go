package API

import "github.com/gin-gonic/gin"

func initRoute(app *gin.Engine) {
	// 静态页面
	app.GET("/", home)
	app.GET("/files", files)
	app.GET("/terminal", terminal)

	// 其他
	app.Static("/static", WORKDIR+"/Web/static")
	app.GET("/favicon.ico", func(context *gin.Context) {
		context.File(WORKDIR + "/Web/static/images/logo.png")
	})
	app.NoRoute(Error404)

	// API
	app.GET("/api/v1/status/cpu_percent", CPUPercent)
	app.GET("/api/v1/status/ram_percent", MemoryPercent)
	app.GET("/api/v1/status/system_info", SystemInfo)
	app.GET("/api/v1/status/disks", Disks)
	app.GET("/api/v1/files/dir", FileDir)
	app.GET("/api/v1/screen/input", screenInput)
	app.GET("/api/v1/screen/create", screenCreate)
	app.GET("/api/v1/screen/close", screenClose)
	app.GET("/api/v1/screen/output", screenOutput)
	app.GET("/api/v1/screen/get_screens", getScreens)
}

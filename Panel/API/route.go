package API

import "github.com/gin-gonic/gin"

func initRoute(app *gin.Engine) {
	// 静态页面
	app.GET("/", Home)
	app.GET("/files", Files)

	// 其他
	app.Static("/static", WORKDIR+"/Web/static")
	app.GET("/favicon.ico", func(context *gin.Context) {
		context.File(WORKDIR + "/Web/static/images/logo.png")
	})
	app.NoRoute(Error404)

	// API
	app.GET("/api/v1/cpu_percent", CPUPercent)
	app.GET("/api/v1/ram_percent", MemoryPercent)
	app.GET("/api/v1/system_info", SystemInfo)
	app.GET("/api/v1/disks", Disks)
	app.GET("/api/v1/dir", FileDir)
}

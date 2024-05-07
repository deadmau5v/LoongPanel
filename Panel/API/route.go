/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：程序路由
 */

package API

import (
	"LoongPanel/Panel/API/v1/clean"
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
	app.GET("/api/ws/screen", func(ctx *gin.Context) {
		terminal.ScreenWs(ctx)
	})

}

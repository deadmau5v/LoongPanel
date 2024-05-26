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
	"LoongPanel/Panel/Service/Auth"
	"LoongPanel/Panel/Service/Log"
	"LoongPanel/Panel/Service/System"
	"github.com/gin-gonic/gin"
)

func SetRoute(Method string, Path string, HandlerFunc gin.HandlerFunc, group *gin.RouterGroup, comment string, Public bool) {
	Log.DEBUG("添加路由", Method, Path, comment)
	if group != nil {
		_, err := Auth.Authenticator.AddPolicy("admin", group.BasePath()+Path, Method)
		if Public {
			_, err = Auth.Authenticator.AddPolicy("user", group.BasePath()+Path, Method)
		}
		if err != nil {
			Log.ERROR("添加权限策略失败", err)
			panic(err)
			return
		}
		group.Handle(Method, Path, HandlerFunc)
	} else {
		_, err := Auth.Authenticator.AddPolicy("admin", Path, Method)
		_, err = Auth.Authenticator.AddPolicy("user", Path, Method)
		if err != nil {
			Log.ERROR("添加权限策略失败", err)
			panic(err)
			return
		}
		App.Handle(Method, Path, HandlerFunc)
	}
}

func initRoute(app *gin.Engine) {
	// 其他
	app.Static("/assets", System.WORKDIR+"/dist/assets")
	// api v1 ws
	v1 := app.Group("/api/v1")
	ws := app.Group("/api/ws")
	//  home 首页
	// -- home -> status 状态监控(实时)
	SetRoute("GET", "/status/system_status", home.SystemStatus, v1, "系统状态", true)
	SetRoute("GET", "/status/system_info", home.SystemInfo, v1, "系统信息", true)
	SetRoute("GET", "/status/disks", home.Disks, v1, "磁盘信息", true)
	// -- home -> clean 清理垃圾
	SetRoute("GET", "/clean/pkg_auto_clean", clean.PkgAutoClean, v1, "清理过期包", false)
	// -- home -> power 电源操作
	SetRoute("POST", "/power/shutdown", home.Reboot, v1, "关机操作", false)
	SetRoute("POST", "/power/reboot", home.Shutdown, v1, "重启操作", false)

	//  files 文件
	SetRoute("GET", "/files/dir", files.FileDir, v1, "获取文件列表", true)

	//  terminal 终端
	SetRoute("GET", "/screen/create", terminal.ScreenCreate, v1, "终端创建", false)
	SetRoute("GET", "/screen/close", terminal.ScreenClose, v1, "终端关闭", false)
	SetRoute("GET", "/screen/get_screens", terminal.GetScreens, v1, "获取终端列表", false)
	// -- terminal -> WebSocket
	SetRoute("GET", "/screen", terminal.ScreenWs, ws, "使用网页终端", false)

	// ping
	SetRoute("GET", "/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "pong",
		})
	}, v1, "权限测试", true)

	// 登录
	SetRoute("POST", "/api/v1/auth/login", AuthAPI.Login, nil, "登录", true)
	SetRoute("POST", "/api/v1/auth/logout", AuthAPI.Logout, nil, "登出", true)

	// 静态页面
	app.NoRoute(func(c *gin.Context) {
		Log.DEBUG("无路由访问...")
		c.File(System.WORKDIR + "/dist/index.html")
	})

	// 信任代理
	err := app.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		Log.DEBUG("设置信任代理失败", err)
		return
	}

}

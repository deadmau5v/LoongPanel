/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：程序路由
 */

package API

import (
	"LoongPanel/Panel/API/v1/appstore"
	AuthAPI "LoongPanel/Panel/API/v1/auth"
	"LoongPanel/Panel/API/v1/clean"
	"LoongPanel/Panel/API/v1/docker"
	"LoongPanel/Panel/API/v1/files"
	"LoongPanel/Panel/API/v1/home"
	"LoongPanel/Panel/API/v1/log"
	"LoongPanel/Panel/API/v1/terminal"
	"LoongPanel/Panel/Service/Auth"
	"LoongPanel/Panel/Service/PanelLog"
	"LoongPanel/Panel/Service/System"
	"path"

	"github.com/gin-gonic/gin"
)

func SetRoute(Method string, Path string, HandlerFunc gin.HandlerFunc, group *gin.RouterGroup, comment string, Public bool) {
	PanelLog.DEBUG("[添加路由]", Method, Path, comment)
	if group != nil {
		_, err := Auth.Authenticator.AddPolicy("admin", group.BasePath()+Path, Method)
		if Public {
			_, err = Auth.Authenticator.AddPolicy("user", group.BasePath()+Path, Method)
		}
		if err != nil {
			PanelLog.ERROR("[权限管理]", "添加权限策略失败", err)
			panic(err)
			return
		}
		group.Handle(Method, Path, HandlerFunc)
	} else {
		_, err := Auth.Authenticator.AddPolicy("admin", Path, Method)
		_, err = Auth.Authenticator.AddPolicy("user", Path, Method)
		if err != nil {
			PanelLog.ERROR("[权限管理]", "添加权限策略失败", err)
			panic(err)
			return
		}
		App.Handle(Method, Path, HandlerFunc)
	}
}

func initRoute(app *gin.Engine) {
	// 其他
	app.Static("/assets", path.Join(System.WORKDIR, "dist", "assets"))
	app.Static("/script/icons", path.Join(System.WORKDIR, "script", "icons"))
	// 路由组 v1 ws
	v1 := app.Group("/api/v1")
	ws := app.Group("/api/ws")

	// 模块路由
	GroupAuth := v1.Group("/auth")
	GroupStatus := v1.Group("/status")
	GroupClean := v1.Group("/clean")
	GroupPower := v1.Group("/power")
	GroupFiles := v1.Group("/files")
	GroupLog := v1.Group("/log")
	GroupAppStore := v1.Group("/appstore")
	GroupDocker := v1.Group("/docker")

	SetRoute("GET", "/system_status", home.SystemStatus, GroupStatus, "系统状态", true)
	SetRoute("GET", "/system_info", home.SystemInfo, GroupStatus, "系统信息", true)
	SetRoute("GET", "/disks", home.Disks, GroupStatus, "磁盘信息", true)
	SetRoute("GET", "/pkg_auto_clean", clean.PkgAutoClean, GroupClean, "清理过期包", false)
	SetRoute("GET", "/temp_dir_remove", clean.TempDirRemove, GroupClean, "清理临时目录", false)
	SetRoute("POST", "/shutdown", home.Reboot, GroupPower, "关机操作", false)
	SetRoute("POST", "/reboot", home.Shutdown, GroupPower, "重启操作", false)
	SetRoute("GET", "/dir", files.FileDir, GroupFiles, "获取文件列表", true)
	SetRoute("GET", "/read", files.FileRead, GroupFiles, "读取文件", true)
	SetRoute("GET", "/screen", terminal.Terminal, ws, "使用网页终端", false)
	SetRoute("GET", "/ping", ping, v1, "权限测试", true)
	SetRoute("POST", "/login", AuthAPI.Login, GroupAuth, "登录", true)
	SetRoute("POST", "/logout", AuthAPI.Logout, GroupAuth, "登出", true)
	SetRoute("GET", "/users", AuthAPI.GetUsers, GroupAuth, "获取全部用户", false)
	SetRoute("DELETE", "/users", AuthAPI.DelUsers, GroupAuth, "获取全部用户", false)
	SetRoute("GET", "/user/:id", AuthAPI.GetUser, GroupAuth, "获取用户", false)
	SetRoute("GET", "/user", AuthAPI.GetUser, GroupAuth, "获取用户", false)
	SetRoute("POST", "/user", AuthAPI.CreateUser, GroupAuth, "创建用户", false)
	SetRoute("PUT", "/user", AuthAPI.UpdateUser, GroupAuth, "更新用户", false)
	SetRoute("DELETE", "/user", AuthAPI.DeleteUser, GroupAuth, "删除用户", false)
	SetRoute("GET", "/role", AuthAPI.GetRoles, GroupAuth, "获取角色", false)
	SetRoute("POST", "/role", AuthAPI.CreateRole, GroupAuth, "创建角色", false)
	SetRoute("DELETE", "/role", AuthAPI.DeleteRole, GroupAuth, "删除角色", false)
	SetRoute("GET", "/logs", log.GetLogs, GroupLog, "获取可用日志", false)
	SetRoute("GET", "/log", log.GetLog, GroupLog, "获取日志", false)
	SetRoute("GET", "/options", log.GetLogStruct, GroupLog, "获取日志结构", false)
	SetRoute("DELETE", "/log", log.ClearLog, GroupLog, "清理日志", false)
	SetRoute("GET", "/apps", appstore.AppList, GroupAppStore, "获取应用列表", true)
	SetRoute("POST", "/app", appstore.InstallApp, GroupAppStore, "安装应用", false)
	SetRoute("DELETE", "/app", appstore.UninstallApp, GroupAppStore, "卸载应用", false)
	SetRoute("POST", "/app/start", appstore.StartApp, GroupAppStore, "启动应用", false)
	SetRoute("POST", "/app/stop", appstore.StopApp, GroupAppStore, "停止应用", false)
	SetRoute("GET", "/containers", docker.GetContainerList, GroupDocker, "获取容器列表", false)
	SetRoute("GET", "/images", docker.GetImageList, GroupDocker, "获取镜像列表", false)

	// 增加参数记得检查Auth匹配 Service.Auth.SESSIONS.PathParse

	app.NoRoute(func(c *gin.Context) {
		c.File(System.WORKDIR + "/dist/index.html")
	})

	// 信任代理
	err := app.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		PanelLog.DEBUG("[权限管理] 设置信任代理失败", err)
		return
	}
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "pong",
	})
}

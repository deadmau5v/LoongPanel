/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：程序路由
 */

package API

import (
	"LoongPanel/Panel/API/v1/appstore"
	AuthAPI "LoongPanel/Panel/API/v1/auth"
	"LoongPanel/Panel/API/v1/clamav"
	"LoongPanel/Panel/API/v1/clean"
	"LoongPanel/Panel/API/v1/docker"
	"LoongPanel/Panel/API/v1/files"
	"LoongPanel/Panel/API/v1/home"
	"LoongPanel/Panel/API/v1/log"
	"LoongPanel/Panel/API/v1/status"
	"LoongPanel/Panel/API/v1/terminal"
	"LoongPanel/Panel/Service/Auth"
	"LoongPanel/Panel/Service/PanelLog"
	"LoongPanel/Panel/Service/System"
	"path"

	"github.com/gin-gonic/gin"
)

func SetRoute(Method, Path string, HandlerFunc gin.HandlerFunc, group *gin.RouterGroup, comment string, Public bool) {
	PanelLog.DEBUG("[添加路由]", Method, Path, comment)

	addPolicy := func(role, path, method string) {
		if _, err := Auth.Authenticator.AddPolicy(role, path, method); err != nil {
			PanelLog.ERROR("[权限管理]", "添加权限策略失败", err)
			panic(err)
		}
	}

	if group != nil {
		basePath := group.BasePath() + Path
		addPolicy("admin", basePath, Method)
		if Public {
			addPolicy("user", basePath, Method)
		}
		group.Handle(Method, Path, HandlerFunc)
	} else {
		addPolicy("admin", Path, Method)
		addPolicy("user", Path, Method)
		App.Handle(Method, Path, HandlerFunc)
	}
}

func initRoute(app *gin.Engine) {
	// 其他
	app.Static("/assets", path.Join(System.WORKDIR, "dist", "assets"))
	app.Static("/script/icons", path.Join(System.WORKDIR, "script", "icons"))

	v1 := app.Group("/api/v1")

	// 公共路由
	SetRoute("GET", "/ping", ping, v1, "权限测试", true)

	// websocket
	ws := app.Group("/api/ws")
	SetRoute("GET", "/screen", terminal.Terminal, ws, "使用网页终端", false)
	SetRoute("GET", "/status", status.GetStatus, ws, "获取状态", false)

	// 权限管理
	GroupAuth := v1.Group("/auth")
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
	SetRoute("POST", "/policy", AuthAPI.AddPolicy, GroupAuth, "获取策略", false)
	SetRoute("DELETE", "/policy", AuthAPI.DeletePolicy, GroupAuth, "删除策略", false)
	SetRoute("POST", "/password", AuthAPI.ChangePassword, GroupAuth, "修改密码", false)

	// 系统状态
	GroupStatus := v1.Group("/status")
	SetRoute("GET", "/system_status", home.SystemStatus, GroupStatus, "系统状态", true)
	SetRoute("GET", "/system_info", home.SystemInfo, GroupStatus, "系统信息", true)
	SetRoute("GET", "/disks", home.Disks, GroupStatus, "磁盘信息", true)
	SetRoute("POST", "/time_step", status.SetStatusStepTime, GroupStatus, "设置状态保存间隔 0为关闭", false)
	SetRoute("POST", "/save_time", status.SetSaveTime, GroupStatus, "设置状态保存时间", false)

	// 清理
	GroupClean := v1.Group("/clean")
	SetRoute("GET", "/pkg_auto_clean", clean.PkgAutoClean, GroupClean, "清理过期包", false)
	SetRoute("GET", "/temp_dir_remove", clean.TempDirRemove, GroupClean, "清理临时目录", false)

	// 系统操作
	GroupPower := v1.Group("/power")
	SetRoute("POST", "/shutdown", home.Reboot, GroupPower, "关机操作", false)
	SetRoute("POST", "/reboot", home.Shutdown, GroupPower, "重启操作", false)

	// 文件操作
	GroupFiles := v1.Group("/files")
	SetRoute("GET", "/dir", files.FileDir, GroupFiles, "获取文件列表", true)
	SetRoute("GET", "/read", files.FileRead, GroupFiles, "读取文件", true)

	// 日志
	GroupLog := v1.Group("/log")
	SetRoute("GET", "/logs", log.GetLogs, GroupLog, "获取可用日志", false)
	SetRoute("GET", "/log", log.GetLog, GroupLog, "获取日志", false)
	SetRoute("GET", "/options", log.GetLogStruct, GroupLog, "获取日志结构", false)
	SetRoute("DELETE", "/log", log.ClearLog, GroupLog, "清理日志", false)

	// 应用商店
	GroupAppStore := v1.Group("/appstore")
	SetRoute("GET", "/apps", appstore.AppList, GroupAppStore, "获取应用列表", true)
	SetRoute("POST", "/app", appstore.InstallApp, GroupAppStore, "安装应用", false)
	SetRoute("DELETE", "/app", appstore.UninstallApp, GroupAppStore, "卸载应用", false)
	SetRoute("POST", "/app/start", appstore.StartApp, GroupAppStore, "启动应用", false)
	SetRoute("POST", "/app/stop", appstore.StopApp, GroupAppStore, "停止应用", false)

	// Docker
	GroupDocker := v1.Group("/docker")
	SetRoute("GET", "/containers", docker.GetContainerList, GroupDocker, "获取容器列表", false)
	SetRoute("GET", "/images", docker.GetImageList, GroupDocker, "获取镜像列表", false)

	// 病毒扫描
	GroupClamav := v1.Group("/clamav")
	SetRoute("GET", "/scan", clamav.ScanFile, GroupClamav, "扫描文件", false)
	SetRoute("GET", "/clamav/scan", clamav.ScanFile, ws, "快速扫描", false)
	SetRoute("GET", "/scan_dir", clamav.ScanDir, GroupClamav, "扫描目录", false)
	SetRoute("GET", "/clamav/scan_dir", clamav.ScanFile, ws, "扫描目录", false)

	// 前端静态文件
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
		"status": 0,
		"msg":    "pong",
	})
}

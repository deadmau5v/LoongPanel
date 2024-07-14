/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：程序路由
 */

package API

import (
	v1 "LoongPanel/Panel/API/v1"
	inspection "LoongPanel/Panel/API/v1/Inspection"
	"LoongPanel/Panel/API/v1/appstore"
	AuthAPI "LoongPanel/Panel/API/v1/auth"
	"LoongPanel/Panel/API/v1/clamav"
	"LoongPanel/Panel/API/v1/clean"
	"LoongPanel/Panel/API/v1/docker"
	"LoongPanel/Panel/API/v1/files"
	"LoongPanel/Panel/API/v1/home"
	"LoongPanel/Panel/API/v1/log"
	"LoongPanel/Panel/API/v1/notice"
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
		v1.AddRouteComment(Method, basePath, comment)
	} else {
		addPolicy("admin", Path, Method)
		addPolicy("user", Path, Method)
		App.Handle(Method, Path, HandlerFunc)
		v1.AddRouteComment(Method, Path, comment)
	}
}

func initRoute(app *gin.Engine) {
	// 其他
	app.Static("/assets", path.Join(System.WORKDIR, "dist", "assets"))
	app.Static("/script/icons", path.Join(System.WORKDIR, "script", "icons"))

	v1 := app.Group("/api/v1")
	ws := app.Group("/api/ws")

	// 公共路由
	SetRoute("GET", "/ping", ping, v1, "基础(检测Cookie是否生效)", true)

	// 权限管理
	GroupAuth := v1.Group("/auth")
	SetRoute("POST", "/login", AuthAPI.Login, GroupAuth, "基础(登录)", true)
	SetRoute("POST", "/logout", AuthAPI.Logout, GroupAuth, "基础(登出)", true)
	SetRoute("GET", "/users", AuthAPI.GetUsers, GroupAuth, "权限管理(获取全部用户)", false)
	SetRoute("DELETE", "/users", AuthAPI.DelUsers, GroupAuth, "权限管理(批量删除用户)", false)
	SetRoute("GET", "/user/:id", AuthAPI.GetUser, GroupAuth, "权限管理(获取用户)", false)
	SetRoute("GET", "/user", AuthAPI.GetUser, GroupAuth, "权限管理(获取用户)", false)
	SetRoute("POST", "/user", AuthAPI.CreateUser, GroupAuth, "权限管理(创建用户)", false)
	SetRoute("PUT", "/user", AuthAPI.UpdateUser, GroupAuth, "权限管理(更新用户)", false)
	SetRoute("DELETE", "/user", AuthAPI.DeleteUser, GroupAuth, "权限管理(删除用户)", false)
	SetRoute("GET", "/role", AuthAPI.GetRoles, GroupAuth, "权限管理(获取角色)", false)
	SetRoute("POST", "/role", AuthAPI.CreateRole, GroupAuth, "权限管理(创建角色)", false)
	SetRoute("DELETE", "/role", AuthAPI.DeleteRole, GroupAuth, "权限管理(删除角色)", false)
	SetRoute("POST", "/policy", AuthAPI.AddPolicy, GroupAuth, "权限管理(获取策略)", false)
	SetRoute("DELETE", "/policy", AuthAPI.DeletePolicy, GroupAuth, "权限管理(删除策略)", false)
	SetRoute("POST", "/password", AuthAPI.ChangePassword, GroupAuth, "权限管理(修改密码)", false)

	// 系统状态
	GroupStatus := v1.Group("/status")
	SetRoute("GET", "/system_status", home.SystemStatus, GroupStatus, "系统状态(系统实时状态)", true)
	SetRoute("GET", "/system_info", home.SystemInfo, GroupStatus, "系统状态(系统信息)", true)
	SetRoute("GET", "/disks", home.Disks, GroupStatus, "系统状态(磁盘信息)", true)
	SetRoute("POST", "/config", status.SetStatusConfig, GroupStatus, "系统状态(设置状态保存时间和间隔)", false)

	SetRoute("GET", "/status", status.GetStatus, ws, "系统状态(流式传输获取状态)", false)

	// 清理
	GroupClean := v1.Group("/clean")
	SetRoute("GET", "/pkg_auto_clean", clean.PkgAutoClean, GroupClean, "包管理工具(清理过期包)", false)
	SetRoute("GET", "/temp_dir_remove", clean.TempDirRemove, GroupClean, "包管理工具(清理临时目录)", false)

	// 系统操作
	GroupPower := v1.Group("/power")
	SetRoute("POST", "/shutdown", home.Reboot, GroupPower, "系统操作(关机操作)", false)
	SetRoute("POST", "/reboot", home.Shutdown, GroupPower, "系统操作(重启操作)", false)

	// 文件操作
	GroupFiles := v1.Group("/files")
	SetRoute("GET", "/dir", files.FileDir, GroupFiles, "文件管理(获取文件列表)", true)
	SetRoute("GET", "/read", files.FileRead, GroupFiles, "文件管理(读取文件)", false)
	SetRoute("POST", "/upload", files.Upload, GroupFiles, "文件管理(上传文件)", false)
	SetRoute("POST", "/download", files.Download, GroupFiles, "文件管理(下载文件)", false)
	SetRoute("POST", "/delete", files.Delete, GroupFiles, "文件管理(删除文件)", false)
	SetRoute("POST", "/move", files.Move, GroupFiles, "文件管理(移动文件)", false)
	SetRoute("POST", "/rename", files.Rename, GroupFiles, "文件管理(重命名文件)", false)
	SetRoute("POST", "/decompress", files.Decompress, GroupFiles, "文件管理(解压文件)", false)
	SetRoute("POST", "/compress", files.Compress, GroupFiles, "文件管理(压缩文件)", false)

	// 日志
	GroupLog := v1.Group("/log")
	SetRoute("GET", "/logs", log.GetLogs, GroupLog, "日志(获取可用日志)", false)
	SetRoute("GET", "/log", log.GetLog, GroupLog, "日志(获取日志)", false)
	SetRoute("GET", "/options", log.GetLogStruct, GroupLog, "日志(获取日志结构)", false)
	SetRoute("DELETE", "/log", log.ClearLog, GroupLog, "日志(清理日志)", false)

	// 应用商店
	GroupAppStore := v1.Group("/appstore")
	SetRoute("GET", "/apps", appstore.AppList, GroupAppStore, "应用商店(获取应用列表)", true)
	SetRoute("POST", "/app", appstore.InstallApp, GroupAppStore, "应用商店(安装应用)", false)
	SetRoute("DELETE", "/app", appstore.UninstallApp, GroupAppStore, "应用商店(卸载应用)", false)
	SetRoute("POST", "/app/start", appstore.StartApp, GroupAppStore, "应用商店(启动应用)", false)
	SetRoute("POST", "/app/stop", appstore.StopApp, GroupAppStore, "应用商店(停止应用)", false)

	// Docker
	GroupDocker := v1.Group("/docker")
	SetRoute("GET", "/containers", docker.GetContainerList, GroupDocker, "Docker(获取容器列表)", false)
	SetRoute("GET", "/images", docker.GetImageList, GroupDocker, "Docker(获取镜像列表)", false)

	// 病毒扫描
	GroupClamav := v1.Group("/clamav")
	SetRoute("GET", "/scan", clamav.ScanFile, GroupClamav, "病毒扫描(扫描文件)", false)
	SetRoute("GET", "/scan_dir", clamav.ScanDir, GroupClamav, "病毒扫描(扫描目录)", false)
	SetRoute("GET", "/set_scan_time", clamav.SetScanTime, GroupClamav, "病毒扫描(设置定时扫描)", false)

	SetRoute("GET", "/clamav/scan", clamav.ScanFile, ws, "病毒扫描(快速扫描)", false)
	SetRoute("GET", "/clamav/scan_dir", clamav.ScanFile, ws, "病毒扫描(快速扫描目录)", false)

	// 网页终端
	SetRoute("GET", "/screen", terminal.Terminal, ws, "网页终端(使用网页终端)", false)

	// 巡检
	SetRoute("GET", "/check", inspection.Check, ws, "巡检(一键巡检)", false)

	// 预警通知
	GroupNotice := v1.Group("/notice")
	SetRoute("GET", "/notices", notice.GetAllSettings, GroupNotice, "预警通知(获取所有通知设置)", true)
	SetRoute("POST", "/notice", notice.AddNotice, GroupNotice, "预警通知(添加通知设置)", false)
	SetRoute("DELETE", "/notice", notice.DeleteNotice, GroupNotice, "预警通知(删除通知设置)", false)
	SetRoute("PUT", "/notice", notice.UpdateNotice, GroupNotice, "预警通知(更新通知设置)", false)

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

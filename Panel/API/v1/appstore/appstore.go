/*
 * 创建人： deadmau5v
 * 创建时间： 2024-6-17
 * 文件作用：应用商店 API
 */

package appstore

import (
	"LoongPanel/Panel/Service/AppStore"
	"LoongPanel/Panel/Service/PanelLog"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AppList(c *gin.Context) {
	type Res struct {
		Name      string   `json:"name"`
		Version   string   `json:"version"`
		Tags      []string `json:"tags"`
		Icon      string   `json:"icon"`
		IsInstall bool     `json:"is_install"`
		IsRunning bool     `json:"is_running"`
	}

	res := make([]Res, 0)

	for _, app := range AppStore.Apps {
		version, err := app.Version()
		if err != nil {
			PanelLog.ERROR("[应用商店]", "获取应用版本失败", err)
			version = "未知"
		}
		res = append(res, Res{
			Name:      app.Name,
			Version:   version,
			Tags:      app.Tags,
			Icon:      app.Icon,
			IsInstall: app.IsInstall(),
			IsRunning: app.IsRunning(),
		})
	}
	PanelLog.INFO("[应用商店]", "获取应用列表")
	c.JSON(200, gin.H{
		"status": 0,
		"data":   res,
	})
}

func StartApp(c *gin.Context) {
	name := c.Query("name")
	if strings.Trim(name, " ") == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "参数错误",
		})
	}

	app := AppStore.FindApp(name)
	ok, err := app.Start()
	// 启动错误
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    err.Error(),
		})
		return
	}
	if !ok {
		// 未启动成功
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "启动失败",
		})
		return
	}
	// 启动成功
	PanelLog.INFO("[应用商店]", "启动应用", name)
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "启动成功",
	})

}

func StopApp(c *gin.Context) {
	name := c.Query("name")
	if strings.Trim(name, " ") == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "参数错误",
		})
	}

	app := AppStore.FindApp(name)
	ok, err := app.Stop()
	// 停止错误
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    err.Error(),
		})
		return
	}
	if !ok {
		// 未停止成功
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "停止失败",
		})
		return
	}
	// 停止成功
	PanelLog.INFO("[应用商店]", "停止应用", name)
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "停止成功",
	})
}

func InstallApp(c *gin.Context) {
	// 长时间阻塞
	name := c.Query("name")
	if strings.Trim(name, " ") == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "参数错误",
		})
	}

	app := AppStore.FindApp(name)
	ok, err := app.Install()
	// 安装错误
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    err.Error(),
		})
		return
	}
	if !ok {
		// 未安装成功
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "安装失败",
		})
		return
	}
	// 安装成功
	PanelLog.INFO("[应用商店]", "安装应用", name)
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "安装成功",
	})
}

func UninstallApp(c *gin.Context) {
	// 长时间阻塞
	name := c.Query("name")
	if strings.Trim(name, " ") == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "参数错误",
		})
	}

	app := AppStore.FindApp(name)
	ok, err := app.Uninstall()
	// 卸载错误
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    err.Error(),
		})
		return
	}
	if !ok {
		// 未卸载成功
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "卸载失败",
		})
		return
	}
	// 卸载成功
	PanelLog.INFO("[应用商店]", "卸载应用", name)
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "卸载成功",
	})
}

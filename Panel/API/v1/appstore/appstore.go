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
			PanelLog.ERROR("获取应用版本失败", err)
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

	c.JSON(200, res)
}

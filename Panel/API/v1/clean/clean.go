/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供清理垃圾的API 主要实现在 Service/Clean 包中
 */

package clean

import (
	"LoongPanel/Panel/Service/Clean"
	"LoongPanel/Panel/Service/System"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func PkgAutoClean(ctx *gin.Context) {
	if System.Data.OSName == "windows" {
		slog.Info("Windows OS 无法使用 PkgAutoClean 函数，跳过")
		ctx.JSON(200, gin.H{
			"msg":    "Windows OS does not support this function, pass",
			"status": 0,
		})
		return
	}
	clean, err := Clean.AptAutoClean()
	if err != nil {
		slog.Error("AptAutoClean", err)
		return
	}
	slog.Debug("AptAutoClean", string(clean))
	remove, err := Clean.AptAutoRemove()
	if err != nil {
		slog.Error("AptAutoRemove", err)
		return
	}
	slog.Debug("AptAutoRemove", string(remove))
	autoClean, err := Clean.YumAutoClean()
	if err != nil {
		slog.Error("YumAutoClean", err)
		return
	}
	slog.Debug("YumAutoClean", string(autoClean))
	return
}

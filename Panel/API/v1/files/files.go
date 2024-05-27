/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供文件管理相关的API 主要实现在 Service/Files 中
 */

package files

import (
	FileService "LoongPanel/Panel/Service/Files"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FileDir(ctx *gin.Context) {
	path := ctx.Query("path")
	if path == "" {
		path = "/"
	}

	files, err := FileService.Dir(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":    err.Error(),
			"status": -1,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"files":  files,
			"status": 0,
			"msg":    "ok",
		})
	}
}

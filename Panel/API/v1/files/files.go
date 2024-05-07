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
	data := map[string]interface{}{}
	var err error
	data["files"], err = FileService.Dir(path)
	if err != nil {
		data["status"] = -1
		data["msg"] = err.Error()
	} else {
		data["status"] = 0
		data["msg"] = ""
	}

	ctx.JSON(http.StatusOK, data)
}

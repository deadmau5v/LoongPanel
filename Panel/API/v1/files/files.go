package files

import (
	FileService "LoongPanel/Panel/Files"
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

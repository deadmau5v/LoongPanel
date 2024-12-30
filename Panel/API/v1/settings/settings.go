package settings

import (
	config "LoongPanel/Panel/Service/Config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMailConfig(ctx *gin.Context) {
	mailConfig := config.GetMailConfig()
	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"data":   mailConfig,
	})
}

func SetMailConfig(ctx *gin.Context) {
	mailConfig := config.MailConfig{}
	if err := ctx.BindJSON(&mailConfig); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    err.Error(),
		})
		return
	}
	config.SetMailConfig(mailConfig)
	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "设置成功",
	})
}

func GetAuthConfig(ctx *gin.Context) {
	authConfig := config.GetAuthConfig()
	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"data":   authConfig,
	})
}

func SetAuthConfig(ctx *gin.Context) {
	authConfig := config.AuthConfig{}
	if err := ctx.BindJSON(&authConfig); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    err.Error(),
		})
		return
	}
	config.SetAuthConfig(authConfig)
	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "设置成功",
	})
}

func GetPanelLogConfig(ctx *gin.Context) {
	panelLogConfig := config.GetPanelLogConfig()
	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"data":   panelLogConfig,
	})
}

func SetPanelLogConfig(ctx *gin.Context) {
	panelLogConfig := config.PanelLogConfig{}
	if err := ctx.BindJSON(&panelLogConfig); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    err.Error(),
		})
		return
	}
	config.SetPanelLogConfig(panelLogConfig)
	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "设置成功",
	})
}

func GetClamavConfig(ctx *gin.Context) {
	clamavConfig := config.GetClamavConfig()
	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"data":   clamavConfig,
	})
}

func SetClamavConfig(ctx *gin.Context) {
	clamavConfig := config.ClamavConfig{}
	if err := ctx.BindJSON(&clamavConfig); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    err.Error(),
		})
		return
	}
	config.SetClamavConfig(clamavConfig)
	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "设置成功",
	})
}

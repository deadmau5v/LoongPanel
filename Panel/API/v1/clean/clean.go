/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供清理垃圾的API 主要实现在 Service/Clean 包中
 */

package clean

import (
	"LoongPanel/Panel/Service/Clean"
	"LoongPanel/Panel/Service/Log"
	"LoongPanel/Panel/Service/System"
	"github.com/gin-gonic/gin"
)

func PkgAutoClean(ctx *gin.Context) {
	if System.Data.OSName == "windows" {
		Log.INFO("Windows OS 无法使用 PkgAutoClean 函数，跳过")
		ctx.JSON(200, gin.H{
			"msg":    "Windows OS does not support this function, pass",
			"status": 0,
		})
		return
	}
	// 检查包管理器
	switch System.Data.PkgManager {
	case "apt":
		msg, err := Clean.AptAutoClean()
		if err != nil {
			Log.ERROR("AptAutoClean() Error: ", err.Error())
			ctx.JSON(200, gin.H{
				"msg":    "AptAutoClean() Error: " + err.Error(),
				"status": 1,
			})
			return
		}
		msg2, err := Clean.AptAutoRemove()
		if err != nil {
			Log.ERROR("AptAutoRemove() Error: ", err.Error())
			ctx.JSON(200, gin.H{
				"msg":    "AptAutoRemove() Error: " + err.Error(),
				"status": 1,
			})
			return
		}
		ctx.JSON(200, gin.H{
			"msg":    msg + msg2,
			"status": 0,
		})

	case "yum":
		msg, err := Clean.YumAutoClean()
		if err != nil {
			Log.ERROR("YumAutoClean() Error: ", err.Error())
			ctx.JSON(200, gin.H{
				"msg":    "YumAutoClean() Error: " + err.Error(),
				"status": 1,
			})
			return
		}
		msg2, err := Clean.YumAutoRemove()
		if err != nil {
			Log.ERROR("YumAutoRemove() Error: ", err.Error())
			ctx.JSON(200, gin.H{
				"msg":    "YumAutoRemove() Error: " + err.Error(),
				"status": 1,
			})
			return
		}
		ctx.JSON(200, gin.H{
			"msg":    msg + msg2,
			"status": 0,
		})
	}
}

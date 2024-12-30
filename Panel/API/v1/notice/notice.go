package notice

import (
	notice "LoongPanel/Panel/Service/Notice"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllSettings(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"data": notice.GetAllSettings(), "status": 0})
}

func AddNotice(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Query("userID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": "参数错误"})
		return
	}
	notice.AddNotice(uint(userID))
	ctx.JSON(http.StatusOK, gin.H{"data": "添加成功", "status": 0})
}

func UpdateNotice(ctx *gin.Context) {
	var notice_ notice.UserNotificationSetting
	if err := ctx.ShouldBindJSON(&notice_); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": "参数错误"})
		return
	}

	notice.UpdateNotice(notice_)
	ctx.JSON(http.StatusOK, gin.H{"data": "更新成功", "status": 0})
}

func DeleteNotice(ctx *gin.Context) {
	var notice_ notice.UserNotificationSetting
	if err := ctx.ShouldBindJSON(&notice_); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": "参数错误"})
		return
	}

	notice.DeleteNotice(notice_)
	ctx.JSON(http.StatusOK, gin.H{"data": "删除成功", "status": 0})
}

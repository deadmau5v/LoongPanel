/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-16
 * 文件作用：储存登陆会话ID
 */

package Auth

import (
	"LoongPanel/Panel/Service/Database"
	"LoongPanel/Panel/Service/PanelLog"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"regexp"
	"time"
)

func NewSESSION(user Database.User) string {
	uuid_ := uuid.New().String()
	session := SESSION{
		KEY:      uuid_,
		User:     user,
		TimeUnix: time.Now().Unix(),
	}
	err := Database.DB.Create(&session).Error
	if err != nil {
		PanelLog.ERROR(err)
		return ""
	}
	return uuid_
}

type SESSION struct {
	KEY      string        `json:"key"`
	UserID   uint          `json:"user_id"`
	User     Database.User `json:"user" gorm:"foreignKey:UserID"`
	TimeUnix int64         `json:"create_time"`
}

func UserAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 放行静态资源
		skipPaths := []string{
			"/assets/*",
		}
		for _, path := range skipPaths {
			// 使用正则匹配
			if match, _ := regexp.MatchString(path, c.Request.URL.Path); match {
				c.Next()
				return
			}
		}
		staticPaths := []string{
			"/",
			"/index.html",
			"/favicon.ico",
			"/api/v1/auth/login",
			"/index",
			"/login",
			"/user",
			"/terminal",
			"/files",
		}
		for _, path := range staticPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		Authorization := c.GetHeader("Authorization")
		PanelLog.DEBUG("Authorization", Authorization)
		var SESSIONS []SESSION
		Database.DB.Find(&SESSIONS)
		var userSession SESSION
		var flag = false
		for _, session := range SESSIONS {
			if session.KEY == Authorization {
				userSession = session
				Database.DB.Preload("User").Find(&userSession)

				flag = true
				break
			}
		}
		if !flag {
			PanelLog.DEBUG("未授权1")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未授权",
			})
			c.Abort()
			return
		}

		ok, err := Authenticator.Enforce(userSession.User.Role, c.Request.URL.Path, c.Request.Method)
		PanelLog.DEBUG("权限验证", userSession.User, c.Request.URL.Path, c.Request.Method, ok, err)
		if ok && err == nil {
			c.Next()
			return
		} else if err != nil || !ok {
			PanelLog.DEBUG("未授权2")

			if err != nil {
				PanelLog.ERROR(err)
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未授权",
			})
			c.Abort()
			return

		}

	}
}

func init() {
	err := Database.DB.AutoMigrate(&SESSION{})
	if err != nil {
		PanelLog.DEBUG("[数据库模块] SESSIONS表创建失败")
		return
	}
}

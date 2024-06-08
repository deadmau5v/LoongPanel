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
			"/auth/global",
			"/auth/role",
			"/auth/user",
		}
		for _, path := range staticPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		Authorization := c.GetHeader("Authorization")
		if Authorization != "" {
			PanelLog.DEBUG("[权限管理] Authorization", Authorization)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未授权",
			})
			return
		}
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
			PanelLog.DEBUG("[权限管理] 未授权 code: 1")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未授权",
			})
			c.Abort()
			return
		}
		path := PathParse(c.Request.URL.Path)
		ok, err := Authenticator.Enforce(userSession.User.Role, path, c.Request.Method)
		status := "未通过"
		if ok {
			status = "通过"
		}
		PanelLog.DEBUG("[权限管理] 权限验证", userSession.User.Name, c.Request.URL.Path, c.Request.Method, status)
		if ok && err == nil {
			c.Next()
			return
		} else if err != nil || !ok {
			PanelLog.DEBUG("[权限管理] 未授权 code: 2")

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

// PathParse 路径解析
func PathParse(path string) string {
	for k, v := range map[string]string{
		`/api/v1/auth/user/(\d+)$`: "/api/v1/auth/user/:id",
	} {
		if match, _ := regexp.MatchString(k, path); match {
			return v
		}
	}
	return path
}

func GetSessionByKey(key string) (*SESSION, error) {
	var Session SESSION
	Database.DB.Model(&SESSION{}).Where("`key` = ?", key).Find(&Session)
	// 获取关联的用户
	Database.DB.Model(&Session).Preload("User").Find(&Session)
	return &Session, nil
}

func init() {
	err := Database.DB.AutoMigrate(&SESSION{})
	if err != nil {
		PanelLog.DEBUG("[数据库模块] SESSIONS表创建失败")
		return
	}
}

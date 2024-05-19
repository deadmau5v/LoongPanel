/*
 * 创建人： deadmau5v
 * 创建时间： 2024-0-0
 * 文件作用：
 */

/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-16
 * 文件作用：
 */

package Auth

import (
	"LoongPanel/Panel/Service/Log"
	"LoongPanel/Panel/Service/User"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"regexp"
)

var SESSIONS = map[string]string{}

func RandomSESSION(username string) string {
	uuid_ := uuid.New().String()
	SESSIONS[uuid_] = username
	return uuid_
}

func AuthUser() gin.HandlerFunc {
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
		}
		for _, path := range staticPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		session, _ := c.Cookie("SESSION")
		Authorization := c.GetHeader("Authorization")
		if Authorization != "" {
			session = Authorization
		}

		if SESSIONS[session] != "" {
			username := SESSIONS[session]
			users := User.Find()
			var user User.User
			for _, u := range users {
				if u.Name == username {
					user = u
					break
				}
			}

			ok, err := Authenticator.Enforce(user.Role, c.Request.URL.Path, c.Request.Method)
			if ok && err == nil {
				c.Next()
				return
			} else {
				if err != nil {
					Log.ERROR(err)
				}
				session = ""
			}
		} else {
			session = ""
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		c.Abort()
	}
}

package Auth

import (
	"LoongPanel/Panel/Service/Database"
	"LoongPanel/Panel/Service/PanelLog"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"sync"
	"time"
)

const (
	sessionIDLength = 32
	CookieName      = "session_token"
	CookieMaxAge    = 86400 * 30 // 30 days
)

var (
	sessions     = make(map[string]*Session)
	sessionMutex sync.RWMutex
)

// Session 表示用户会话
type Session struct {
	Username string
	Expiry   time.Time
}

// IsExpired 检查会话是否已过期
func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

// CreateSession 创建新的会话并返回会话ID
func CreateSession(username string) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		PanelLog.ERROR("[权限管理]", "创建会话失败: "+err.Error())
		return "", err
	}

	expiryTime := time.Now().Add(30 * 24 * time.Hour)

	sessionMutex.Lock()
	sessions[sessionID] = &Session{
		Username: username,
		Expiry:   expiryTime,
	}
	sessionMutex.Unlock()

	return sessionID, nil
}

// GetSession 根据会话ID获取会话信息
func GetSession(sessionID string) (*Session, bool) {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()
	session, exists := sessions[sessionID]
	return session, exists
}

// DeleteSession 删除指定的会话
func DeleteSession(sessionID string) {
	sessionMutex.Lock()
	delete(sessions, sessionID)
	sessionMutex.Unlock()
}

// SetSessionCookie 设置会话cookie
func SetSessionCookie(c *gin.Context, sessionID string) {
	c.SetCookie(
		CookieName,
		sessionID,
		CookieMaxAge,
		"/",
		"",
		false,
		true,
	)
	c.Next()
}

// ClearSessionCookie 清除会话cookie
func ClearSessionCookie(c *gin.Context) {
	c.SetCookie(
		CookieName,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
	c.Next()
}

// generateSessionID 生成随机的会话ID
func generateSessionID() (string, error) {
	b := make([]byte, sessionIDLength)
	_, err := rand.Read(b)
	if err != nil {
		PanelLog.ERROR("[权限管理]", "生成会话ID失败: "+err.Error())
		return "", errors.New("failed to generate session ID")
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// UserAuth 用户认证 Gin 中间件
func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 放行静态资源
		skipPaths := []string{
			"/assets/*",
			"/script/icons/*",
			"/api/ws/*",
		}
		for _, path := range skipPaths {
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
			"/log",
			"/appstore",
		}
		for _, path := range staticPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		// 从请求中获取session token
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "未授权访问", "status": 401})
			c.Abort()
			return
		}

		// 验证session
		session, exists := GetSession(sessionToken)
		if !exists || session.IsExpired() {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "会话已过期或无效", "status": 401})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("username", session.Username)

		// 检查用户权限
		if !CheckUserPolicy(c) {
			c.JSON(http.StatusForbidden, gin.H{"msg": "没有权限访问", "status": 403})
			c.Abort()
			return
		}
		// 继续处理请求
		c.Next()
	}
}

func CheckUserPolicy(ctx *gin.Context) bool {
	username, exists := ctx.Get("username")
	if !exists {
		return false
	}

	user := Database.User{}
	result := Database.DB.Model(&Database.User{}).Where("name = ?", username).First(&user)
	if result.Error != nil {
		return false
	}

	enforce, err := Authenticator.Enforce(user.Role, ctx.Request.URL.Path, ctx.Request.Method)
	//PanelLog.DEBUG("[调试]", enforce, err, user.Role, ctx.Request.URL.Path, ctx.Request.Method)
	if err != nil {
		PanelLog.ERROR("[权限管理] 验证权限失败", err.Error())
		return false
	}

	return enforce
}

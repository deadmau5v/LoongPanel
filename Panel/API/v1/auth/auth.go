/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-18
 * 文件作用：使用Casbin进行权限验证
 */

package Auth

import (
	"LoongPanel/Panel/Service/Auth"
	"LoongPanel/Panel/Service/Database"
	"LoongPanel/Panel/Service/PanelLog"
	"errors"
	"github.com/gin-gonic/gin"
)

const (
	successCode      = 200
	errorCode        = 400
	unauthorizedCode = 401
)

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(errorCode, gin.H{"code": errorCode, "msg": "参数错误"})
		return
	}

	user, err := authenticateUser(req.Username, req.Email, req.Password)
	if err != nil || user == nil {
		if err != nil {
			PanelLog.DEBUG("[权限管理]", "登录失败: "+err.Error())
		} else {
			PanelLog.DEBUG("[权限管理]", "登录失败: 用户不存在")
		}
		c.JSON(unauthorizedCode, gin.H{"code": unauthorizedCode, "msg": "用户名或密码错误"})
		return
	}

	session, err := Auth.CreateSession(user.Name)
	if err != nil {
		PanelLog.DEBUG("[权限管理]", "登录失败: "+err.Error())
		c.JSON(errorCode, gin.H{"code": errorCode, "msg": "登录失败"})
		return
	}
	PanelLog.DEBUG("[权限管理]", user.Name+": 登录成功")
	Auth.SetSessionCookie(c, session)

	c.JSON(successCode, gin.H{
		"code": successCode,
		"msg":  "登录成功",
	})
}

func Logout(c *gin.Context) {
	session, err := c.Cookie(Auth.CookieName)
	if err != nil {
		c.JSON(errorCode, gin.H{"code": errorCode, "msg": "参数错误"})
		return
	}

	Auth.DeleteSession(session)
	PanelLog.DEBUG("[权限管理]", "用户注销成功")
	Auth.ClearSessionCookie(c)
	c.JSON(successCode, gin.H{"code": successCode, "msg": "注销成功"})
}

func authenticateUser(username, email, password string) (*Database.User, error) {
	users := Database.UserFind()
	for _, user := range users {
		if user.Name == username {
			// 通过用户名
			ok, err := Auth.AuthenticateUser(user.Name, password)
			if ok && err == nil {
				return &user, nil
			} else if !ok && err == nil {
				return nil, errors.New("密码错误")
			} else {
				return nil, err
			}
		} else if user.Mail == email {
			// 通过邮箱
			ok, err := Auth.AuthenticateUserByEmail(user.Mail, password)
			if ok && err == nil {
				return &user, nil
			} else if !ok && err == nil {
				return nil, errors.New("密码错误")
			} else {
				return nil, err
			}
		}
	}
	return nil, errors.New("用户名或密码错误")
}

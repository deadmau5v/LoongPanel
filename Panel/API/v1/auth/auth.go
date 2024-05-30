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
	"github.com/gin-gonic/gin"
)

var filed = map[string]interface{}{
	"code": 401,
	"msg":  "用户名或密码错误",
}

func Login(c *gin.Context) {
	// 读取请求的Json参数
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	err := c.BindJSON(&req)

	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		c.Abort()
		return
	}

	// 验证用户名和密码
	for _, user := range Database.UserFind() {
		// 使用邮箱登录
		if req.Email != "" {
			if user.Mail == req.Email && user.Password == req.Password {
				// 登录成功
				PanelLog.DEBUG(req.Email + ": 登录成功")
				c.JSON(200, gin.H{
					"code":    200,
					"msg":     "登录成功",
					"session": Auth.NewSESSION(user),
				})
				return
			} else {
				PanelLog.DEBUG(req.Email+": 登录失败", req.Password, req.Email, " != ", user.Password, user.Mail)
				continue
			}
		}
		// 使用用户名登录
		if req.Username != "" {
			if user.Name == req.Username && user.Password == req.Password {
				// 登录成功
				PanelLog.DEBUG(req.Username + ": 登录成功")
				c.JSON(200, gin.H{
					"code":    200,
					"msg":     "登录成功",
					"session": Auth.NewSESSION(user),
				})
				return
			} else {
				PanelLog.DEBUG(req.Username+": 登录失败", req.Password, req.Username, " != ", user.Password, user.Name)
				continue
			}
		}
	}

	// 登录失败
	PanelLog.DEBUG(req.Username + ": 登录失败")
	c.JSON(401, filed)
	return

}

func Logout(c *gin.Context) {
	// 读取请求的Json参数
	var req struct {
		Session string `json:"session"`
	}

	err := c.BindJSON(&req)

	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	return
}

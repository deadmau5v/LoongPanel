/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-18
 * 文件作用：使用Casbin进行权限验证
 */

package Auth

import (
	"LoongPanel/Panel/Service/Auth"
	"LoongPanel/Panel/Service/Database"
	"LoongPanel/Panel/Service/Log"
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
		if user.Name == req.Username && user.Password == req.Password {
			// 登录成功
			Log.DEBUG(req.Username + ": 登录成功")
			c.JSON(200, gin.H{
				"code":    200,
				"msg":     "登录成功",
				"session": Auth.RandomSESSION(req.Username),
			})
			c.Abort()
			return
		} else {
			Log.DEBUG(req.Username+": 登录失败", req.Password, " != ", user.Password)
			continue
		}
	}

	// 登录失败
	Log.DEBUG(req.Username + ": 登录失败")
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

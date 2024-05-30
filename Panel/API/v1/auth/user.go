/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-28
 * 文件作用：用户管理
 */

package Auth

import (
	"LoongPanel/Panel/Service/Auth"
	"LoongPanel/Panel/Service/Database"
	"LoongPanel/Panel/Service/PanelLog"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetUsers 获取用户列表
func GetUsers(ctx *gin.Context) {
	var Users []Database.User
	id, err := strconv.Atoi(ctx.Query("id"))
	if err == nil {
		PanelLog.DEBUG("查询用户ID: ", id)
		Database.DB.First(&Users, id)
	} else {
		PanelLog.DEBUG("查询所有用户")
		Database.DB.Find(&Users)
	}
	PanelLog.INFO("[权限管理] 获取用户列表")
	ctx.JSON(200, Users)
}

// CheckUserExist 检查用户是否存在
func CheckUserExist(user Database.User) gin.H {

	var count int64
	Database.DB.Model(&Database.User{}).Where("name = ?", user.Name).Count(&count)
	PanelLog.DEBUG("用户名重复数量: ", count)
	if count > 0 {
		return gin.H{"msg": "用户名已存在"}
	}
	Database.DB.Model(&Database.User{}).Where("email = ?", user.Mail).Count(&count)
	PanelLog.DEBUG("邮箱重复数量: ", count)
	if count > 0 {
		return gin.H{"msg": "邮箱已存在"}
	}
	return nil
}

// CheckUserUpdate 检查用户是否存在 用于更新
func CheckUserUpdate(user Database.User) gin.H {
	// 忽略自己
	var count int64
	Database.DB.Model(&Database.User{}).Where("name = ? AND id != ?", user.Name, user.ID).Count(&count)
	PanelLog.DEBUG("用户名重复数量: ", count)
	if count > 1 {
		return gin.H{"msg": "用户名已存在"}
	}
	Database.DB.Model(&Database.User{}).Where("email = ? AND id != ?", user.Mail, user.ID).Count(&count)
	PanelLog.DEBUG("邮箱重复数量: ", count)
	if count > 1 {
		return gin.H{"msg": "邮箱已存在"}
	}
	return nil
}

// CreateUser 创建用户
func CreateUser(ctx *gin.Context) {
	var user Database.User
	err := ctx.BindJSON(&user)

	if err != nil {
		PanelLog.ERROR("[权限管理] 绑定JSON失败", err)
		return
	}

	msg := CheckUserExist(user)
	if msg != nil {
		ctx.JSON(400, msg)
		return
	}

	PanelLog.INFO("[权限管理] 创建用户: ", user.Name)
	Database.DB.Create(&user)
	ctx.JSON(200, user)
}

// UpdateUser 更新用户
func UpdateUser(ctx *gin.Context) {
	var user Database.User
	var updateUser Database.User
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"msg": "ID错误"})
		return
	}
	Database.DB.First(&user, id)

	if user.ID == 0 && user.Name == "" {
		ctx.JSON(404, gin.H{"msg": "用户不存在"})
		return
	}

	err = ctx.BindJSON(&updateUser)
	if err != nil {
		PanelLog.ERROR("[权限管理] 绑定JSON失败", err)
		return
	}

	user.Name = updateUser.Name
	user.Mail = updateUser.Mail
	user.Password = updateUser.Password

	msg := CheckUserUpdate(user)
	if msg != nil {
		ctx.JSON(400, msg)
		return
	}

	Database.DB.Save(&user)
	PanelLog.INFO("[权限管理] 更新用户: ", user.Name)
	ctx.JSON(200, user)
}

// DeleteUser 删除用户
func DeleteUser(ctx *gin.Context) {
	var user Database.User
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"msg": "ID错误"})
		return
	}
	Database.DB.First(&user, id)
	Database.DB.Delete(&user)
	PanelLog.INFO("[权限管理] 删除用户: ", user.Name)
	ctx.JSON(200, gin.H{"msg": "删除成功"})
}

// GetRoles 获取角色列表
func GetRoles(ctx *gin.Context) {
	roles, err := Auth.Authenticator.GetAllRoles()
	if err != nil {
		PanelLog.ERROR("[权限管理] 获取角色列表失败", err)
		return
	}
	PanelLog.INFO("[权限管理] 获取角色列表")
	ctx.JSON(200, roles)
}

// CreateRole 创建角色
func CreateRole(ctx *gin.Context) {
	RoleName := ctx.Query("name")
	ok, err := Auth.Authenticator.AddRoleForUser(RoleName, RoleName)
	if err != nil {
		PanelLog.ERROR("[权限管理] 创建角色失败", err)
		return
	}
	PanelLog.INFO("[权限管理] 创建角色: ", RoleName)
	if ok {
		ctx.JSON(200, gin.H{
			"msg":    "创建成功",
			"status": 0,
		})
	} else {
		ctx.JSON(200, gin.H{
			"msg":    "创建失败",
			"status": 1,
		})
	}
}

// DeleteRole 删除角色
func DeleteRole(ctx *gin.Context) {
	RoleName := ctx.Query("name")
	ok, err := Auth.Authenticator.DeleteRole(RoleName)
	if err != nil {
		PanelLog.ERROR("[权限管理] 删除角色失败", err)
		return
	}
	PanelLog.INFO("[权限管理] 删除角色: ", RoleName)
	if ok {
		ctx.JSON(200, gin.H{
			"msg":    "删除成功",
			"status": 0,
		})
	} else {
		ctx.JSON(200, gin.H{
			"msg":    "删除失败",
			"status": 1,
		})
	}
}

type policy struct {
	Role   string `json:"role"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

// GetPolicy 获取权限列表
func GetPolicy(ctx *gin.Context) {
	var policy_list []policy

	Policy, err := Auth.Authenticator.GetPolicy()
	if err != nil {
		PanelLog.ERROR("[权限管理] 获取权限列表失败", err)
		return
	}

	for _, strings := range Policy {
		policy_list = append(policy_list, policy{
			Role:   strings[0],
			Path:   strings[1],
			Method: strings[2],
		})
	}

	PanelLog.INFO("[权限管理] 获取权限列表")
	ctx.JSON(200, policy_list)
}

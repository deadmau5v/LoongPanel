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
	"errors"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
	"unicode"
)

// GetUsers 获取用户列表
func GetUsers(ctx *gin.Context) {
	var Users []Database.User

	PanelLog.DEBUG("查询所有用户")
	Database.DB.Find(&Users)

	PanelLog.INFO("[权限管理] 获取用户列表")
	ctx.JSON(200, Users)
}

// DelUsers 批量删除用户
func DelUsers(ctx *gin.Context) {
	var ids []int
	err := ctx.BindJSON(&ids)
	if err != nil {
		PanelLog.ERROR("[权限管理] 绑定JSON失败", err)
		return
	}

	PanelLog.INFO("[权限管理] 批量删除用户: ", ids)
	Database.DB.Where("id IN (?)", ids).Delete(&Database.User{})
	ctx.JSON(200, gin.H{"msg": "删除成功", "status": 0})
}

// GetUser 获取用户
func GetUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		session, err := Auth.GetSessionByKey(ctx.GetHeader("Authorization"))
		if err != nil {
			PanelLog.ERROR("[权限管理] 获取SESSION失败", err)
			ctx.JSON(400, gin.H{"msg": err.Error(), "status": 1})
			return
		} else {
			PanelLog.DEBUG("[权限管理] 获取SESSION成功 使用默认ID: ", session.User.ID)
			id = session.User.ID
		}
	}

	var user Database.User
	Database.DB.Model(&Database.User{}).Where("id = ?", id).First(&user)
	PanelLog.INFO("[权限管理] 获取用户: ", user.Name)
	ctx.JSON(200, user)
}

// CheckUserExist 检查用户是否存在
func CheckUserExist(user Database.User) gin.H {
	var count int64
	err := CheckUserName(user.Name)
	if err != nil {
		return gin.H{"msg": err.Error(), "status": 1}
	}
	err = CheckMail(user.Mail)
	if err != nil {
		return gin.H{"msg": err.Error(), "status": 1}
	}
	Database.DB.Model(&Database.User{}).Where("name = ?", user.Name).Count(&count)
	PanelLog.DEBUG("用户名重复数量: ", count)
	if count > 0 {
		return gin.H{"msg": "用户名已存在", "status": 1}
	}
	Database.DB.Model(&Database.User{}).Where("mail = ?", user.Mail).Count(&count)
	PanelLog.DEBUG("邮箱重复数量: ", count)
	if count > 0 {
		return gin.H{"msg": "邮箱已存在", "status": 1}
	}
	return nil
}

// CheckUserUpdate 检查用户是否存在 用于更新
func CheckUserUpdate(user Database.User) gin.H {
	// 忽略自己
	var count int64
	err := CheckUserName(user.Name)
	if err != nil {
		return gin.H{"msg": err.Error(), "status": 1}
	}
	err = CheckMail(user.Mail)
	if err != nil {
		return gin.H{"msg": err.Error(), "status": 1}
	}
	Database.DB.Model(&Database.User{}).Where("name = ? AND id != ?", user.Name, user.ID).Count(&count)
	PanelLog.DEBUG("用户名重复数量: ", count)
	if count > 1 {
		return gin.H{"msg": "用户名已存在", "status": 1}
	}
	Database.DB.Model(&Database.User{}).Where("mail = ? AND id != ?", user.Mail, user.ID).Count(&count)
	PanelLog.DEBUG("邮箱重复数量: ", count)
	if count > 1 {
		return gin.H{"msg": "邮箱已存在", "status": 1}
	}
	return nil
}

// isChinese 检查是否为中文
func isChinese(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

// CheckUserName 检查用户名是否重复
func CheckUserName(name string) error {
	// 长度
	if len(name) < 3 || len(name) > 20 {
		return errors.New("用户名长度不合法")
	}

	// 字符
	for _, c := range name {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && !isChinese(c) {
			return errors.New("用户名包含非法字符")
		}

	}
	return nil
}

// CheckMail 检查邮箱格式是否合法
func CheckMail(mail string) error {
	if len(mail) < 5 || len(mail) > 50 {
		return errors.New("邮箱长度不合法")
	}
	// 正则检查邮箱是否合法
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regex.MatchString(mail) {
		return errors.New("邮箱格式不合法")
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
		// 更新SESSION获取到的用户ID
		sessionKey := ctx.GetHeader("Authorization")
		if sessionKey != "" {
			session, err := Auth.GetSessionByKey(sessionKey)
			if err != nil {
				PanelLog.ERROR("[权限管理] 获取SESSION失败", err)
				ctx.JSON(400, gin.H{"msg": err.Error()})
				return
			}
			id = session.User.ID
		} else {
			PanelLog.ERROR("[权限管理] ID错误")
			ctx.JSON(400, gin.H{"msg": "ID错误"})
			return
		}
	}
	Database.DB.Model(&Database.User{}).Where("id = ?", id).First(&user)

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
		PanelLog.ERROR("[权限管理] ID错误", err)
		ctx.JSON(400, gin.H{"msg": "ID错误", "status": 1})
		return
	}
	Database.DB.Model(&Database.User{}).Where("id = ?", id).First(&user)
	if user.Name == "admin" {
		ctx.JSON(401, gin.H{"msg": "无法删除内置用户", "status": 1})
		return
	}
	Database.DB.Delete(&user)
	PanelLog.INFO("[权限管理] 删除用户: ", user.Name)
	ctx.JSON(200, gin.H{"msg": "删除成功", "status": 0})
}

type role struct {
	Name       string   `json:"name"`
	PolicyList []policy `json:"policy_list"`
}

// GetRoles 获取角色列表
func GetRoles(ctx *gin.Context) {
	roles, err := Auth.Authenticator.GetAllRoles()
	if err != nil {
		PanelLog.ERROR("[权限管理] 获取角色列表失败", err)
		return
	}
	PanelLog.INFO("[权限管理] 获取角色列表")
	roles_ := make([]role, 0)

	for _, r := range roles {
		role_ := role{
			Name:       r,
			PolicyList: getPolicy(r),
		}
		roles_ = append(roles_, role_)
	}

	ctx.JSON(200, roles_)
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
	if RoleName == "admin" || RoleName == "user" {
		ctx.JSON(401, gin.H{
			"msg":    "无法删除内置角色",
			"status": 1,
		})
		return
	}

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

// getPolicy 获取权限列表
func getPolicy(role string) []policy {
	policyList := make([]policy, 0)

	allPolicy, err := Auth.Authenticator.GetPolicy()
	for _, policy_ := range allPolicy {
		if policy_[0] == role {
			policyList = append(policyList, policy{
				Role:   policy_[0],
				Path:   policy_[1],
				Method: policy_[2],
			})
		}
	}

	if err != nil {
		PanelLog.ERROR("[权限管理] 获取权限列表失败", err)
		return nil
	}

	PanelLog.INFO("[权限管理] 获取权限列表")
	return policyList
}

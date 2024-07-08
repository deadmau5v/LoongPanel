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
	"net/http"
	"regexp"
	"strconv"
	"unicode"

	"github.com/gin-gonic/gin"
)

// GetUsers 获取用户列表
func GetUsers(ctx *gin.Context) {
	var users []Database.User
	result := Database.DB.Model(&Database.User{}).Find(&users)
	if result.Error != nil {
		PanelLog.ERROR("[权限管理] 获取用户列表失败", result.Error)
		ctx.JSON(500, gin.H{"status": 1, "msg": "获取用户列表失败"})
		return
	}

	PanelLog.INFO("[权限管理] 获取用户列表")
	ctx.JSON(200, gin.H{"status": 0, "data": users})
}

// DelUsers 批量删除用户
func DelUsers(ctx *gin.Context) {
	var ids []int
	if err := ctx.BindJSON(&ids); err != nil {
		PanelLog.ERROR("[权限管理] 绑定JSON失败", err)
		ctx.JSON(400, gin.H{"status": 1, "msg": "无效的请求数据"})
		return
	}

	result := Database.DB.Where("id IN ? AND name != 'admin'", ids).Delete(&Database.User{})
	if result.Error != nil {
		PanelLog.ERROR("[权限管理] 批量删除用户失败", result.Error)
		ctx.JSON(500, gin.H{"status": 1, "msg": "删除用户失败"})
		return
	}

	PanelLog.INFO("[权限管理] 批量删除用户: ", ids)
	ctx.JSON(200, gin.H{"status": 0, "msg": "删除成功", "count": result.RowsAffected})
}

// GetUser 获取用户
func GetUser(ctx *gin.Context) {
	id, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	var user Database.User
	result := Database.DB.Model(&Database.User{}).Where("id = ?", id).First(&user)
	if result.Error != nil {
		PanelLog.ERROR("[权限管理] 获取用户失败", result.Error)
		ctx.JSON(404, gin.H{"status": 1, "msg": "用户不存在"})
		return
	}

	PanelLog.INFO("[权限管理] 获取用户: ", user.Name)
	ctx.JSON(200, gin.H{"status": 0, "data": user})
}

// CreateUser 创建用户
func CreateUser(ctx *gin.Context) {
	var user Database.User
	if err := ctx.BindJSON(&user); err != nil {
		PanelLog.ERROR("[权限管理] 绑定JSON失败", err)
		ctx.JSON(400, gin.H{"status": 1, "msg": "无效的请求数据"})
		return
	}

	if msg := CheckUserExist(user); msg != nil {
		ctx.JSON(400, msg)
		return
	}

	err := Auth.CreateUser(user)
	if err != nil {
		return
	}

	PanelLog.INFO("[权限管理] 创建用户: ", user.Name)
}

// UpdateUser 更新用户
func UpdateUser(ctx *gin.Context) {
	id, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"status": 1, "msg": err.Error()})
		return
	}

	var updateUser Database.User
	if err := ctx.BindJSON(&updateUser); err != nil {
		PanelLog.ERROR("[权限管理] 绑定JSON失败", err)
		ctx.JSON(400, gin.H{"status": 1, "msg": "无效的请求数据"})
		return
	}

	updateUser.ID = id
	if msg := CheckUserUpdate(updateUser); msg != nil {
		ctx.JSON(400, msg)
		return
	}

	result := Database.DB.Model(&Database.User{}).Where("id = ?", id).Updates(Database.User{
		Name:     updateUser.Name,
		Mail:     updateUser.Mail,
		Password: updateUser.Password,
	})

	if result.Error != nil {
		PanelLog.ERROR("[权限管理] 更新用户失败", result.Error)
		ctx.JSON(500, gin.H{"status": 1, "msg": "更新用户失败"})
		return
	}

	PanelLog.INFO("[权限管理] 更新用户: ", updateUser.Name)
	ctx.JSON(200, gin.H{"status": 0, "msg": "更新成功", "data": updateUser})
}

// DeleteUser 删除用户
func DeleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		PanelLog.ERROR("[权限管理] ID错误", err)
		ctx.JSON(400, gin.H{"status": 1, "msg": "无效的用户ID"})
		return
	}

	var user Database.User
	result := Database.DB.Where("id = ? AND name != 'admin'", id).Delete(&user)
	if result.Error != nil {
		PanelLog.ERROR("[权限管理] 删除用户失败", result.Error)
		ctx.JSON(500, gin.H{"status": 1, "msg": "删除用户失败"})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"status": 1, "msg": "用户不存在或无法删除"})
		return
	}

	PanelLog.INFO("[权限管理] 删除用户: ", id)
	ctx.JSON(200, gin.H{"status": 0, "msg": "删除成功"})
}

// 辅助函数

func getUserID(ctx *gin.Context) (int, error) {

	username, exists := ctx.Get("username")
	if !exists {
		return 0, errors.New("无效的用户ID")
	}

	user := Database.User{}
	result := Database.DB.Model(&Database.User{}).Where("name = ?", username).First(&user)
	if result.Error != nil {
		return 0, errors.New("无效的用户ID")
	}

	return user.ID, nil
}

func CheckUserExist(user Database.User) gin.H {
	if err := CheckUserName(user.Name); err != nil {
		return gin.H{"status": 1, "msg": err.Error()}
	}
	if err := CheckMail(user.Mail); err != nil {
		return gin.H{"status": 1, "msg": err.Error()}
	}

	var count int64
	Database.DB.Model(&Database.User{}).Where("name = ? OR mail = ?", user.Name, user.Mail).Count(&count)
	if count > 0 {
		return gin.H{"status": 1, "msg": "用户名或邮箱已存在"}
	}
	return nil
}

func CheckUserUpdate(user Database.User) gin.H {
	if err := CheckUserName(user.Name); err != nil {
		return gin.H{"status": 1, "msg": err.Error()}
	}
	if err := CheckMail(user.Mail); err != nil {
		return gin.H{"status": 1, "msg": err.Error()}
	}

	var count int64
	Database.DB.Model(&Database.User{}).Where("(name = ? OR mail = ?) AND id != ?", user.Name, user.Mail, user.ID).Count(&count)
	if count > 0 {
		return gin.H{"status": 1, "msg": "用户名或邮箱已存在"}
	}
	return nil
}

func CheckUserName(name string) error {
	if len(name) < 3 || len(name) > 20 {
		return errors.New("用户名长度应在3-20之间")
	}
	for _, c := range name {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) && !isChinese(c) {
			return errors.New("用户名只能包含字母、数字和中文")
		}
	}
	return nil
}

func CheckMail(mail string) error {
	if len(mail) < 5 || len(mail) > 50 {
		return errors.New("邮箱长度应在5-50之间")
	}
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regex.MatchString(mail) {
		return errors.New("邮箱格式不正确")
	}
	return nil
}

func isChinese(r rune) bool {
	return unicode.Is(unicode.Han, r)
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

// AddPolicy 添加策略
func AddPolicy(c *gin.Context) {
	req := policy{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(errorCode, gin.H{"status": 1, "msg": "参数错误"})
		return
	}

	if req.Role == "" || req.Path == "" || req.Method == "" {
		c.JSON(errorCode, gin.H{"status": 1, "msg": "参数错误"})
		return
	}

	exists, err := Auth.Authenticator.AddPolicy(req.Role, req.Path, req.Method)
	if err != nil {
		c.JSON(errorCode, gin.H{"status": 1, "msg": "添加失败"})
		return
	}
	if exists {
		c.JSON(errorCode, gin.H{"status": 1, "msg": "已存在"})
		return
	} else {
		PanelLog.DEBUG("[权限管理]", "添加权限: "+req.Role+" "+req.Path+" "+req.Method)
	}

	c.JSON(successCode, gin.H{"status": 0, "msg": "添加成功"})
}

// DeletePolicy 删除策略
func DeletePolicy(c *gin.Context) {
	req := policy{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(errorCode, gin.H{"status": 1, "msg": "参数错误"})
		return
	}

	ok, err := Auth.Authenticator.RemovePolicy(req.Role, req.Path, req.Method)
	if err != nil {
		c.JSON(errorCode, gin.H{"status": 1, "msg": "删除失败"})
		return
	}
	if ok {
		PanelLog.DEBUG("[权限管理]", "删除权限: "+req.Role+" "+req.Path+" "+req.Method)
		c.JSON(successCode, gin.H{"status": 0, "msg": "删除成功"})
	} else {
		c.JSON(errorCode, gin.H{"status": 1, "msg": "删除失败"})
	}
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	req := struct {
		Password string `json:"password"`
	}{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(errorCode, gin.H{"status": 1, "msg": "参数错误"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(errorCode, gin.H{"status": 1, "msg": "无效的用户ID"})
		return
	}
	user := Database.User{}
	result := Database.DB.Model(&Database.User{}).Where("name = ?", username).First(&user)
	if result.Error != nil {
		c.JSON(errorCode, gin.H{"status": 1, "msg": "无效的用户ID"})
		return
	}

	user.Password = Auth.HashPassword(req.Password)
	err := Auth.ValidateCredentials(user.Name, req.Password)
	if err != nil {
		c.JSON(errorCode, gin.H{"status": 1, "msg": "密码格式错误" + err.Error()})
		return
	}
	user.Update()
	PanelLog.INFO("[权限管理]", user.Name, "用户修改密码")
	c.JSON(successCode, gin.H{"status": 0, "msg": "修改成功"})
}

// Register 注册
func Register(ctx *gin.Context) {
	req := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Mail     string `json:"mail"`
	}{}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": "参数错误"})
		return
	}

	if len(req.Username) < 5 || len(req.Password) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": 1, "msg": "用户名或密码长度不符合要求"})
		return
	}

	hashedPassword := Auth.HashPassword(req.Password)
	newUser := Database.User{
		Name:     req.Username,
		Password: hashedPassword,
		Mail:     req.Mail,
	}

	result := Database.DB.Create(&newUser)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": 1, "msg": "注册失败"})
		return
	}

	PanelLog.INFO("[用户管理]", "注册新用户:", req.Username)
	ctx.JSON(http.StatusOK, gin.H{"status": 0, "msg": "注册成功"})
}

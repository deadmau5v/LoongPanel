package Auth

import (
	"LoongPanel/Panel/Service/Database"
	"LoongPanel/Panel/Service/PanelLog"
	"LoongPanel/Panel/Service/System"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"os"
	"strings"
)

func init() {
	// 初始化casbin模型
	stat, err := os.Stat(System.WORKDIR + "/resource/model.conf")
	if err != nil || stat.IsDir() || stat.Size() == 0 {
		_ = os.Mkdir(System.WORKDIR+"/resource", os.ModePerm)
		File, err := os.Create(System.WORKDIR + "/resource/model.conf")
		if err != nil {
			panic(err)
		}
		_, _ = File.WriteString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
	}

	// 初始化casbin
	a, _ := gormadapter.NewAdapterByDB(Database.DB)
	Authenticator, _ = casbin.NewEnforcer(System.WORKDIR+"/resource/model.conf", a)
	policy, err := Authenticator.GetPolicy()

	// 获取策略组 g
	groups, _ := Authenticator.GetGroupingPolicy()
	if len(groups) == 0 {
		PanelLog.INFO("[权限管理] 初始化策略组")
		_, _ = Authenticator.AddGroupingPolicy("admin", "admin")
		_, _ = Authenticator.AddGroupingPolicy("admin", "user")
		_, _ = Authenticator.AddGroupingPolicy("user", "user")
	} else {
		for _, group := range groups {
			PanelLog.DEBUG("[权限管理] 策略组: ", group)
		}
	}

	Authenticator.EnableAutoSave(true)
	if err != nil {
		return
	}
	// 如果没有权限策略，添加默认策略
	if len(policy) != 0 {
		for _, v := range policy {
			msg := strings.Join(v, " ")
			PanelLog.DEBUG("[权限管理] 策略 ", msg, " 已加载")
		}
		PanelLog.INFO("[权限管理] 策略已加载完成")
	}

	// 生成默认用户
	haveAdmin := Database.DB.Where("name = ?", "admin").Find(&Database.User{}).RowsAffected

	if haveAdmin == 0 {
		admin := Database.User{
			Name:     "admin",
			Password: "12345678",
			Role:     "admin",
		}
		user := Database.User{
			Name:     "user",
			Password: "12345678",
			Role:     "user",
		}

		err := CreateUser(admin)
		if err != nil {
			PanelLog.ERROR("[权限管理]", "默认用户创建失败", err.Error())
			return
		}
		err = CreateUser(user)
		if err != nil {
			PanelLog.ERROR("[权限管理]", "默认用户创建失败", err.Error())
			return
		}

		PanelLog.INFO("[权限管理] 默认用户 admin 创建成功")
	}
}

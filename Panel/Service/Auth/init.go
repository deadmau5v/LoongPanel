/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-16
 * 文件作用：初始化casbin权限管理
 */

package Auth

import (
	"LoongPanel/Panel/Service/Database"
	"LoongPanel/Panel/Service/Log"
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
	Authenticator.EnableAutoSave(true)
	if err != nil {
		return
	}
	// 如果没有权限策略，添加默认策略
	if len(policy) == 0 {
		// 系统监控 默认管理员权限、用户权限

		_, err = Authenticator.AddPolicy("admin", "/api/v1/status/system_status", "GET")
		_, err = Authenticator.AddPolicy("admin", "/api/v1/status/system_info", "GET")
		_, err = Authenticator.AddPolicy("admin", "/api/v1/status/disks", "GET")
		_, err = Authenticator.AddPolicy("user", "/api/v1/status/system_status", "GET")
		_, err = Authenticator.AddPolicy("user", "/api/v1/status/system_info", "GET")
		_, err = Authenticator.AddPolicy("user", "/api/v1/status/disks", "GET")
		// 清理垃圾 默认管理员权限
		_, err = Authenticator.AddPolicy("admin", "/api/v1/clean/pkg_auto_clean", "GET")
		// 电源操作 默认管理员权限
		_, err = Authenticator.AddPolicy("admin", "/api/v1/power/shutdown", "GET")
		_, err = Authenticator.AddPolicy("admin", "/api/v1/power/reboot", "GET")
		// 文件操作 默认管理员权限、用户权限
		_, err = Authenticator.AddPolicy("admin", "/api/v1/files/dir", "GET")
		_, err = Authenticator.AddPolicy("user", "/api/v1/files/dir", "GET")
		// 终端操作 默认管理员权限
		_, err = Authenticator.AddPolicy("admin", "/api/v1/screen/input", "GET")
		_, err = Authenticator.AddPolicy("admin", "/api/v1/screen/create", "GET")
		_, err = Authenticator.AddPolicy("admin", "/api/v1/screen/close", "GET")
		_, err = Authenticator.AddPolicy("admin", "/api/v1/screen/output", "GET")
		_, err = Authenticator.AddPolicy("admin", "/api/v1/screen/get_screens", "GET")
		_, err = Authenticator.AddPolicy("admin", "/api/ws/screen", "GET")
		// Ping
		_, err = Authenticator.AddPolicy("admin", "/api/v1/ping", "GET")
		_, err = Authenticator.AddPolicy("user", "/api/v1/ping", "GET")
		// 登录
		_, err = Authenticator.AddPolicy("admin", "/api/v1/auth/login", "POST")
		_, err = Authenticator.AddPolicy("user", "/api/v1/auth/login", "POST")
		// 默认角色
		_, err = Authenticator.AddGroupingPolicy("admin", "admin")
		_, err = Authenticator.AddGroupingPolicy("user", "user")
		err := Authenticator.SavePolicy()
		if err != nil {
			Log.ERROR("[权限管理] 添加默认策略错误", err.Error())
			return
		}
	} else {
		Log.INFO("[权限管理] 策略已加载")

		for _, v := range policy {
			msg := strings.Join(v, " ")
			Log.DEBUG("[权限管理] 策略 [", msg, "]")
		}
	}

	// 测试可删除
	ok, err := Authenticator.Enforce("admin", "/api/v1/auth/login", "POST")
	if err != nil {
		Log.ERROR("[权限管理] 初始化错误", err.Error())
		return
	}

	if ok {
		Log.INFO("[权限管理]Casbin初始化成功")
	}

}

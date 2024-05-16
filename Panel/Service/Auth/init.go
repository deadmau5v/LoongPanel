/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-16
 * 文件作用：
 */

package Auth

import (
	"LoongPanel/Panel/Service/Database"
	"LoongPanel/Panel/Service/Log"
	"LoongPanel/Panel/Service/System"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"os"
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
	Authenticator, _ := casbin.NewEnforcer(System.WORKDIR+"/resource/model.conf", a)

	_, err = Authenticator.AddPolicy("admin", "/api/v1/auth/login", "POST")
	if err != nil {
		Log.ERROR(err.Error())
		return
	}

	// 测试可删除
	ok, err := Authenticator.Enforce("admin", "/api/v1/auth/login", "POST")
	if err != nil {
		return
	}

	if ok {
		Log.INFO("[权限管理]Casbin初始化成功")
	}

}

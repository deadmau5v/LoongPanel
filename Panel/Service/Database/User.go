/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-25
 * 文件作用：
 */

package Database

import (
	"LoongPanel/Panel/Service/Log"
)

func (user *User) Save() {
	// 保存到数据库
	DB.Create(&user)

}

func (user *User) Delete() {
	DB.Delete(&user)
}

func UserFind() []User {
	var Users []User
	DB.Find(&Users)

	return Users
}

func init() {
	err := DB.AutoMigrate(&User{})
	if err != nil {
		Log.ERROR("初始化SQLite数据库失败")
		return
	}

	// 初始化管理员
	if len(UserFind()) == 0 {
		Log.INFO("[数据库模块]初始化管理员")
		admin := User{
			Name:     "admin",
			Password: "123456",
			Role:     "admin",
		}
		admin.Save()
	}
	// 初始化用户
	if len(UserFind()) == 1 {
		Log.INFO("[数据库模块]初始化用户")
		user := User{
			Name:     "user",
			Password: "123456",
			Role:     "user",
		}
		user.Save()
	}
}

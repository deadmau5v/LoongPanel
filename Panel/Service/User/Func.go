/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：面板用户模块
 */

package User

import (
	"LoongPanel/Panel/Service/Database"
	"LoongPanel/Panel/Service/Log"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (user *User) Save() {
	// 保存到数据库
	Database.DB.Create(&user)

}

func (user *User) Delete() {
	Database.DB.Delete(&user)
}

func Find() []User {
	var Users []User
	Database.DB.Find(&Users)

	return Users
}

func init() {
	err := Database.DB.AutoMigrate(&User{})
	if err != nil {
		Log.ERROR("初始化SQLite数据库失败")
		return
	}

	// 初始化管理员
	if len(Find()) == 0 {
		Log.INFO("初始化管理员")
		admin := User{
			Name:     "admin",
			Password: "123456",
			Role:     "admin",
		}
		admin.Save()
	}
	// 初始化用户
	if len(Find()) == 1 {
		Log.INFO("初始化用户")
		user := User{
			Name:     "user",
			Password: "123456",
			Role:     "user",
		}
		user.Save()
	}
}

/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：面板用户模块
 */

package User

import (
	"LoongPanel/Panel/Service/Database"
	"fmt"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (user *User) Save() {
	// 保存到数据库
	switch Database.UseDB {
	case "sqlite":
		Database.DB.Create(&user)
	}
}

func (user *User) Delete() {
	switch Database.UseDB {
	case "sqlite":
		Database.DB.Delete(&user)
	}
}

func Find(where map[string]interface{}) []User {
	var Users = make([]User, 0)

	switch Database.UseDB {
	case "sqlite":
		if where != nil {
			query := ""
			keys := make([]interface{}, 0)
			for k, v := range where {
				query = query + fmt.Sprintf(" %v = ? ", k)
				keys = append(keys, v)
			}
			args := make([]interface{}, 0)
			args = append(args, query)
			args = append(args, keys...)
			Database.DB.Find(&Users, args)
		} else {
			Database.DB.Find(&Users)
		}
	}
	return Users
}

func init() {
	err := Database.DB.AutoMigrate(&User{})
	if err != nil {
		fmt.Println("初始化SQLite数据库失败")
		return
	}
}

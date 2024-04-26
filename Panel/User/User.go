package User

import (
	Database2 "LoongPanel/Panel/Database"
	"fmt"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (user *User) Save() {
	// 保存到数据库
	switch Database2.UseDB {
	case "sqlite":
		Database2.SQLite.Create(&user)
	case "mysql":
		// Todo 数据库操作兼容
	case "redis":
		// Todo 数据库操作兼容
	}
}

func (user *User) Delete() {
	switch Database2.UseDB {
	case "sqlite":
		Database2.SQLite.Delete(&user)
	case "mysql":
		// Todo 数据库操作兼容
	case "redis":
		// Todo 数据库操作兼容
	}
}

func Find(where map[string]interface{}) []User {
	var Users = make([]User, 0)

	switch Database2.UseDB {
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
			Database2.SQLite.Find(&Users, args)
		} else {
			Database2.SQLite.Find(&Users)
		}
	case "mysql":
		// Todo 数据库操作兼容
	case "redis":
		// Todo 数据库操作兼容
	}
	return Users
}

func init() {
	err := Database2.SQLite.AutoMigrate(&User{})
	if err != nil {
		fmt.Println("初始化SQLite数据库失败")
		return
	}
}

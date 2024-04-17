package Database

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var UseDB = "sqlite"
var Address = "localhost"
var User = "root"
var Password = ""
var UseDatabase = "LoongPanel"

func init() {
	err := errors.New("None")

	switch UseDB {
	case "mysql":
		// Todo Mysql支持
	case "redis":
		// Todo Redis 支持
	case "sqlite":
		SQLite, err = gorm.Open(sqlite.Open("LoongPanel.db"))
	}

	if err != nil {
		fmt.Printf("使用的数据库: %v 使用的库: %v 详细: \n", UseDB, UseDatabase)
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println("连接成功")
	}
}

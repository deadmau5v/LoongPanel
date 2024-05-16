/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：数据库配置
 */

package Database

import (
	"LoongPanel/Panel/Service/Log"
	"errors"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var UseDB = "sqlite"
var Address = "localhost"
var User = "root"
var Password = ""
var UseDatabase = "LoongPanel"

func init() {
	err := errors.New("none")
	_ = os.Mkdir("/resource", os.ModePerm)
	switch UseDB {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open("/resource/LoongPanel.db"))
	}

	if err != nil {
		Log.INFO("使用的数据库: %v 使用的库: %v 详细: \n", UseDB, UseDatabase)
		Log.ERROR(err.Error())
		return
	} else {
		Log.INFO("[数据库模块]连接成功")
	}
}

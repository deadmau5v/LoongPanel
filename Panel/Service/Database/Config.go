/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：数据库配置
 */

package Database

import (
	"LoongPanel/Panel/Service/Log"
	"LoongPanel/Panel/Service/System"
	"errors"
	"os"
	"path"

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
	_ = os.Mkdir(path.Join(System.WORKDIR, "resource"), os.ModePerm)

	switch UseDB {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(path.Join(System.WORKDIR, "resource", "LoongPanel.db")))
	}

	if err != nil {
		Log.ERROR("使用的数据库: ", UseDB, "使用的库:", UseDatabase, "详细: \n")
		Log.ERROR(err.Error())
		return
	} else {
		Log.INFO("[数据库模块]连接成功")
	}
}

/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：数据库配置
 */

package Database

import (
	"LoongPanel/Panel/Service/PanelLog"
	"LoongPanel/Panel/Service/System"
	"errors"
	"flag"
	"fmt"
	"gorm.io/driver/mysql"
	"os"
	"path"
	"strings"

	"gorm.io/gorm"
)

var (
	dbAddress     = "localhost"
	dbUser        = "root"
	dbPassword    = ""
	dbUseDatabase = "LoongPanel" // 使用的数据库名
	dbPort        = 4000
)

func init() {

	flag.StringVar(&dbAddress, "a", dbAddress, "数据库地址")
	flag.StringVar(&dbUser, "u", dbUser, "数据库用户名")
	flag.StringVar(&dbPassword, "p", dbPassword, "数据库密码")
	flag.StringVar(&dbUseDatabase, "d", dbUseDatabase, "数据库名")
	flag.IntVar(&dbPort, "P", dbPort, "数据库端口")

	flag.Parse()

	// 创建资源目录
	err := errors.New("none")
	_ = os.Mkdir(path.Join(System.WORKDIR, "resource"), os.ModePerm)

	// 配置 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbAddress, dbPort)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		if strings.Contains(err.Error(), "Unknown database") {
			PanelLog.ERROR("[数据库] 数据库不存在")
		} else if strings.Contains(err.Error(), "Access denied") {
			PanelLog.ERROR("[数据库] 用户名或密码错误")
		} else {
			PanelLog.ERROR("[数据库] 连接数据库失败")
		}
		panic(err.Error())
		// Todo 自动下载TIDB运行环境
	} else {
		PanelLog.INFO("[数据库] 连接成功")
	}

	// 创建 LoongPanel使用的数据库
	DB.Exec("CREATE DATABASE IF NOT EXISTS " + dbUseDatabase)
	// 重新链接
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbAddress, dbPort, dbUseDatabase)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

}

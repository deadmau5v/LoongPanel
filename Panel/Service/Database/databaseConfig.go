/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：数据库配置
 */

package Database

import (
	"LoongPanel/Panel/Service/PanelLog"
	"LoongPanel/Panel/Service/System"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

type config struct {
	DbAddress     string `json:"db_address"`
	DbUser        string `json:"db_user"`
	DbPassword    string `json:"db_password"`
	DbUseDatabase string `json:"db_use_database"`
	DbPort        int    `json:"db_port"`
}

var (
	DB     *gorm.DB
	Config *config
)

func Download(url string) {
	resp, err := http.Get(url)
	if err != nil {
		PanelLog.ERROR("[数据库] 下载失败: " + err.Error())
		return
	}
	defer resp.Body.Close()

	// 获取文件名
	urlParts := strings.Split(url, "/")
	fileName := urlParts[len(urlParts)-1]
	filePath := path.Join(System.WORKDIR, "resource", fileName)

	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		PanelLog.ERROR("[数据库] 文件创建失败: " + err.Error())
		return
	}
	defer out.Close()

	// 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		PanelLog.ERROR("[数据库] 文件写入失败: " + err.Error())
		return
	}

	PanelLog.INFO("[数据库] 下载成功: " + filePath)
}

// GetTIDB 自动下载TiDB
func GetTIDB() {
	f, err := os.Stat(path.Join(System.WORKDIR, "resource", "tidb-server"))
	if err == nil && !f.IsDir() {
		return
	}

	PanelLog.INFO("[数据库] 首次启动需要自动下载TiDB... 请稍等")
	// 检查架构
	if System.Data.OSArch == "amd64" {
		Download("https://cdn1.d5v.cc/CDN/Project/LoongPanel/applications/tidb-server-amd64")
		os.Rename(path.Join(System.WORKDIR, "resource", "tidb-server-amd64"), path.Join(System.WORKDIR, "resource", "tidb-server"))
	} else {
		Download("https://cdn1.d5v.cc/CDN/Project/LoongPanel/applications/tidb-server-loong64")
		os.Rename(path.Join(System.WORKDIR, "resource", "tidb-server-loong64"), path.Join(System.WORKDIR, "resource", "tidb-server"))
	}

	// 设置运行权限
	os.Chmod(path.Join(System.WORKDIR, "resource", "tidb-server"), 0755)
}

func TIDBStatus() bool {
	pid, err := exec.Command("pgrep", "-f", "tidb-server").Output()
	if err != nil {
		return false
	} else if len(pid) == 0 {
		return false
	} else {
		return true
	}
}

func RunTIDB() {
	if TIDBStatus() {
		return
	}
	PanelLog.INFO("[数据库] 启动TiDB...")
	cmd := exec.Command(path.Join(System.WORKDIR, "resource", "tidb-server"), "start")
	cmd.Start()
}

// SaveFile 保存文件
func SaveConfig() {
	configData, err := json.Marshal(Config)
	if err != nil {
		PanelLog.ERROR("[数据库] 配置文件保存失败: " + err.Error())
		return
	}
	os.WriteFile(path.Join(System.WORKDIR, "resource", "config.json"), configData, 0644)
}

// LoadConfig 加载配置文件
func LoadConfig() {
	Config = &config{}
	configPath := path.Join(System.WORKDIR, "resource", "config.json")
	_, err := os.Stat(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			PanelLog.INFO("[数据库] 首次启动, 自动创建配置文件...")
			GetTIDB()
			// 初始化TIDB
			Config.DbAddress = "127.0.0.1"
			Config.DbPort = 4000
			Config.DbUser = "root"
			Config.DbPassword = GenerateRandomPassword()
			Config.DbUseDatabase = "LoongPanel"
			RunTIDB()
			dsn := fmt.Sprintf("%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
				Config.DbUser, Config.DbAddress, Config.DbPort)
			time.Sleep(2 * time.Second)
			for {
				// 等待启动完成
				DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				if err != nil {
					if strings.Contains(err.Error(), "Access denied") {
						PanelLog.ERROR("[数据库] 用户名或密码错误")
						exec.Command("killall", "tidb-server").Start()
						err = os.RemoveAll("/tmp/tidb")
						if err != nil {
							PanelLog.ERROR("[数据库] 删除TiDB数据目录失败: " + err.Error())
						}
						RunTIDB()
						continue
					} else {
						time.Sleep(1 * time.Second)
					}
				}
				break
			}
			// 设置密码
			// Start of Selection
			DB.Exec("ALTER USER 'root'@'%' IDENTIFIED BY '" + Config.DbPassword + "'")
			DB.Exec("CREATE DATABASE IF NOT EXISTS " + Config.DbUseDatabase)
			SaveConfig()
		} else {
			PanelLog.ERROR("[数据库] 配置文件读取失败: " + err.Error())
			panic(err.Error())
		}
	}
	configData, err := os.ReadFile(configPath)
	if err != nil {
		PanelLog.ERROR("[数据库] 配置文件读取失败: " + err.Error())
		panic(err.Error())
	}
	err = json.Unmarshal(configData, Config)
	if err != nil {
		PanelLog.ERROR("[数据库] 配置文件解析失败: " + err.Error())
		panic(err.Error())
	}
	PanelLog.DEBUG("[数据库] 配置文件加载成功: " + Config.DbPassword)
}

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		PanelLog.ERROR("[数据库] 生成随机密码失败")
		return ""
	}
	return hex.EncodeToString(bytes)
}

func Connect() {
	var err error

	LoadConfig()
	RunTIDB()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Config.DbUser, Config.DbPassword, Config.DbAddress, Config.DbPort, Config.DbUseDatabase)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		if strings.Contains(err.Error(), "Unknown database") {
			PanelLog.DEBUG("[数据库] 数据库不存在")
			os.RemoveAll(path.Join(System.WORKDIR, "resource", "config.json"))
		} else if strings.Contains(err.Error(), "Access denied") {
			PanelLog.DEBUG("[数据库] 用户名或密码错误, 删除配置文件并重新连接...")
			os.RemoveAll(path.Join(System.WORKDIR, "resource", "config.json"))
		} else {
			PanelLog.DEBUG("[数据库] 数据库启动中... 请稍等")
			time.Sleep(1 * time.Second)
		}
		Connect()
	} else {
		PanelLog.INFO("[数据库] 连接成功")
	}
}

func init() {
	// 创建资源目录
	_ = os.Mkdir(path.Join(System.WORKDIR, "resource"), os.ModePerm)

	Connect()
}

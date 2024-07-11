package config

import (
	"LoongPanel/Panel/Service/Database"
	"LoongPanel/Panel/Service/PanelLog"
	"time"
)

var (
	Config ConfigStruct
)

type ConfigStruct struct {
	// 病毒扫描配置
	Clamav struct {
		// 是否开启定时扫描
		CronScan bool `json:"cron_scan"`
		// 定时扫描时间间隔
		CronScanTime time.Duration `json:"cron_scan_time"`
	} `gorm:"embedded" json:"clamav"` // 使用 gorm:"embedded" 来嵌入结构体

	// 权限管理配置
	Auth struct {
		// 允许注册
		CanRegister bool `json:"can_register"` // 修正字段名拼写错误
		// 允许mail登录
		CanMailLogin bool `json:"can_mail_login"`
	} `gorm:"embedded" json:"auth"` // 使用 gorm:"embedded" 来嵌入结构体

	// 面板日志
	PanelLog struct {
		// 是否开启调试模式
		Debug bool `json:"debug"`
		// 是否保存到文件
		SaveToFile bool `json:"save_to_file"`
	} `gorm:"embedded" json:"panel_log"` // 使用 gorm:"embedded" 来嵌入结构体

	MailConfig struct {
		MailHost  string `json:"mail_host" valid:"required~邮件服务器未配置,dns"`   // 邮件服务器
		MailPort  int    `json:"mail_port" valid:"required~邮件端口未配置,port"`   // 邮件端口
		MailUser  string `json:"mail_user" valid:"required~邮件用户名未配置"`       // 邮件用户名
		MailPass  string `json:"mail_pass" valid:"required~邮件密码未配置"`        // 邮件密码
		MailFrom  string `json:"mail_from" valid:"required~邮件发送者未配置,email"` // 邮件发送者
		MailSSL   bool   `json:"mail_ssl" valid:"required~邮件SSL未配置"`        // 邮件SSL
		MailTo    string `json:"mail_to" valid:"required~邮件接收者未配置,email"`   // 邮件接收者
		MailBody  string `json:"mail_body" valid:"required~邮件内容未配置"`        // 邮件内容
		MailTitle string `json:"mail_title" valid:"required~邮件标题未配置"`       // 邮件标题
	} `gorm:"embedded" json:"mail_config"`
}

// 保存或更新配置
func (c *ConfigStruct) Save() error {
	if err := Database.DB.Model(&ConfigStruct{}).Updates(c).Error; err != nil {
		return err
	}
	return nil
}

// 初始化配置
func init() {
	Database.DB.AutoMigrate(&ConfigStruct{})

	var count int64
	Database.DB.Model(&ConfigStruct{}).Count(&count)
	if count == 0 {
		// 默认配置
		Config.Clamav.CronScan = false
		Config.Clamav.CronScanTime = time.Hour * 24 // 默认一天扫描一次
		Config.Auth.CanRegister = true              // 允许注册
		Config.Auth.CanMailLogin = true             // 允许邮箱登录
		Config.PanelLog.Debug = false
		Config.PanelLog.SaveToFile = true
		Database.DB.Create(&Config)
	}
	Database.DB.First(&Config)
	PanelLog.IsDebug = Config.PanelLog.Debug
	PanelLog.IsSaveToFile = Config.PanelLog.SaveToFile
}

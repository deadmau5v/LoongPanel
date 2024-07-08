package notice

import (
	"LoongPanel/Panel/Service/Database"
	"fmt"

	"github.com/asaskevich/govalidator"
)

var (
	mailConfig = MailConfig{}
)

type MailConfig struct {
	MailHost  string `json:"mail_host" valid:"required~邮件服务器未配置,dns"`   // 邮件服务器
	MailPort  int    `json:"mail_port" valid:"required~邮件端口未配置,port"`   // 邮件端口
	MailUser  string `json:"mail_user" valid:"required~邮件用户名未配置"`       // 邮件用户名
	MailPass  string `json:"mail_pass" valid:"required~邮件密码未配置"`        // 邮件密码
	MailFrom  string `json:"mail_from" valid:"required~邮件发送者未配置,email"` // 邮件发送者
	MailSSL   bool   `json:"mail_ssl" valid:"required~邮件SSL未配置"`        // 邮件SSL
	MailTo    string `json:"mail_to" valid:"required~邮件接收者未配置,email"`   // 邮件接收者
	MailBody  string `json:"mail_body" valid:"required~邮件内容未配置"`        // 邮件内容
	MailTitle string `json:"mail_title" valid:"required~邮件标题未配置"`       // 邮件标题
}

// SendMail 发送邮件
func SendMail(to, title, body string) error {
	config := mailConfig
	config.MailTo = to
	config.MailTitle = title
	config.MailBody = body
	err := config.Check()
	if err != nil {
		return fmt.Errorf("SendMail -> %w", err)
	}
	return nil
}

// Check 检查邮件配置
func (c MailConfig) Check() error {
	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		return fmt.Errorf("Check -> %w", err)
	}
	return nil
}

// SendSMS 发送短信
// func SendSMS(to, title, body string) {
// 成本太贵，暂时不做
// }

func init() {
	Database.DB.AutoMigrate(&MailConfig{})
	Database.DB.Last(&mailConfig)
}

package notice

import (
	config "LoongPanel/Panel/Service/Config"
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"gopkg.in/gomail.v2"
)

// SendMail 发送邮件
func SendMail(to, title, body string) error {
	mailConfig := config.Config.MailConfig

	mailConfig.MailTo = to
	mailConfig.MailTitle = title
	mailConfig.MailBody = body
	res, err := govalidator.ValidateStruct(mailConfig)
	if err != nil {
		return fmt.Errorf("SendMail -> %w", err)
	}
	if !res {
		return fmt.Errorf("SendMail -> %w", errors.New("邮件配置错误"))
	}

	d := gomail.NewDialer(mailConfig.MailHost, mailConfig.MailPort, mailConfig.MailUser, mailConfig.MailPass)
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: !mailConfig.MailSSL,
		ServerName:         mailConfig.MailHost,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", mailConfig.MailFrom)
	m.SetHeader("To", mailConfig.MailTo)
	m.SetHeader("Subject", mailConfig.MailTitle)
	m.SetBody("text/plain", mailConfig.MailBody)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("SendMail -> %w", err)
	}

	return nil
}

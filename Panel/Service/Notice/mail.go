package notice

import (
	config "LoongPanel/Panel/Service/Config"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/asaskevich/govalidator"
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

	auth := smtp.PlainAuth("", mailConfig.MailUser, mailConfig.MailPass, mailConfig.MailHost)
	toList := []string{mailConfig.MailTo}
	msg := []byte("To: " + mailConfig.MailTo + "\r\n" +
		"Subject: " + mailConfig.MailTitle + "\r\n" +
		"\r\n" +
		mailConfig.MailBody + "\r\n")
	err = smtp.SendMail(fmt.Sprintf("%s:%d", mailConfig.MailHost, mailConfig.MailPort), auth, mailConfig.MailFrom, toList, msg)
	if err != nil {
		return fmt.Errorf("SendMail -> %w", err)
	}

	return nil
}

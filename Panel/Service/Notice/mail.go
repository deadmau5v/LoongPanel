package notice

import (
	config "LoongPanel/Panel/Service/Config"
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
)

var (
	mailConfig = config.Config.MailConfig
)

// SendMail 发送邮件
func SendMail(to, title, body string) error {
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
	return nil
}

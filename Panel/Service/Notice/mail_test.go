package notice

import (
	config "LoongPanel/Panel/Service/Config"
	"testing"
)

func TestSendMail(t *testing.T) {
	to := "2020268674@qq.com"
	title := "测试邮件"
	body := "这是一封测试邮件。"

	config.SetMailConfig(config.MailConfig{
		MailHost:  "smtp.qq.com",
		MailPort:  465,
		MailUser:  "deadmau5v",
		MailPass:  "uabhdvespuxbceab",
		MailFrom:  "deadmau5v@qq.com",
		MailSSL:   true,
		MailTo:    to,
		MailBody:  body,
		MailTitle: title,
	})

	err := SendMail(to, title, body)
	if err != nil {
		t.Errorf("发送邮件失败: %v", err)
	} else {
		t.Logf("邮件发送成功")
	}
}

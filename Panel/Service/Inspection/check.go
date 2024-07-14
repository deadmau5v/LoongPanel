package inspection

import (
	"LoongPanel/Panel/Service/Database"
	notice "LoongPanel/Panel/Service/Notice"
	"LoongPanel/Panel/Service/PanelLog"
	"encoding/json"
	"fmt"
)

type Item struct {
	Name       string
	check_func func() (string, string)
}

type Res struct {
	Status  string `json:"status"` // msg ERROR WARN OK
	Message string `json:"message"`
}

func (i *Item) Check() string {
	status, message := i.check_func()
	res, _ := json.Marshal(Res{Status: status, Message: message})
	return string(res)
}

func (i *Item) Start() string {
	res, _ := json.Marshal(Res{Status: "msg", Message: "正在检查" + i.Name})
	return string(res)
}

func GetItemsRes(res []Res) string {
	resTemplate := `
<h2>巡检结果<h2>
检查磁盘容量: %s
检查CPU状态: %s
检查内存状态: %s
检查网络访问: %s
检查系统运行时间: %s
检查时间同步: %s
检查SSH服务状态: %s
检查SSH : %s
检查SSH密码登录: %s
检查SSH端口: %s
检查SSH协议版本: %s
检查SSH最大认证尝试次数: %s
检查SSH登录宽限时间: %s
检查防火墙状态: %s
检查防火墙规则: %s
`
	resStr := []interface{}{}
	for _, v := range res {
		switch v.Status {
		case "OK":
			resStr = append(resStr, "正常 "+v.Message)
		case "WARN":
			resStr = append(resStr, "警告 "+v.Message)
		case "ERROR":
			resStr = append(resStr, "异常 "+v.Message)
		}
	}
	return fmt.Sprintf(resTemplate, resStr...)
}

// GetAllItems 获取所有检查项
func GetAllItems() []Item {
	return []Item{
		{Name: "检查磁盘容量", check_func: CheckDiskUsage},
		{Name: "检查CPU状态", check_func: CheckCPUStatus},
		{Name: "检查内存状态", check_func: CheckMemoryStatus},
		{Name: "检查网络访问", check_func: CheckInternetAccess},
		{Name: "检查系统运行时间", check_func: CheckUptime},
		{Name: "检查时间同步", check_func: CheckTimeSync},
		{Name: "检查SSH服务状态", check_func: CheckSSHService},
		{Name: "检查SSH Root登录", check_func: CheckSSHRootLogin},
		{Name: "检查SSH密码登录", check_func: CheckSSHPasswordLogin},
		{Name: "检查SSH端口", check_func: CheckSSHPort},
		{Name: "检查SSH协议版本", check_func: CheckSSHProtocol},
		{Name: "检查SSH最大认证尝试次数", check_func: CheckSSHMaxAuthTries},
		{Name: "检查SSH登录宽限时间", check_func: CheckSSHLoginGraceTime},
		{Name: "检查防火墙状态", check_func: CheckFirewallStatus},
		{Name: "检查防火墙规则", check_func: CheckFirewallRules},
	}
}

// 执行巡检
func Check() chan string {
	c := make(chan string, 100)

	go func() {
		results := []Res{}

		for _, item := range GetAllItems() {
			c <- item.Start()

			checkResult := item.Check()
			var checkResultJson Res
			json.Unmarshal([]byte(checkResult), &checkResultJson)
			results = append(results, checkResultJson)

			c <- checkResult
		}
		close(c)
		mailBody := GetItemsRes(results)
		var settings []notice.UserNotificationSetting
		Database.DB.Preload("User").Find(&settings)
		for _, v := range settings {
			if v.InspectionNotify {
				err := notice.SendMail(v.User.Mail, "巡检结果", mailBody)
				if err != nil {
					PanelLog.ERROR("[巡检通知]", err.Error())
				}
			}
		}
	}()

	return c
}

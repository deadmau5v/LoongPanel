// 巡检项

package inspection

import (
	"LoongPanel/Panel/Service/PanelLog"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// CheckDiskUsage 检查服务器磁盘的容量占比
func CheckDiskUsage() (string, string) {
	// 获取磁盘使用情况
	usageStat, err := disk.Usage("/")
	if err != nil {
		PanelLog.ERROR("获取磁盘使用情况失败: %v", err)
	}

	// 如果大于 85% 报警
	if usageStat.UsedPercent > 95 {
		return "ERROR", "磁盘使用率大于 95%"
	} else if usageStat.UsedPercent > 85 {
		return "WARN", "磁盘使用率大于 85%"
	} else if usageStat.UsedPercent > 70 {
		return "WARN", "磁盘使用率大于 70%"
	} else {
		return "OK", "磁盘剩余空间正常"
	}
}

// CheckCPUStatus 检查每个CPU的状态和温度
func CheckCPUStatus() (string, string) {
	// 获取CPU使用情况
	cpuStats, err := cpu.Percent(0, true)
	if err != nil {
		PanelLog.ERROR("获取CPU使用情况失败: %v", err)
	}

	// 获取CPU温度
	temperatureStats, err := host.SensorsTemperatures()
	if err != nil {
		PanelLog.ERROR("获取CPU温度失败: %v", err)
	}

	for _, percent := range cpuStats {
		if percent > 90 {
			return "ERROR", "CPU使用率大于 90%"
		} else if percent > 75 {
			return "WARN", "CPU使用率大于 75%"
		}
	}

	for _, temp := range temperatureStats {
		if temp.SensorKey == "coretemp" && temp.Temperature > 80 {
			return "ERROR", "CPU温度大于 80°C"
		} else if temp.SensorKey == "coretemp" && temp.Temperature > 70 {
			return "WARN", "CPU温度大于 70°C"
		}
	}

	return "OK", "CPU状态正常"
}

// CheckMemoryStatus 检查内存状态
func CheckMemoryStatus() (string, string) {
	// 获取内存使用情况
	virtualMemoryStat, err := mem.VirtualMemory()
	if err != nil {
		PanelLog.ERROR("获取内存使用情况失败: %v", err)
	}

	// 获取交换分区使用情况
	swapMemoryStat, err := mem.SwapMemory()
	if err != nil {
		PanelLog.ERROR("获取交换分区使用情况失败: %v", err)
	}

	if virtualMemoryStat.UsedPercent >= 90 {
		return "ERROR", "物理机内存使用率不低于90% 请清理内存"
	} else if swapMemoryStat.UsedPercent >= 10 {
		return "WARN", "物理机交换分区使用率不低于10% 请清理内存"
	}

	return "OK", "内存状态正常"
}

// CheckInternetAccess 检查是否能访问外网
func CheckInternetAccess() (string, string) {
	// 检查是否能访问外网
	_, err := http.Get("http://www.baidu.com")
	if err != nil {
		return "ERROR", "无法访问外网 请检查网络"
	}
	return "OK", "可以访问外网"
}

// CheckUptime 检查系统运行时间
func CheckUptime() (string, string) {
	// 获取系统运行时间
	uptimeStat, err := host.Uptime()
	if err != nil {
		PanelLog.ERROR("获取系统运行时间失败: %v", err)
		return "ERROR", "无法获取系统运行时间"
	}

	// 将运行时间转换为小时
	uptimeHours := uptimeStat / 3600

	if uptimeHours > 720 {
		return "WARN", "系统运行时间超过720小时 请及时重启"
	}

	return "OK", "系统运行时间正常"
}

// CheckTimeSync 检查时间是否同步
func CheckTimeSync() (string, string) {
	// 获取本地时间
	localTime := time.Now()

	// 获取网络时间
	resp, err := http.Get("http://worldtimeapi.org/api/timezone/Etc/UTC")
	if err != nil {
		PanelLog.ERROR("获取网络时间失败: %v", err)
		return "ERROR", "无法获取网络时间 请检查网络"
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		PanelLog.ERROR("解析网络时间失败: %v", err)
		return "ERROR", "无法解析网络时间"
	}

	utcTimeStr, ok := result["utc_datetime"].(string)
	if !ok {
		PanelLog.ERROR("获取UTC时间失败")
		return "ERROR", "无法获取UTC时间"
	}

	utcTime, err := time.Parse(time.RFC3339, utcTimeStr)
	if err != nil {
		PanelLog.ERROR("解析UTC时间失败: %v", err)
		return "ERROR", "无法解析UTC时间"
	}

	// 计算时间差
	timeDiff := localTime.Sub(utcTime)
	if timeDiff > time.Minute || timeDiff < -time.Minute {
		return "WARN", "系统时间不同步 请检查时间设置"
	}

	return "OK", "系统时间同步"
}

// CheckSSHService 检查SSH服务是否在运行
func CheckSSHService() (string, string) {
	// 执行命令检查SSH服务状态
	cmd := exec.Command("systemctl", "is-active", "ssh")
	output, err := cmd.Output()
	if err != nil {
		PanelLog.ERROR("检查SSH服务状态失败: %v", err)
		return "ERROR", "无法检查SSH服务状态"
	}

	// 判断SSH服务是否在运行
	status := strings.TrimSpace(string(output))
	if status == "active" {
		return "OK", "SSH服务正在运行"
	}

	return "WARN", "SSH服务未运行 请检查服务状态"
}

// CheckSSHRootLogin 检查SSH设置 是否允许Root登录
func CheckSSHRootLogin() (string, string) {
	// 读取SSH配置文件
	config, err := os.ReadFile("/etc/ssh/sshd_config")
	if err != nil {
		PanelLog.ERROR("读取SSH配置文件失败: %v", err)
		return "ERROR", "无法读取SSH配置文件"
	}

	// 检查是否允许Root登录
	if strings.Contains(string(config), "PermitRootLogin yes") {
		return "WARN", "SSH允许Root登录 请修改配置以提高安全性"
	}

	return "OK", "SSH不允许Root登录"
}

// CheckSSHPasswordLogin 检查SSH是否允许密码登录
func CheckSSHPasswordLogin() (string, string) {
	// 读取SSH配置文件
	config, err := os.ReadFile("/etc/ssh/sshd_config")
	if err != nil {
		PanelLog.ERROR("读取SSH配置文件失败: %v", err)
		return "ERROR", "无法读取SSH配置文件"
	}

	// 检查是否允许密码登录
	if strings.Contains(string(config), "PasswordAuthentication yes") {
		return "WARN", "SSH允许密码登录 请修改配置以提高安全性"
	}

	return "OK", "SSH不允许密码登录"
}

// CheckSSHPort 检查SSH服务端口
func CheckSSHPort() (string, string) {
	// 读取SSH配置文件
	config, err := os.ReadFile("/etc/ssh/sshd_config")
	if err != nil {
		PanelLog.ERROR("读取SSH配置文件失败: %v", err)
		return "ERROR", "无法读取SSH配置文件"
	}

	// 检查SSH服务端口
	var port string
	lines := strings.Split(string(config), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Port ") {
			port = strings.TrimSpace(strings.TrimPrefix(line, "Port "))
			break
		}
	}

	if port == "22" {
		return "WARN", "SSH服务端口为默认端口 请修改配置以提高安全性"
	}

	return "OK", "SSH端口安全"
}

// CheckSSHProtocol 检查SSH协议版本
func CheckSSHProtocol() (string, string) {
	// 读取SSH配置文件
	config, err := os.ReadFile("/etc/ssh/sshd_config")
	if err != nil {
		PanelLog.ERROR("读取SSH配置文件失败: %v", err)
		return "ERROR", "无法读取SSH配置文件"
	}

	// 检查SSH协议版本
	if strings.Contains(string(config), "Protocol 2") {
		return "OK", "SSH使用协议版本2"
	}

	return "WARN", "SSH未使用协议版本2 请修改配置以提高安全性"
}

// CheckSSHMaxAuthTries 检查SSH最大认证尝试次数
func CheckSSHMaxAuthTries() (string, string) {
	// 读取SSH配置文件
	config, err := os.ReadFile("/etc/ssh/sshd_config")
	if err != nil {
		PanelLog.ERROR("读取SSH配置文件失败: %v", err)
		return "ERROR", "无法读取SSH配置文件"
	}

	// 检查最大认证尝试次数
	maxAuthTries := "6" // 默认值
	lines := strings.Split(string(config), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MaxAuthTries ") {
			maxAuthTries = strings.TrimSpace(strings.TrimPrefix(line, "MaxAuthTries "))
			break
		}
	}

	return "OK", "SSH最大认证尝试次数为 " + maxAuthTries
}

// CheckSSHLoginGraceTime 检查SSH登录宽限时间
func CheckSSHLoginGraceTime() (string, string) {
	// 读取SSH配置文件
	config, err := os.ReadFile("/etc/ssh/sshd_config")
	if err != nil {
		PanelLog.ERROR("读取SSH配置文件失败: %v", err)
		return "ERROR", "无法读取SSH配置文件"
	}

	// 检查登录宽限时间
	loginGraceTime := "120" // 默认值
	lines := strings.Split(string(config), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "LoginGraceTime ") {
			loginGraceTime = strings.TrimSpace(strings.TrimPrefix(line, "LoginGraceTime "))
			break
		}
	}

	return "OK", "SSH登录宽限时间为 " + loginGraceTime + " 秒"
}

// CheckFirewallStatus 检查防火墙状态
func CheckFirewallStatus() (string, string) {
	// 检查防火墙状态
	out, err := exec.Command("ufw", "status").Output()
	if err != nil {
		PanelLog.ERROR("获取防火墙状态失败: %v", err)
		return "ERROR", "无法获取防火墙状态"
	}

	status := string(out)
	if strings.Contains(status, "inactive") {
		return "WARN", "防火墙未激活"
	}

	return "OK", "防火墙状态正常"
}

// CheckFirewallRules 检查防火墙规则
func CheckFirewallRules() (string, string) {
	// 获取防火墙规则
	out, err := exec.Command("ufw", "status", "numbered").Output()
	if err != nil {
		PanelLog.ERROR("获取防火墙规则失败: %v", err)
		return "ERROR", "无法获取防火墙规则"
	}

	rules := string(out)
	if rules == "" {
		return "WARN", "防火墙未设置任何规则"
	}

	return "OK", "防火墙规则正常"
}

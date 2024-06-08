/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-29
 * 文件作用：获取日志内容
 */

package SystemLog

import (
	"LoongPanel/Panel/Service/Log"
	Log2 "LoongPanel/Panel/Service/PanelLog"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type LogEntry interface {
	ProcessLogLine(logLine string) (any, error)
}

// GetLog 获取日志
func GetLog(log *Log.Log_, line int, entry LogEntry) interface{} {
	if !log.Ok {
		return nil
	}

	Log2.DEBUG("获取日志", log.Name, log.Path)
	file, err := os.Open(log.Path)
	if err != nil {
		log.Ok = false
		Log2.ERROR("打开日志文件失败", log.Name, log.Path)
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Ok = false
			Log2.ERROR("关闭日志文件失败", log.Name, log.Path)
		}
	}(file)
	all, err := io.ReadAll(file)
	if err != nil {
		log.Ok = false
		Log2.ERROR("读取日志文件失败", log.Name, log.Path)
		return nil
	}

	if line != 0 {
		allStr := string(all)
		allStrSplit := strings.Split(allStr, "\n")
		if len(allStrSplit) > line {
			allStrSplit = allStrSplit[len(allStrSplit)-line:]
			all = []byte(strings.Join(allStrSplit, "\n"))
		}
	}

	sp := strings.Split(string(all), "\n")
	var res []interface{}
	for _, line := range sp {
		logEntry, err := entry.ProcessLogLine(line)
		if err != nil {
			continue
		}
		if logEntry != nil {
			res = append(res, logEntry)
		}
	}
	return res
}

// ClearLog 清空日志
func ClearLog(log *Log.Log_) {
	if !log.Ok {
		return
	}

	file, err := os.Open(log.Path)
	if err != nil {
		log.Ok = false
		Log2.ERROR("[系统日志] 打开日志文件错误", err.Error())
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Ok = false
			Log2.ERROR("[系统日志] 关闭日志文件错误", err.Error())
		}
	}(file)
	// 截断文件
	err = file.Truncate(0)
	if err != nil {
		log.Ok = false
		Log2.ERROR("[系统日志] 清空日志文件错误", err.Error())
		return
	}
}

// creatLog 创建日志对象 简化流程
func createLog(path, name string, customGetLog func(log *Log.Log_, line int) interface{}, entry LogEntry) *Log.Log_ {
	log := &Log.Log_{
		Path: path,
		Name: name,
	}
	log.Ok = log.CheckLogExist()

	if !log.Ok {
		return nil
	}

	log.GetLog = func(line int) interface{} {
		if customGetLog != nil {
			return customGetLog(log, line)
		}
		return GetLog(log, line, entry)
	}
	log.ClearLog = func() {
		ClearLog(log)
	}

	return log
}

type BootLog struct {
	Level   string `json:"level"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	Module  string `json:"module"`
	Content string `json:"content"`
}

func (b BootLog) ProcessLogLine(logLine string) (any, error) {
	pattern := `^\[(.*?)\] (.*?) - (.*?) \[(.*?)\] (.*)$`

	checkLogLine := strings.Replace(logLine, " ", "", -1)
	checkLogLine = strings.Replace(checkLogLine, "\t", "", -1)
	checkLogLine = strings.Replace(checkLogLine, "\n", "", -1)
	checkLogLine = strings.Replace(checkLogLine, "\r", "", -1)
	if strings.Trim(checkLogLine, " ") == "" {
		return nil, errors.New("空行")
	}
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(logLine)
	if matches == nil {
		fmt.Println("[日志管理] 无法解析日志行", logLine)
		return nil, errors.New("无法解析日志行")
	}

	entry := BootLog{
		Level:   matches[1],
		Date:    matches[2],
		Time:    matches[3],
		Module:  matches[4],
		Content: matches[5],
	}

	return &entry, nil
}

// GetBootLog 获取启动日志
func GetBootLog() *Log.Log_ {
	log := createLog("/var/log/boot.log", "系统启动日志", nil, BootLog{})
	log.Struct = append(log.Struct, []map[string]string{
		{
			"title":     "日志等级",
			"dataIndex": "level",
			"key":       "1",
		},
		{
			"title":     "日期",
			"dataIndex": "date",
			"key":       "2",
		},
		{
			"title":     "时间",
			"dataIndex": "time",
			"key":       "3",
		},
		{
			"title":     "模块",
			"dataIndex": "module",
			"key":       "4",
		},
		{
			"title":     "内容",
			"dataIndex": "content",
			"key":       "5",
		},
	})
	return log
}

type CronLog struct {
	Date    string `json:"date"`
	Time    string `json:"time"`
	Host    string `json:"host"`
	Module  string `json:"module"`
	Content string `json:"content"`
}

func (c CronLog) ProcessLogLine(logLine string) (any, error) {
	pattern := `^(\w+\s+\d+\s+\d+:\d+:\d+)\s+(\S+)\s+(\S+)\[(\d+)\]:\s+(.*)$`

	checkLogLine := strings.ReplaceAll(logLine, " ", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\t", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\n", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\r", "")
	if strings.Trim(checkLogLine, " ") == "" {
		return nil, errors.New("空行")
	}

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(logLine)
	if matches == nil {
		fmt.Println("[日志管理] 无法解析日志行:", logLine)
		return nil, errors.New("无法解析日志行")
	}

	entry := &CronLog{
		Date:    matches[1],
		Host:    matches[2],
		Module:  matches[3],
		Content: matches[5],
	}

	return entry, nil
}

// GetCronLog 获取定时任务日志
func GetCronLog() *Log.Log_ {
	log := createLog("/var/log/cron", "计划任务日志", nil, CronLog{})
	log.Struct = append(log.Struct, []map[string]string{
		{
			"title":     "日志等级",
			"dataIndex": "level",
			"key":       "1",
		},
		{
			"title":     "日期",
			"dataIndex": "date",
			"key":       "2",
		},
		{
			"title":     "时间",
			"dataIndex": "time",
			"key":       "3",
		},
		{
			"title":     "模块",
			"dataIndex": "module",
			"key":       "4",
		},
		{
			"title":     "内容",
			"dataIndex": "content",
			"key":       "5",
		},
	})
	return log
}

type FirewalldLog struct {
	Date    string `json:"date"`
	Time    string `json:"time"`
	Level   string `json:"level"`
	Module  string `json:"module"`
	Content string `json:"content"`
}

func (c FirewalldLog) ProcessLogLine(logLine string) (any, error) {
	pattern := `^(\d{4}-\d{2}-\d{2})\s+(\d{2}:\d{2}:\d{2})\s+(WARNING|INFO|ERROR):\s+(.*)$`
	checkLogLine := strings.ReplaceAll(logLine, " ", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\t", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\n", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\r", "")
	if strings.Trim(checkLogLine, " ") == "" {
		return nil, errors.New("空行")
	}

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(logLine)
	if matches == nil {
		fmt.Println("[日志管理] 无法解析日志行:", logLine)
		return nil, errors.New("无法解析日志行")
	}

	entry := &FirewalldLog{
		Date:    matches[1],
		Time:    matches[2],
		Level:   matches[3],
		Module:  "systemd",
		Content: matches[5],
	}

	return entry, nil
}

// GetFirewalldLog 获取Firewalld日志
func GetFirewalldLog() *Log.Log_ {
	log := createLog("/var/log/firewalld", "防火墙日志", nil, FirewalldLog{})
	log.Struct = append(log.Struct, []map[string]string{
		{
			"title":     "日志等级",
			"dataIndex": "level",
			"key":       "1",
		},
		{
			"title":     "日期",
			"dataIndex": "date",
			"key":       "2",
		},
		{
			"title":     "时间",
			"dataIndex": "time",
			"key":       "3",
		},
		{
			"title":     "模块",
			"dataIndex": "module",
			"key":       "4",
		},
		{
			"title":     "内容",
			"dataIndex": "content",
			"key":       "5",
		},
	})
	return log
}

type MessagesLog struct {
	Date    string `json:"date"`
	Time    string `json:"time"`
	Host    string `json:"host"`
	Module  string `json:"module"`
	Content string `json:"content"`
}

func (c MessagesLog) ProcessLogLine(logLine string) (any, error) {
	pattern := `^(\w+\s+\d+)\s+(\d{2}:\d{2}:\d{2})\s+(\S+)\s+(\S+)\[(\d+)\]:\s+(.*)$`
	checkLogLine := strings.ReplaceAll(logLine, " ", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\t", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\n", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\r", "")
	if strings.Trim(checkLogLine, " ") == "" {
		return nil, errors.New("空行")
	}

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(logLine)
	if matches == nil {
		fmt.Println("[日志管理] 无法解析日志行:", logLine)
		return nil, errors.New("无法解析日志行")
	}

	entry := &MessagesLog{
		Date:    matches[1],
		Time:    matches[2],
		Host:    matches[3],
		Module:  matches[4],
		Content: matches[5],
	}

	return entry, nil
}

// GetMessagesLog 获取系统消息日志
func GetMessagesLog() *Log.Log_ {
	log := createLog("/var/log/messages", "系统消息日志", nil, MessagesLog{})
	log.Struct = append(log.Struct, []map[string]string{
		{
			"title":     "日期",
			"dataIndex": "date",
			"key":       "1",
		},
		{
			"title":     "时间",
			"dataIndex": "time",
			"key":       "2",
		},
		{
			"title":     "主机名",
			"dataIndex": "host",
			"key":       "3",
		},
		{
			"title":     "模块",
			"dataIndex": "module",
			"key":       "4",
		},
		{
			"title":     "内容",
			"dataIndex": "content",
			"key":       "5",
		},
	})
	return log
}

type SecureLog struct {
	Date    string `json:"date"`
	Time    string `json:"time"`
	Host    string `json:"host"`
	Module  string `json:"module"`
	Content string `json:"content"`
}

func (c SecureLog) ProcessLogLine(logLine string) (any, error) {
	pattern := `^(\w+\s+\d+)\s+(\d{2}:\d{2}:\d{2})\s+(\S+)\s+(\S+)\[(\d+)\]:\s+(.*)$`
	checkLogLine := strings.ReplaceAll(logLine, " ", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\t", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\n", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\r", "")
	if strings.Trim(checkLogLine, " ") == "" {
		return nil, errors.New("空行")
	}

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(logLine)
	if matches == nil {
		fmt.Println("[日志管理] 无法解析日志行:", logLine)
		return nil, errors.New("无法解析日志行")
	}

	entry := &SecureLog{
		Date:    matches[1],
		Time:    matches[2],
		Host:    matches[3],
		Module:  matches[4],
		Content: matches[5],
	}

	return entry, nil
}

// GetSecureLog 获取安全日志
func GetSecureLog() *Log.Log_ {
	log := createLog("/var/log/secure", "安全日志", nil, SecureLog{})
	log.Struct = append(log.Struct, []map[string]string{
		{
			"title":     "日期",
			"dataIndex": "date",
			"key":       "1",
		},
		{
			"title":     "时间",
			"dataIndex": "time",
			"key":       "2",
		},
		{
			"title":     "主机名",
			"dataIndex": "host",
			"key":       "3",
		},
		{
			"title":     "模块",
			"dataIndex": "module",
			"key":       "4",
		},
		{
			"title":     "内容",
			"dataIndex": "content",
			"key":       "5",
		},
	})
	return log
}

type WtmpLog struct {
	Level     string `json:"level"`
	PID       string `json:"pid"`
	Type      string `json:"type"`
	User      string `json:"user"`
	Terminal  string `json:"terminal"`
	SrcIP     string `json:"src_ip"`
	DestIP    string `json:"dest_ip"`
	Timestamp string `json:"timestamp"`
}

func (c WtmpLog) ProcessLogLine(logLine string) (any, error) {
	pattern := `^\[(\d+)\] \[(\d+)\] \[(.*?)\] \[(.*?)\] \[(.*?)\] \[(.*?)\] \[(.*?)\] \[(.*?)\]$`
	checkLogLine := strings.ReplaceAll(logLine, " ", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\t", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\n", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\r", "")
	if strings.Trim(checkLogLine, " ") == "" {
		return nil, errors.New("空行")
	}

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(logLine)
	if matches == nil {
		fmt.Println("[日志管理] 无法解析日志行:", logLine)
		return nil, errors.New("无法解析日志行")
	}

	entry := &SecureLog{
		Date:    matches[1],
		Time:    matches[2],
		Host:    matches[3],
		Module:  matches[4],
		Content: matches[5],
	}

	return entry, nil
}

// GetWtmpLog 获取登录日志
func GetWtmpLog() *Log.Log_ {
	log := createLog("/var/log/wtmp", "登录日志", func(log *Log.Log_, line int) interface{} {
		output, err := exec.Command("utmpdump", log.Path).Output()
		if err != nil {
			log.Ok = false
			Log2.ERROR("执行 utmpdump 命令失败", log.Name, log.Path)
			return nil
		}

		// 将输出转换为字符串并按行分割
		outputStr := string(output)
		outputStrSplit := strings.Split(outputStr, "\n")

		if len(outputStrSplit) > line {
			outputStrSplit = outputStrSplit[len(outputStrSplit)-line:]
		}

		var res []interface{}
		for _, line := range outputStrSplit {
			logEntry, err := WtmpLog{}.ProcessLogLine(line)
			if err != nil {
				continue
			}
			if logEntry != nil {
				res = append(res, logEntry)
			}

		}
		return res
	}, WtmpLog{})

	log.Struct = append(log.Struct, []map[string]string{
		{
			"title":     "级别",
			"dataIndex": "level",
			"key":       "1",
		},
		{
			"title":     "进程ID",
			"dataIndex": "pid",
			"key":       "2",
		},
		{
			"title":     "类型",
			"dataIndex": "type",
			"key":       "3",
		},
		{
			"title":     "用户",
			"dataIndex": "user",
			"key":       "4",
		},
		{
			"title":     "终端",
			"dataIndex": "terminal",
			"key":       "5",
		},
		{
			"title":     "源IP",
			"dataIndex": "src_ip",
			"key":       "6",
		},
		{
			"title":     "目的IP",
			"dataIndex": "dest_ip",
			"key":       "7",
		},
		{
			"title":     "时间戳",
			"dataIndex": "timestamp",
			"key":       "8",
		},
	})
	return log
}

type KernelLog struct {
	Date    string `json:"date"`
	Time    string `json:"time"`
	Host    string `json:"host"`
	Module  string `json:"module"`
	Content string `json:"content"`
}

func (c KernelLog) ProcessLogLine(logLine string) (any, error) {
	pattern := `^(\d{1,2}月\s+\d{2})\s+(\d{2}:\d{2}:\d{2})\s+(\S+)\s+(\S+):\s+(.*)$`
	checkLogLine := strings.ReplaceAll(logLine, " ", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\t", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\n", "")
	checkLogLine = strings.ReplaceAll(checkLogLine, "\r", "")
	if strings.Trim(checkLogLine, " ") == "" {
		return nil, errors.New("空行")
	}

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(logLine)
	if matches == nil {
		fmt.Println("[日志管理] 无法解析日志行:", logLine)
		return nil, errors.New("无法解析日志行")
	}

	entry := &KernelLog{
		Date:    matches[1],
		Time:    matches[2],
		Host:    matches[3],
		Module:  matches[4],
		Content: matches[5],
	}

	return entry, nil
}

// GetKernelLog 获取内核日志
func GetKernelLog() *Log.Log_ {
	log := Log.Log_{
		Name: "系统日志",
		Ok:   true,
	}

	log.Struct = append(log.Struct, []map[string]string{
		{
			"title":     "日期",
			"dataIndex": "date",
			"key":       "1",
		},
		{
			"title":     "时间",
			"dataIndex": "time",
			"key":       "2",
		},
		{
			"title":     "主机名",
			"dataIndex": "host",
			"key":       "3",
		},
		{
			"title":     "模块",
			"dataIndex": "module",
			"key":       "4",
		},
		{
			"title":     "内容",
			"dataIndex": "content",
			"key":       "5",
		},
	})
	log.GetLog = func(line int) interface{} {
		output, err := exec.Command("journalctl", "-k").Output()
		if err != nil {
			log.Ok = false
			Log2.ERROR("执行 journalctl 命令失败", log.Name)
			return nil
		}

		// 将输出转换为字符串并按行分割
		outputStr := string(output)
		outputStrSplit := strings.Split(outputStr, "\n")

		if len(outputStrSplit) > line {
			outputStrSplit = outputStrSplit[len(outputStrSplit)-line:]
		}
		var res []interface{}
		for _, line := range outputStrSplit {
			logEntry, err := KernelLog{}.ProcessLogLine(line)
			if err != nil {
				continue
			}
			if logEntry != nil {
				res = append(res, logEntry)
			}

		}
		return res
	}

	log.ClearLog = func() {
		err := exec.Command("journalctl", "--rotate").Run()
		if err != nil {
			log.Ok = false
			Log2.ERROR("执行 journalctl --rotate 命令失败", log.Name)
		}
		return
	}
	return &log
}

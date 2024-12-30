/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-29
 * 文件作用：获取日志内容
 */

package SystemLog

import (
	"LoongPanel/Panel/Service/Log"
	"LoongPanel/Panel/Service/PanelLog"
	"errors"
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

	PanelLog.DEBUG("[日志管理] 获取日志", log.Name, log.Path)
	file, err := os.Open(log.Path)
	if err != nil {
		log.Ok = false
		PanelLog.ERROR("打开日志文件失败", log.Name, log.Path)
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Ok = false
			PanelLog.ERROR("关闭日志文件失败", log.Name, log.Path)
		}
	}(file)
	all, err := io.ReadAll(file)
	if err != nil {
		log.Ok = false
		PanelLog.ERROR("读取日志文件失败", log.Name, log.Path)
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
	err := exec.Command("echo", ">", log.Path).Run()
	if err != nil {
		PanelLog.ERROR("[日志管理] 清空日志失败", log.Name, log.Path)
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
	Content string `json:"content"`
}

func (b BootLog) ProcessLogLine(logLine string) (any, error) {
	checkLogLine := strings.Replace(logLine, " ", "", -1)
	checkLogLine = strings.Replace(checkLogLine, "\t", "", -1)
	checkLogLine = strings.Replace(checkLogLine, "\n", "", -1)
	checkLogLine = strings.Replace(checkLogLine, "\r", "", -1)
	if strings.Trim(checkLogLine, " ") == "" {
		return nil, errors.New("空行")
	}

	entry := BootLog{
		Content: logLine,
	}

	return &entry, nil
}

// GetBootLog 获取启动日志
func GetBootLog() *Log.Log_ {
	log := createLog("/var/log/boot.log", "系统启动日志", nil, BootLog{})
	if log == nil {
		return nil
	}
	log.Struct = append(log.Struct, []map[string]string{
		{
			"title":     "内容",
			"dataIndex": "content",
			"key":       "1",
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
	pattern := `(?P<date>\w+\s+\d+)\s+(?P<time>\d{2}:\d{2}:\d{2})\s+(?P<host>\S+)\s+(?P<module>[\w\[\]\d]+):\s+(?P<content>.*)`

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
		return nil, errors.New("无法解析日志行")
	}

	md := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			md[name] = matches[i]
		}
	}

	entry := &CronLog{
		Date:    md["date"],
		Time:    md["time"],
		Host:    md["host"],
		Module:  md["module"],
		Content: md["content"],
	}

	return entry, nil
}

// GetCronLog 获取定时任务日志
func GetCronLog() *Log.Log_ {
	log := createLog("/var/log/cron", "计划任务日志", nil, CronLog{})
	if log == nil {
		return nil
	}
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
		return nil, errors.New("无法解析日志行")
	}

	entry := &FirewalldLog{
		Date:    matches[1],
		Time:    matches[2],
		Level:   matches[3],
		Module:  "systemd",
		Content: matches[4],
	}

	return entry, nil
}

// GetFirewalldLog 获取Firewalld日志
func GetFirewalldLog() *Log.Log_ {
	log := createLog("/var/log/firewalld", "防火墙日志", nil, FirewalldLog{})
	if log == nil {
		return nil
	}
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
	pattern := `(?P<date>\w+\s+\d+)\s+(?P<time>\d{2}:\d{2}:\d{2})\s+(?P<host>[^\s]+)\s+(?P<module>[^\s]+):\s+(?P<content>.*)`
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
		return nil, errors.New("无法解析日志行")
	}

	entry := MessagesLog{
		Date:    matches[re.SubexpIndex("date")],
		Time:    matches[re.SubexpIndex("time")],
		Host:    matches[re.SubexpIndex("host")],
		Module:  matches[re.SubexpIndex("module")],
		Content: matches[re.SubexpIndex("content")],
	}

	return entry, nil
}

// GetMessagesLog 获取系统消息日志
func GetMessagesLog() *Log.Log_ {
	log := createLog("/var/log/messages", "系统消息日志", nil, MessagesLog{})
	if log == nil {
		return nil
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
	pattern := `(?P<date>\w+\s+\d+)\s+(?P<time>\d{2}:\d{2}:\d{2})\s+(?P<host>[^\s]+)\s+(?P<module>[^\s]+):\s+(?P<content>.*)`
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
		return nil, errors.New("无法解析日志行")
	}

	entry := &SecureLog{
		Date:    matches[re.SubexpIndex("date")],
		Time:    matches[re.SubexpIndex("time")],
		Host:    matches[re.SubexpIndex("host")],
		Module:  matches[re.SubexpIndex("module")],
		Content: matches[re.SubexpIndex("content")],
	}

	return entry, nil
}

// GetSecureLog 获取安全日志
func GetSecureLog() *Log.Log_ {
	log := createLog("/var/log/secure", "安全日志", nil, SecureLog{})
	if log == nil {
		return nil
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
		return nil, errors.New("无法解析日志行")
	}

	entry := &WtmpLog{
		Level:     matches[1],
		PID:       matches[2],
		Type:      matches[3],
		User:      matches[4],
		Terminal:  matches[5],
		SrcIP:     matches[6],
		DestIP:    matches[7],
		Timestamp: matches[8],
	}

	return entry, nil
}

// GetWtmpLog 获取登录日志
func GetWtmpLog() *Log.Log_ {
	log := createLog("/var/log/wtmp", "登录日志", func(log *Log.Log_, line int) interface{} {
		output, err := exec.Command("utmpdump", log.Path).Output()
		if err != nil {
			log.Ok = false
			PanelLog.ERROR("执行 utmpdump 命令失败", log.Name, log.Path)
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

	if log == nil {
		return nil
	}

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
		output, err := exec.Command("journalctl", "-k", "--no-pager").Output()
		if err != nil {
			log.Ok = false
			PanelLog.ERROR("执行 journalctl 命令失败", log.Name)
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
		// 删除 /run/log/journal/
		err := exec.Command("rm", "-rf", "/run/log/journal/").Run()
		if err != nil {
			log.Ok = false
			PanelLog.ERROR("删除 /run/log/journal/ 失败", log.Name)
		}
		return
	}
	return &log
}

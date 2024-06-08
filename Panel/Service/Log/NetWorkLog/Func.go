/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-30
 * 文件作用：网络日志
 */

package NetWorkLog

import (
	"LoongPanel/Panel/Service/Log"
	Log2 "LoongPanel/Panel/Service/PanelLog"
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type LogEntry struct {
	Level     string `json:"level"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Hostname  string `json:"hostname"`
	Process   string `json:"process"`
	PID       int    `json:"pid"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func GetNetWorkLog() *Log.Log_ {
	log := Log.Log_{
		Ok:   true,
		Name: "网络日志",
	}

	log.GetLog = func(line int) interface{} {
		output, err := exec.Command("journalctl", "-u", "NetworkManager", "-n", strconv.Itoa(line), "--no-pager").Output()
		if err != nil {
			log.Ok = false
			Log2.ERROR("[日志管理] 获取网络日志失败：" + err.Error())
			return nil
		}

		logs := strings.Split(string(output), "\n")
		var logEntries []LogEntry
		for _, logLine := range logs {
			entry, err := parseLog(logLine)
			if err != nil {
				continue
			}
			logEntries = append(logEntries, *entry)
		}

		return logEntries
	}

	log.ClearLog = func() {
		_, err := exec.Command("journalctl", "-u", "NetworkManager", "--rotate").Output()
		// Todo 未测试
		if err != nil {
			log.Ok = false
			Log2.ERROR("[日志管理] 清空网络日志失败：" + err.Error())
			return
		}
		return
	}

	// 测试是否可用
	log.GetLog(1)

	if !log.Ok {
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
			"title":     "主机",
			"dataIndex": "hostname",
			"key":       "4",
		},
		{
			"title":     "进程",
			"dataIndex": "process",
			"key":       "5",
		},
		{
			"title":     "PID",
			"dataIndex": "pid",
			"key":       "6",
		},
		{
			"title":     "时间戳",
			"dataIndex": "timestamp",
			"key":       "7",
		},
		{
			"title":     "内容",
			"dataIndex": "message",
			"key":       "8",
		},
	})

	return &log
}

func init() {
}

// parseLog 解析日志
func parseLog(logLine string) (*LogEntry, error) {
	pattern := `(\w+ \d+) (\d{2}:\d{2}:\d{2}) (\w+) (\w+)\[(\d+)\]: <(\w+)> \[(\d+\.\d+)\] (.+)`
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
		Log2.DEBUG("[日志管理] 无法解析日志行", logLine)
		return nil, errors.New("无法解析日志行")
	}
	pid, err := strconv.Atoi(matches[5])
	if err != nil {
		return nil, errors.New("无法解析PID")
	}
	entry := LogEntry{
		Date:      matches[1],
		Time:      matches[2],
		Hostname:  matches[3],
		Process:   matches[4],
		PID:       pid,
		Level:     matches[6],
		Timestamp: matches[7],
		Message:   matches[8],
	}
	return &entry, nil
}

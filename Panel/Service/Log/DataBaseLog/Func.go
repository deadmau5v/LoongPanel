/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-30
 * 文件作用：
 */

package DataBaseLog

import (
	"LoongPanel/Panel/Service/Log"
	Log2 "LoongPanel/Panel/Service/PanelLog"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
)

type LogEntry struct {
	Date    string `json:"date"`
	Time    string `json:"time"`
	Level   string `json:"level"`
	Module  string `json:"module"`
	Content string `json:"content"`
}

func GetDataBaseLog() *Log.Log_ {

	log := &Log.Log_{
		Path: "/var/log/tidb.log",
		Name: "数据库日志",
		Ok:   true,
	}

	log.Ok = log.CheckLogExist()
	if !log.Ok {
		return nil
	}

	log.GetLog = func(line int) interface{} {
		file, err := os.Open(log.Path)
		if err != nil {
			Log2.ERROR("[面板日志] 打开日志文件失败: ", err.Error())
			return nil
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				Log2.ERROR("[面板日志] 关闭日志文件失败: ", err.Error())
			}
		}(file)

		all, err := io.ReadAll(file)
		if err != nil {
			Log2.ERROR("[面板日志] 读取日志文件失败: ", err.Error())
			return nil
		}

		allStr := string(all)
		allStrSplite := strings.Split(allStr, "\n")
		if len(allStrSplite) > line {
			allStrSplite = allStrSplite[len(allStrSplite)-line:]
			all = []byte(strings.Join(allStrSplite, "\n"))
		}

		var res []LogEntry
		for _, line := range allStrSplite {
			entry, err := parseDataBaseLog(line)
			if err != nil {
				continue
			}
			if entry != nil {
				res = append(res, *entry)
			}
		}

		return res
	}

	log.Struct = append(log.Struct, []map[string]string{
		{
			"title":     "日期",
			"dataIndex": "time",
			"key":       "1",
		}, {
			"title":     "时间",
			"dataIndex": "date",
			"key":       "1",
		},
		{
			"title":     "日志等级",
			"dataIndex": "level",
			"key":       "2",
		},
		{
			"title":     "模块",
			"dataIndex": "module",
			"key":       "3",
		},
		{
			"title":     "内容",
			"dataIndex": "content",
			"key":       "4",
		},
	})

	log.ClearLog = func() {
		if !log.Ok {
			return
		}
		err := os.Truncate(log.Path, 0)
		if err != nil {
			Log2.ERROR("[面板日志] 清空日志文件失败: ", err.Error())
		}
	}

	return log
}

// parseDataBaseLog 解析数据库日志
func parseDataBaseLog(logLine string) (*LogEntry, error) {
	pattern := `^\[(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3} \+\d{2}:\d{2})\] \[([A-Z]+)\] \[([\w\.]+:\d+)\] \[(.*?)\]`

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

	entry := LogEntry{
		Date:    matches[1][:10],
		Time:    matches[1][11:],
		Level:   matches[2],
		Module:  matches[3],
		Content: matches[4],
	}

	return &entry, nil
}

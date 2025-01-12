/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-29
 * 文件作用：面板日志
 */

package PanelLog

import (
	"LoongPanel/Panel/Service/Log"
	PanelLog "LoongPanel/Panel/Service/PanelLog"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
)

type LogEntry struct {
	Level   string `json:"level"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	Module  string `json:"module"`
	Content string `json:"content"`
}

func GetPanelLog() *Log.Log_ {

	log := &Log.Log_{
		Path: "./temp.log",
		Name: "面板日志",
		Ok:   true,
	}

	log.Ok = log.CheckLogExist()
	if !log.Ok {
		return nil
	}

	log.GetLog = func(line int) interface{} {
		file, err := os.Open(log.Path)
		if err != nil {
			PanelLog.ERROR("[面板日志] 打开日志文件失败: ", err.Error())
			return nil
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				PanelLog.ERROR("[面板日志] 关闭日志文件失败: ", err.Error())
			}
		}(file)

		all, err := io.ReadAll(file)
		if err != nil {
			PanelLog.ERROR("[面板日志] 读取日志文件失败: ", err.Error())
			return nil
		}

		allStr := string(all)
		allStrSplite := strings.Split(allStr, "\n")
		if len(allStrSplite) > line {
			allStrSplite = allStrSplite[len(allStrSplite)-line:]
		}

		var res []LogEntry
		for _, line := range allStrSplite {
			entry, err := parseLog(line)
			if err != nil {
				continue
			}
			if entry != nil {
				res = append(res, *entry)
			}
		}

		return res
	}

	log.ClearLog = func() {
		if !log.Ok {
			return
		}
		err := os.Truncate(log.Path, 0)
		if err != nil {
			PanelLog.ERROR("[面板日志] 清空日志文件失败: ", err.Error())
		}
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
		}, {
			"title":     "模块",
			"dataIndex": "module",
			"key":       "4",
		}, {
			"title":     "内容",
			"dataIndex": "content",
			"key":       "5",
		},
	})

	return log
}

// ParseLog 解析日志
func parseLog(logLine string) (*LogEntry, error) {
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
		return nil, errors.New("无法解析日志行")
	}

	entry := LogEntry{
		Level:   matches[1],
		Date:    matches[2],
		Time:    matches[3],
		Module:  matches[4],
		Content: matches[5],
	}

	return &entry, nil
}

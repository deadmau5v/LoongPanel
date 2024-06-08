/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-29
 * 文件作用：包管理工具日志
 */

package PkgLog

import (
	"LoongPanel/Panel/Service/Log"
	Log2 "LoongPanel/Panel/Service/PanelLog"
	"errors"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type LogEntry struct {
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func GetDnfLog() *Log.Log_ {
	path := "/var/log/dnf.log"
	name := "包管理工具"
	log := Log.Log_{
		Path: path,
		Name: name,
		Ok:   true,
	}

	if !log.CheckLogExist() {
		return nil
	}

	log.GetLog = func(line int) interface{} {
		if !log.Ok {
			Log2.DEBUG("[包管理日志] 错误跳过读取")
			return nil
		}
		file, err := os.Open(log.Path)
		if err != nil {
			Log2.ERROR("[包管理日志] 打开日志文件失败: ", err.Error())
			log.Ok = false
			return nil
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				Log2.ERROR("[包管理日志] 关闭日志文件失败: ", err.Error())
			}
		}(file)

		all, err := io.ReadAll(file)
		if err != nil {
			Log2.ERROR("[包管理日志] 读取日志文件失败: ", err.Error())
			log.Ok = false
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
			entry, err := parseLog(line)
			if err != nil {
				continue
			}
			if entry != nil {
				res = append(res, *entry)
			}
		}

		Log2.INFO("[包管理日志] 读取日志成功")

		return res
	}

	log.Struct = append(log.Struct, []map[string]string{
		{
			"title":     "日志等级",
			"dataIndex": "level",
			"key":       "1",
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

	log.ClearLog = func() {
		if !log.Ok {
			Log2.DEBUG("[包管理日志] 错误跳过清除")
			return
		}
		err := exec.Command("echo", ">", log.Path).Run()
		if err != nil {
			Log2.ERROR("[包管理日志] 清除日志失败: ", err.Error())
			return
		}
		Log2.INFO("[包管理日志] 清除日志成功")
	}

	return &log
}

// parseLog 解析日志
func parseLog(logLine string) (*LogEntry, error) {
	pattern := `^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{4}) (\w+) (.+)$`
	checkLogLine := strings.TrimSpace(logLine)
	if checkLogLine == "" {
		return nil, errors.New("空行")
	}

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(logLine)
	if matches == nil {
		return nil, errors.New("无法解析日志行")
	}

	entry := LogEntry{
		Timestamp: matches[1], // timestamp
		Level:     matches[2], // level
		Message:   matches[3], // message
	}

	return &entry, nil
}

//// GetYumLog 获取yum日志
//func GetYumLog() *Log.Log_ {
//	return createLog("/var/log/yum.log", "yum")
//}

//// GetAptLog 获取apt日志
//func GetAptLog() *Log.Log_ {
//	return createLog("/var/log/apt/history.log", "apt")
//}

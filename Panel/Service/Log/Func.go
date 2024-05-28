/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-21
 * 文件作用：统一日至模块
 */

package Log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
	"time"
)

var logColors = map[string]string{
	"INFO":  "\033[1;32m[INF]\033[0m %s",
	"WARN":  "\033[1;33m[WAR]\033[0m %s",
	"ERROR": "\033[1;31m[ERR]\033[0m %s",
	"DEBUG": "\033[0;33m[DEB]\033[0m %s", // 黄色（橙色）
}

func logWithColor(level string, args ...interface{}) {
	color, ok := logColors[level]
	if !ok {
		color = "%s"
	}
	output := make([]string, 0)
	time_ := time.Now().Format("2006-01-02 - 15:04:05 |")
	output = append(output, fmt.Sprintf("%v", time_))
	for _, arg := range args {
		output = append(output, fmt.Sprintf("%v", arg))
	}
	outputStr := strings.Join(output, " ")
	fmt.Println(fmt.Sprintf(color, outputStr))
}

func logToFile(args ...interface{}) {
	if !IsSaveToFile {
		return
	}
	_, err := os.Stat("./temp.log")
	if os.IsNotExist(err) {
		_, err := os.Create("./temp.log")
		if err != nil {
			ERROR("[日志模块]创建日志文件失败")
			return
		}
	}
	f, _ := os.OpenFile("./temp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			ERROR("[日志模块]关闭日志文件失败")
		}
	}(f)
	args = append([]interface{}{time.Now().Format("2006-01-02 - 15:04:05")}, args...)
	_, err = f.WriteString(fmt.Sprintln(args...))
	if err != nil {
		ERROR("[日志模块]写入日志文件失败")
		return
	}
}

func INFO(args ...interface{}) {
	logWithColor("INFO", args...)
	logToFile(args...)
}

func WARN(args ...interface{}) {
	logWithColor("WARN", args...)
	logToFile(args...)
}

func ERROR(args ...interface{}) {
	logWithColor("ERROR", args...)
	logToFile(args...)
}

func DEBUG(args ...interface{}) {
	if IsDebug {
		logWithColor("DEBUG", args...)
		logToFile(args...)
	}
}

func GinLogToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		logToFile(fmt.Sprintf("%s %s %s", c.ClientIP(), c.Request.Method, c.Request.URL.Path))
		c.Next()
	}
}

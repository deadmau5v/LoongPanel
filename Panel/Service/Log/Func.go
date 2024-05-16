/*
 * 创建人： deadmau5v
 * 创建时间： 2024-0-0
 * 文件作用：
 */

package Log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

var logColors = map[string]string{
	"INFO":  "[INFO]\033[1;34m%s\033[0m",
	"WARN":  "[WARN]\033[1;33m%s\033[0m",
	"ERROR": "[ERROR]\033[1;31m%s\033[0m",
	"DEBUG": "[DEBUG]\033[1;36m%s\033[0m",
}

func logWithColor(level string, args ...interface{}) {
	color, ok := logColors[level]
	if !ok {
		color = "%s"
	}
	fmt.Println(fmt.Sprintf(color, args...))
}

func logToFile(args ...interface{}) {
	_, err := os.Stat("./log.txt")
	if os.IsNotExist(err) {
		os.Create("./log.txt")
	}
	f, _ := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(fmt.Sprintln(args...))
}

func INFO(args ...interface{}) {
	logWithColor("INFO", args...)
	logToFile(args...)
}

func WARN(args ...interface{}) {
	logWithColor("WARN", args)
	logToFile(args)
}

func ERROR(args ...interface{}) {
	logWithColor("ERROR", args)
	logToFile(args)
}

func DEBUG(message string) {
	if isDebug {
		logWithColor("DEBUG", message)
		logToFile(message)
	}
}

func GinLogToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		logToFile(fmt.Sprintf("[%s] -> [%s] -> [%s]", c.ClientIP(), c.Request.Method, c.Request.URL.Path))
		c.Next()
	}
}

/*
 * 创建人： deadmau5v
 * 创建时间： 2024-6-4
 * 文件作用：路由页面 API
 */

package log

import (
	"LoongPanel/Panel/Service/Log"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetLogs 获取所有日志
func GetLogs(c *gin.Context) {
	allLog := Log.AllLog
	var allOkLog []string

	for _, log := range allLog {
		if log.Ok {
			allOkLog = append(allOkLog, log.Name)
		}
	}
	c.JSON(200, gin.H{
		"status": 0,
		"data":   allOkLog,
	})
}

// GetLog 获取日志
func GetLog(c *gin.Context) {
	name := c.Query("name")
	line := c.Query("line")
	if name == "" || line == "" {
		c.JSON(400, gin.H{"msg": "缺少参数", "status": 1})
		return
	}
	lineInt, err := strconv.Atoi(line)
	if err != nil {
		c.JSON(400, gin.H{"msg": "行数参数错误", "status": 1})
		return
	}

	for _, log := range Log.AllLog {
		if log.Name == name {
			output := log.GetLog(lineInt)
			if output != nil && log.Ok {
				c.JSON(200, gin.H{"data": output, "status": 0})
				return
			} else {
				c.JSON(400, gin.H{"msg": "日志获取错误", "status": 1})
				return
			}
		}
	}

	c.JSON(400, gin.H{"msg": "未找到日志", "status": 1})
}

// GetLogStruct 获取日志结构
func GetLogStruct(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(400, gin.H{"msg": "缺少参数"})
		return
	}

	for _, log := range Log.AllLog {
		if log.Name == name {
			c.JSON(200, gin.H{"data": log.Struct, "status": 0})
			return
		}
	}

	c.JSON(400, gin.H{"msg": "未找到日志", "status": 1})
}

// ClearLog 清空日志
func ClearLog(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(400, gin.H{"msg": "缺少参数", "status": 1})
		return
	}

	for _, log := range Log.AllLog {
		if log.Name == name {
			log.ClearLog()
			if log.Ok {
				c.JSON(200, gin.H{"msg": "清空成功", "status": 0})
			} else {
				c.JSON(400, gin.H{"msg": "清空失败", "status": 1})
			}
			return
		}
	}

	c.JSON(400, gin.H{"msg": "未找到日志", "status": 1})
}

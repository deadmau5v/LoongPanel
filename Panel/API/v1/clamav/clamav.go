/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供ClamAV相关的API 主要实现在 Service/Clamav 中
 */

package clamav

import (
	clamav "LoongPanel/Panel/Service/Clamav"
	"LoongPanel/Panel/Service/PanelLog"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func scan(c *gin.Context, scanDir bool, conn *websocket.Conn) {
	path := c.Query("path")
	scanType := c.Query("type")
	if strings.TrimSpace(path) == "" {
		path = "/"
	}
	PanelLog.INFO("[病毒扫描]", "开始扫描")
	var res *clamav.ScanResult
	var err error
	if scanType == "fast" {
		res, err = clamav.FastScan(conn)
	} else if scanType == "full" {
		res, err = clamav.FullScan(conn)
	} else {
		res, err = clamav.Scan(nil, []string{path}, scanDir, false)
	}
	if err != nil {
		PanelLog.ERROR("[病毒扫描]", err)
		if errors.Is(err, clamav.ErrorPath) {
			c.JSON(http.StatusOK, gin.H{"msg": "路径错误", "status": 1})
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": err.Error(), "status": 1})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res, "status": 0})
}

func ScanFile(c *gin.Context) {
	if c.Query("type") == "fast" {
		PanelLog.DEBUG("[病毒扫描]", "快速扫描")
		FastScan(c)
	} else if c.Query("type") == "full" {
		PanelLog.DEBUG("[病毒扫描]", "全盘扫描")
		FullScan(c)
	} else {
		PanelLog.DEBUG("[病毒扫描]", "扫描文件")
		scan(c, false, nil)
	}
}

func ScanDir(c *gin.Context) {
	PanelLog.DEBUG("[病毒扫描]", "扫描目录")
	scan(c, true, nil)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func FastScan(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		PanelLog.ERROR("[病毒扫描]", err)
		return
	}
	defer conn.Close() // 确保 WebSocket 连接在函数退出时关闭

	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"msg":"加载病毒库中..."}`))
	if err != nil {
		PanelLog.ERROR("[病毒扫描]", "发送消息失败:", err)
		return
	}

	// 执行扫描逻辑...
	// 假设 scan 返回扫描结果
	result, err := clamav.FastScan(conn)
	if err != nil {
		errMsg, err := json.Marshal(gin.H{"error": err.Error()})
		if err != nil {
			PanelLog.ERROR("[病毒扫描]", "转Json失败:", err)
		}
		conn.WriteMessage(websocket.TextMessage, errMsg)
		return
	}

	// 发送扫描结果
	resultMsg, err := json.Marshal(gin.H{
		"data":   result,
		"status": 0,
	})
	PanelLog.DEBUG("[病毒扫描]", "结果:", string(resultMsg))
	if err != nil {
		PanelLog.ERROR("[病毒扫描]", "转Json失败:", err)
	}
	conn.WriteMessage(websocket.TextMessage, resultMsg)
}

func FullScan(c *gin.Context) {
	// 升级 HTTP 连接到 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		PanelLog.ERROR("[病毒扫描]", "WebSocket连接升级失败:", err)
		return
	}
	defer conn.Close() // 确保 WebSocket 连接在函数退出时关闭

	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"msg":"加载病毒库中..."}`))
	if err != nil {
		PanelLog.ERROR("[病毒扫描]", "发送消息失败:", err)
		return
	}

	// 执行全盘扫描
	result, err := clamav.FullScan(conn)
	if err != nil {
		errMsg, err := json.Marshal(gin.H{"error": err.Error()})
		if err != nil {
			PanelLog.ERROR("[病毒扫描]", "转Json失败:", err)
		}
		conn.WriteMessage(websocket.TextMessage, errMsg)
		return
	}

	// 发送扫描结果
	resultMsg, err := json.Marshal(gin.H{
		"data":   result,
		"status": 0,
	})
	if err != nil {
		PanelLog.ERROR("[病毒扫描]", "转Json失败:", err)
	}
	conn.WriteMessage(websocket.TextMessage, resultMsg)
}

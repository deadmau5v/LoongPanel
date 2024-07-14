package clamav

import (
	config "LoongPanel/Panel/Service/Config"
	"LoongPanel/Panel/Service/Cron"
	"LoongPanel/Panel/Service/Database"
	notice "LoongPanel/Panel/Service/Notice"
	"LoongPanel/Panel/Service/PanelLog"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
)

var (
	ErrorPath = errors.New("path is not a file or directory")
)

var (
	output    []byte
	scan_over bool = false
	cronID    cron.EntryID
)

type WsReaderWriter struct {
	*websocket.Conn
}

type ScanResult struct {
	DataRead           string `json:"data_read"`           // 扫描读取的数据量
	DataScanned        string `json:"data_scanned"`        // 扫描扫描的数据量
	EndDate            string `json:"end_date"`            // 扫描结束时间
	EngineVersion      string `json:"engine_version"`      // 引擎版本
	InfectedFiles      string `json:"infected_files"`      // 感染的文件数量
	KnownViruses       string `json:"known_viruses"`       // 已知病毒数量
	ScannedDirectories string `json:"scanned_directories"` // 扫描的目录数量
	ScannedFiles       string `json:"scanned_files"`       // 扫描的文件数量
	StartDate          string `json:"start_date"`          // 扫描开始时间
	Time               string `json:"time"`                // 扫描耗时
}

type WriteOutput struct {
}

func (w *WriteOutput) Write(p []byte) (n int, err error) {
	if scan_over {
		return 0, nil // 如果扫描已结束，不再写入
	}
	output = append(output, p...)
	if strings.Contains(string(p), "SCAN SUMMARY") {
		scan_over = true
		return len(p), nil // 确保返回正确的写入字节数
	}
	return len(p), nil
}

func (w *WsReaderWriter) Write(p []byte) (n int, err error) {
	if scan_over {
		PanelLog.DEBUG("[调试]", "尝试写入但扫描已结束")
		return 0, nil // 如果扫描已结束，不再写入
	}
	output = append(output, p...)
	if strings.Contains(string(p), "SCAN SUMMARY") {
		scan_over = true
		PanelLog.DEBUG("[调试]", "扫描总结已写入，标记扫描结束")
		return len(p), nil // 确保返回正确的写入字节数
	}
	if w.Conn == nil || w.Conn.UnderlyingConn() == nil {
		PanelLog.DEBUG("[错误]", "WebSocket连接已关闭")
		return 0, fmt.Errorf("WebSocket连接已关闭")
	}
	// 尝试使用连接进行操作，例如写入空数据，来检测连接是否有效
	_, err = w.Conn.UnderlyingConn().Write([]byte{})
	if err != nil {
		PanelLog.DEBUG("[错误]", "WebSocket连接已关闭或出现其他错误:", err)
		return 0, fmt.Errorf("WebSocket连接已关闭或出现其他错误: %w", err)
	}
	writer, err := w.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		PanelLog.DEBUG("[错误]", "获取WebSocket写入器失败:", err)
		return 0, err
	}
	defer writer.Close()
	n, err = writer.Write(p)
	if err != nil {
		PanelLog.DEBUG("[错误]", "WebSocket写入失败:", err)
		return n, err // 处理写入错误
	}
	// PanelLog.DEBUG("[调试]", "WebSocket写入成功:", n, "字节")
	return n, nil
}

func (w *WsReaderWriter) Read(p []byte) (n int, err error) {
	var msgType int
	var reader io.Reader
	for {
		msgType, reader, err = w.Conn.NextReader()
		if err != nil {
			return 0, err
		}
		if msgType != websocket.TextMessage {
			continue
		}
		return reader.Read(p)
	}
}

// Check 检查文件路径是否存在 是否扫描目录
func Check(filePaths []string, scanDir bool) error {
	for _, filePath := range filePaths {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return fmt.Errorf("CheckFile -> %w", err)
		}

		if scanDir != fileInfo.IsDir() {
			return fmt.Errorf("CheckFile -> %w", ErrorPath)
		}
	}

	return nil
}

// Parse 解析扫描结果
func Parse(output string) (*ScanResult, error) {
	re := regexp.MustCompile(`(?m)Known viruses: (\d+)\nEngine version: ([\d.]+)\nScanned directories: (\d+)\nScanned files: (\d+)\nInfected files: (\d+)\nData scanned: ([\d.]+) MB\nData read: ([\d.]+) MB \(ratio [\d.]+:[\d.]+\)\nTime: ([\d.]+) sec \(\d+ m \d+ s\)\nStart Date: (\d{4}:\d{2}:\d{2} \d{2}:\d{2}:\d{2})\nEnd Date:   (\d{4}:\d{2}:\d{2} \d{2}:\d{2}:\d{2})`)
	matches := re.FindStringSubmatch(output)

	if len(matches) != 11 {
		return nil, fmt.Errorf("Parse -> %w", errors.New("output format error"))
	}

	summary := ScanResult{
		KnownViruses:       matches[1],
		EngineVersion:      matches[2],
		ScannedDirectories: matches[3],
		ScannedFiles:       matches[4],
		InfectedFiles:      matches[5],
		DataScanned:        matches[6],
		DataRead:           matches[7],
		Time:               matches[8],
		StartDate:          matches[9],
		EndDate:            matches[10],
	}

	return &summary, nil
}

func getOutput() string {
	return string(output)
}

// Scan 扫描文件或目录
func Scan(c *websocket.Conn, args []string, scanDir, skipCheck bool) (*ScanResult, error) {
	PanelLog.DEBUG("[调试]", "开始扫描")

	output = []byte{}
	defer func() {
		scan_over = false
	}()
	if !skipCheck {
		err := Check(args, scanDir)
		if err != nil {
			return nil, fmt.Errorf("ScanFile -> %w", err)
		}
	}

	PanelLog.DEBUG("[调试]", "创建连接")
	var conn *WsReaderWriter
	if c != nil {
		conn = &WsReaderWriter{Conn: c}
	}

	PanelLog.DEBUG("[调试]", "创建命令")
	cmd := exec.Command("clamscan", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if conn != nil {
		cmd.Stdout = conn
	} else {
		cmd.Stdout = &WriteOutput{}
	}
	// 执行命令
	PanelLog.DEBUG("[调试]", "执行命令")
	err := cmd.Run()
	if err != nil {
		if strings.Contains(stderr.String(), "short write") {
			return nil, fmt.Errorf("ScanFile:Execution -> %w: %s", err, stderr.String())
		}
		return nil, fmt.Errorf("ScanFile:Execution -> %w: %s", err, stderr.String())
	}
	// 解析结果
	PanelLog.DEBUG("[调试]", "开始解析output")
	summary, err := Parse(string(getOutput()))
	if err != nil {
		return nil, fmt.Errorf("ScanFile:Parse -> %w", err)
	}

	var settings []notice.UserNotificationSetting
	Database.DB.Preload("User").Find(&settings)
	for _, v := range settings {
		if v.ClamAVScanNotify {
			notice.SendMail(v.User.Mail, "病毒扫描结果", fmt.Sprintf("扫描结果: %s", summary.InfectedFiles))
		}
	}

	return summary, nil
}

// FastScan 快速扫描 扫描系统关键位置
func FastScan(conn *websocket.Conn) (*ScanResult, error) {
	args := []string{
		"/tmp",
		"/var/tmp",
		"/dev/shm",
		"/bin",
		"/etc",
		"/boot",
		"/home",
		"/root/.bashrc",
	}

	return Scan(conn, args, true, true)
}

func FullScan(conn *websocket.Conn) (*ScanResult, error) {
	args := []string{
		"-r",
		"/",
	}

	return Scan(conn, args, true, true)
}

func SetCronScan(duration time.Duration) error {
	if cronID != 0 {
		Cron.Cron.Remove(cronID)
	}
	id, err := Cron.Cron.AddFunc(Cron.DurationToCron(duration), func() {
		ScanResult, err := FastScan(nil)
		if err != nil {
			return
		}
		PanelLog.INFO("[病毒扫描]", "扫描结果:", ScanResult.InfectedFiles, "个病毒")
	})
	if err != nil {
		return fmt.Errorf("SetCronScan -> %w", err)
	}
	cronID = id
	return nil
}

func init() {
	if config.Config.Clamav.CronScan {
		if config.Config.Clamav.CronScanTime < 4*time.Hour {
			// 太短的扫描时间太影响性能
			config.Config.Clamav.CronScanTime = 4 * time.Hour
		}
		SetCronScan(config.Config.Clamav.CronScanTime)
	}
}

package clamav

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	ErrorPath = errors.New("path is not a file or directory")
)

var (
	output    []byte
	scan_over bool = false
)

type WsReaderWriter struct {
	*websocket.Conn
}

type ScanResult struct {
	DataRead           string `json:"data_read"`
	DataScanned        string `json:"data_scanned"`
	EndDate            string `json:"end_date"`
	EngineVersion      string `json:"engine_version"`
	InfectedFiles      string `json:"infected_files"`
	KnownViruses       string `json:"known_viruses"`
	ScannedDirectories string `json:"scanned_directories"`
	ScannedFiles       string `json:"scanned_files"`
	StartDate          string `json:"start_date"`
	Time               string `json:"time"`
}

func (w *WsReaderWriter) Write(p []byte) (n int, err error) {
	output = append(output, p...)
	if scan_over {
		return 0, nil
	}
	if strings.Contains(string(p), "SCAN SUMMARY") {
		scan_over = true
		return 0, nil
	}
	writer, err := w.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}
	defer writer.Close()
	return writer.Write(p)
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

	defer func() {
		scan_over = false
		output = []byte{}
	}()
	if !skipCheck {
		err := Check(args, scanDir)
		if err != nil {
			return nil, fmt.Errorf("ScanFile -> %w", err)
		}
	}

	var conn *WsReaderWriter
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if c != nil {
		conn = &WsReaderWriter{Conn: c}
	}

	cmd := exec.Command("clamscan", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if conn != nil {
		cmd.Stdout = conn
	}

	err := cmd.Run()
	if err != nil {
		if strings.Contains(err.Error(), "short write") {
			return nil, nil
		}
		return nil, fmt.Errorf("ScanFile:Execution -> %w: %s", err, stderr.String())
	}

	summary, err := Parse(getOutput())
	if err != nil {
		return nil, fmt.Errorf("ScanFile:Parse -> %w", err)
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

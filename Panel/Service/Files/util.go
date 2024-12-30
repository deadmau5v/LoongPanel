/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：文件操作工具类
 */

package Files

import (
	"LoongPanel/Panel/Service/PanelLog"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CheckFileName 检查文件名是否合法
func CheckFileName(path string) bool {
	words := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, word := range words {
		if strings.Contains(path, word) {
			return false
		}
	}
	return true
}

// RuneIndex 查找字符在切片中的位置
func RuneIndex(runes []rune, str string) int {
	for i, v := range runes {
		if string(v) == str {
			return i
		}
	}
	return -1
}

// RuneIndexR 从右向左查找
//func RuneIndexR(runes []rune, str string) int {
//	for i := len(runes) - 1; i >= 0; i-- {
//		if string(runes[i]) == str {
//			return i
//		}
//	}
//	return -1
//}

// getFilePath 获取文件的路径
func getFilePath(path string) string {
	return filepath.Dir(path)
}

// getFileName 获取文件名
func getFileName(path string) string {
	if strings.TrimSpace(path) == "" {
		return ""
	}
	return filepath.Base(path)
}

// copyFileConflictRename 复制文件冲突时改名
func copyFileConflictRename(path string) string {
	filePath := getFilePath(path)
	fileName := getFileName(path)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "_copy" + filepath.Ext(fileName)
	return filepath.Join(filePath, fileName)
}

// copyDirConflictRename 复制文件夹冲突将目标文件夹名改为原文件夹名_copy
func copyDirConflictRename(path string) string {
	filePath := getFilePath(path)
	fileName := getFileName(path)
	fileName = fileName + "_copy"
	return filepath.Join(filePath, fileName)
}

// checkFilePath 检查文件路径是否合法
func checkFilePath(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	words := []string{"\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, word := range words {
		if strings.Contains(path, word) {
			return false
		}
	}
	return true
}

// isDir 判断是否是目录
func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// copyFile 复制文件 已经确保参数合法
func copyFile(srcPath, targetPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		PanelLog.ERROR("[文件管理]", "打开源文件失败", srcPath, err.Error())
		return fmt.Errorf("Copy:os.Open -> %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(targetPath)
	if err != nil {
		PanelLog.ERROR("[文件管理]", "创建目标文件失败", targetPath, err.Error())
		return fmt.Errorf("Copy:os.Create -> %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		PanelLog.ERROR("[文件管理]", "复制文件失败", srcPath, err.Error())
		return fmt.Errorf("Copy:io.Copy -> %w", err)
	}
	return nil
}

// isTar 判断是否是tar文件
func isTar(path string) bool {
	return filepath.Ext(path) == ".tar"
}

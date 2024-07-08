/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供文件操作的实现
 */

package Files

import (
	"LoongPanel/Panel/Service/PanelLog"
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Dir 读取目录
func Dir(path string) ([]File, error) {
	path, err := filepath.Abs(path)
	path = filepath.Clean(path)
	if err != nil {
		PanelLog.ERROR("[文件管理]", err.Error())
		return nil, err
	}

	files := make([]File, 0)
	readDir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// 解析文件
	for _, file_ := range readDir {
		fileStat, _ := file_.Info()

		file := File{Name: file_.Name()}
		if path != "/" {
			file.Path = path + string(os.PathSeparator) + file_.Name()
		} else {
			file.Path = path + file_.Name()
		}
		// 路径

		file.User, file.Group = getUidGid(fileStat) // 所属用户 所属组

		file.Size = fileStat.Size()                                    // 大小
		file.Mode = fileStat.Mode()                                    // 权限
		file.Time = fileStat.ModTime()                                 // 时间
		file.IsDir = fileStat.IsDir()                                  // 是否是目录
		file.IsHidden = file_.Name()[0] == '.'                         // 是否隐藏
		file.IsLink = fileStat.Mode()&os.ModeSymlink == os.ModeSymlink // 是否为链接
		file.ShowEdit = true                                           // 显示编辑按钮
		file.ShowTime = true                                           // 显示时间
		file.ShowSize = true                                           // 显示大小

		if file.IsDir {
			file.ShowSize = false // 显示大小
			file.ShowTime = false // 显示时间
		}

		// 文件类型
		if !file.IsDir && strings.Contains(file.Name[1:], ".") {
			runeName := []rune(file.Name)
			file.Ext = string(runeName[RuneIndex(runeName, ".")+1:])
		} else {
			file.Ext = ""
		}
		files = append(files, file)
	}

	return files, err
}

// Content 读取文件内容
func Content(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		PanelLog.ERROR("[文件管理]", err.Error())
		return "", err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		PanelLog.ERROR("[文件管理]", err.Error())
		return "", err
	}
	fileStr := string(file)
	return fileStr, err
}

// Delete 删除
func Delete(path string) error {
	// 获取绝对路径
	path, err := filepath.Abs(path)
	if err != nil {
		PanelLog.ERROR("[文件管理]", err.Error())
		return err
	}
	// 删除文件
	err = os.Remove(path)
	if err != nil {
		PanelLog.ERROR("[文件管理]", err.Error())
		return err
	}
	PanelLog.INFO("[文件管理]", "删除文件", path)

	return nil
}

// Copy 复制文件
func Copy(path string, dest string) error {
	// 检查路径
	if !checkFilePath(path) || !checkFilePath(dest) {
		PanelLog.ERROR("[文件管理]", "路径不合法")
		return fmt.Errorf("Copy:checkFilePath -> %w", errors.New("路径不合法"))
	}

	// 获取绝对路径
	path, err := filepath.Abs(path)
	if err != nil {
		PanelLog.ERROR("[文件管理]", err.Error())
		return fmt.Errorf("Copy:filepath.Abs -> %w", err)
	}
	dest, err = filepath.Abs(dest) // 这里应该是 dest
	if err != nil {
		PanelLog.ERROR("[文件管理]", err.Error())
		return fmt.Errorf("Copy:filepath.Abs -> %w", err)
	}

	// 文件夹复制
	if isDir(path) {
		if !isDir(dest) {
			PanelLog.ERROR("[文件管理]", "目标路径不是文件夹")
			return fmt.Errorf("Copy:isDir -> %w", errors.New("目标路径不是文件夹"))
		}

		if getFilePath(dest) == getFilePath(path) {
			dest = copyDirConflictRename(dest)
		}

		// 创建目标文件夹
		err = os.MkdirAll(dest, 0755)
		if err != nil {
			PanelLog.ERROR("[文件管理]", "创建目标文件夹失败", err.Error())
			return fmt.Errorf("Copy:os.MkdirAll -> %w", err)
		}

		// 遍历源文件夹
		err = filepath.Walk(path, func(srcPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 计算目标文件路径
			relPath, err := filepath.Rel(path, srcPath)
			if err != nil {
				return err
			}
			targetPath := filepath.Join(dest, relPath)

			// 如果是目录则创建目录
			if info.IsDir() {
				err = os.MkdirAll(targetPath, info.Mode())
				if err != nil {
					PanelLog.ERROR("[文件管理]", "创建目录失败", targetPath, err.Error())
					return fmt.Errorf("Copy:os.MkdirAll -> %w", err)
				}
			} else {
				// 复制文件
				copyFile(srcPath, targetPath)
			}
			return nil
		})

		if err != nil {
			PanelLog.ERROR("[文件管理]", "复制文件夹失败", err.Error())
			return fmt.Errorf("Copy:filepath.Walk -> %w", err)
		}
	} else {
		// 打开源文件
		srcFile, err := os.Open(path)
		if err != nil {
			PanelLog.ERROR("[文件管理]", "打开源文件失败", err.Error())
			return fmt.Errorf("Copy:os.Open -> %w", err)
		}
		defer srcFile.Close()

		// 如果文件复制到同一目录下，则改名 加上_copy
		// 详细看 copyFileConflictRename 函数
		if isDir(dest) {
			dest = filepath.Join(dest, getFileName(path))
		}
		if getFilePath(dest) == getFilePath(path) {
			dest = copyFileConflictRename(dest)
		}

		// 创建目标文件
		destFile, err := os.Create(dest)
		if err != nil {
			PanelLog.ERROR("[文件管理]", "创建目标文件失败", err.Error())
			return fmt.Errorf("Copy:os.Create -> %w", err)
		}

		// 复制文件
		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			PanelLog.ERROR("[文件管理]", "复制文件失败", err.Error())
			return fmt.Errorf("Copy:io.Copy -> %w", err)
		}
	}

	PanelLog.INFO("[文件管理]", "文件复制成功", path, "到", dest)
	return nil
}

// Move 移动
func Move(path string, dest string) error {
	// 检查路径
	if !checkFilePath(path) || !checkFilePath(dest) {
		PanelLog.ERROR("[文件管理]", "路径不合法")
		return fmt.Errorf("Move:checkFilePath -> %w", errors.New("路径不合法"))
	}

	if isDir(path) {
		err := os.Rename(path, dest)
		if err != nil {
			PanelLog.ERROR("[文件管理]", "移动文件夹失败", path, "到", dest, err.Error())
			return fmt.Errorf("Move:os.Rename -> %w", err)
		}
	} else {
		err := os.Rename(path, dest)
		if err != nil {
			PanelLog.ERROR("[文件管理]", "移动文件失败", path, "到", dest, err.Error())
			return fmt.Errorf("Move:os.Rename -> %w", err)
		}
	}
	PanelLog.INFO("[文件管理]", "文件移动成功", path, "到", dest)
	return nil
}

// Rename 重命名
func Rename(path string, newName string) error {
	// 检查路径
	if !checkFilePath(path) {
		PanelLog.ERROR("[文件管理]", "路径不合法")
		return fmt.Errorf("Rename:checkFilePath -> %w", errors.New("路径不合法"))
	}
	if !checkFilePath(newName) {
		return fmt.Errorf("Rename:checkNewFileName -> %w", errors.New("文件名不可用"))
	}

	newPath := filepath.Join(getFilePath(path), newName)
	Move(path, newPath)
	return nil
}

// Compress 压缩
func Compress(path string) error {
	if !checkFilePath(path) {
		PanelLog.ERROR("[文件管理]", "压缩路径不合法", path)
		return fmt.Errorf("Compress:checkFilePath -> %w", errors.New("压缩路径不合法"))
	}

	// 创建压缩文件
	tarFile := path + ".tar"
	file, err := os.Create(tarFile)
	if err != nil {
		PanelLog.ERROR("[文件管理]", "创建压缩文件失败", tarFile, err.Error())
		return fmt.Errorf("Compress:os.Create -> %w", err)
	}
	defer file.Close()

	tarWriter := tar.NewWriter(file)
	defer tarWriter.Close()

	info, err := os.Stat(path)
	if err != nil {
		PanelLog.ERROR("[文件管理]", "获取文件信息失败", path, err.Error())
		return fmt.Errorf("Compress:os.Stat -> %w", err)
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(path)
	}

	err = filepath.Walk(path, func(fileP string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fileInfo, "")
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(fileP, path))
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !fileInfo.IsDir() {
			fileContents, err := os.Open(fileP)
			if err != nil {
				return err
			}
			defer fileContents.Close()

			if _, err := io.Copy(tarWriter, fileContents); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		PanelLog.ERROR("[文件管理]", "压缩过程中出错", path, err.Error())
		return fmt.Errorf("Compress:filepath.Walk -> %w", err)
	}

	PanelLog.INFO("[文件管理]", "文件压缩成功", path, "到", tarFile)
	return nil
}

// Decompress 解压
func Decompress(path string) error {
	// 检查路径
	if !checkFilePath(path) {
		PanelLog.ERROR("[文件管理]", "解压路径不合法", path)
		return fmt.Errorf("Decompress:checkFilePath -> %w", errors.New("解压路径不合法"))
	}

	// 获取绝对路径
	path, err := filepath.Abs(path)
	if err != nil {
		PanelLog.ERROR("[文件管理]", err.Error())
		return fmt.Errorf("Decompress:filepath.Abs -> %w", err)
	}

	// 判断是否为文件夹
	if isDir(path) {
		PanelLog.ERROR("[文件管理]", "解压路径是文件夹", path)
		return fmt.Errorf("Decompress:isDir -> %w", errors.New("解压路径是文件夹"))
	}

	// 判断是否是tar文件
	if !isTar(path) {
		PanelLog.ERROR("[文件管理]", "解压路径不是tar文件", path)
		return fmt.Errorf("Decompress:isTar -> %w", errors.New("解压路径不是tar文件"))
	}

	// 打开tar文件
	tarFile, err := os.Open(path)
	if err != nil {
		PanelLog.ERROR("[文件管理]", "打开tar文件失败", path, err.Error())
		return fmt.Errorf("Decompress:os.Open -> %w", err)
	}
	defer tarFile.Close()

	// 创建tar阅读器
	tarReader := tar.NewReader(tarFile)

	// 解压缩每个文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 文件结束
		}
		if err != nil {
			PanelLog.ERROR("[文件管理]", "读取tar文件失败", path, err.Error())
			return fmt.Errorf("Decompress:tarReader.Next -> %w", err)
		}

		// 目标文件路径
		targetPath := filepath.Join(filepath.Dir(path), header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// 创建文件夹
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				PanelLog.ERROR("[文件管理]", "创建文件夹失败", targetPath, err.Error())
				return fmt.Errorf("Decompress:os.MkdirAll -> %w", err)
			}
		case tar.TypeReg:
			// 创建文件
			outFile, err := os.Create(targetPath)
			if err != nil {
				PanelLog.ERROR("[文件管理]", "创建文件失败", targetPath, err.Error())
				return fmt.Errorf("Decompress:os.Create -> %w", err)
			}
			defer outFile.Close()

			// 写入文件
			if _, err := io.Copy(outFile, tarReader); err != nil {
				PanelLog.ERROR("[文件管理]", "写入文件失败", targetPath, err.Error())
				return fmt.Errorf("Decompress:io.Copy -> %w", err)
			}
		}
	}

	return nil
}

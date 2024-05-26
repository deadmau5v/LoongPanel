/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：提供文件操作的实现
 */

package Files

import (
	"LoongPanel/Panel/Service/Log"
	"os"
	"path/filepath"
	"strings"
)

func Dir(path string) ([]File, error) {
	path, err := filepath.Abs(path)
	path = filepath.Clean(path)
	if err != nil {
		Log.ERROR(err.Error())
		return nil, err
	}

	back := NewObj()
	// 上一级
	back.Name = ".."
	back.IsDir = true
	back.IsLink = false
	back.Path = filepath.Join(path, "..")
	back.ShowEdit = false
	back.ShowSize = false
	back.ShowTime = false

	files := make([]File, 0)
	if path != "/" && path != "" {
		files = append(files, *back)
	}

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

		file.User, file.Group = getUidGid() // 所属用户 所属组

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
			file.Ext = string(runeName[RuneIndex(runeName, "."):])
		} else {
			file.Ext = ""
		}
		files = append(files, file)
	}

	return files, err
}

package Files

import (
	"os"
	"strings"
)

func Dir(path string) ([]File, error) {

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

		fileInfo := fileStat.Sys()
		file.User, file.Group = getUidGid(&fileInfo) // 所属用户 所属组

		file.Size = fileStat.Size()                                // 大小
		file.Mod = fileStat.Mode().String()                        // 权限
		file.Time = fileStat.ModTime().Format("2024-04-1 2:15:05") // 时间
		file.IsDir = fileStat.IsDir()                              // 是否是目录
		file.IsHidden = file_.Name()[0] == '.'                     // 是否隐藏
		// 文件类型
		if !file.IsDir && strings.Contains(".", file.Name[1:]) {
			file.Ext = file.Name[strings.LastIndex(file.Name, ".")+1:]
		} else {
			file.Ext = ""
		}
		files = append(files, file)
	}

	return files, err
}

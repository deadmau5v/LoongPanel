package Files

import (
	"os"
	"strconv"
	"syscall"
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
		file.Path = path + string(os.PathSeparator) + file_.Name()
		file.User = strconv.Itoa(int(fileStat.Sys().(*syscall.Stat_t).Uid))
		file.Group = strconv.Itoa(int(fileStat.Sys().(*syscall.Stat_t).Gid))
		file.Size = fileStat.Size()
		file.Mod = fileStat.Mode().String()
		file.Time = fileStat.ModTime().Format("2024-04-1 2:15:05")
		file.IsDir = fileStat.IsDir()
		file.IsHidden = file_.Name()[0] == '.'
	}

	return files, err
}

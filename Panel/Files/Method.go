package Files

import (
	"errors"
	"io"
	"os"
)

// Delete 删除当前文件
func (f *File) Delete() error {
	err := os.Remove(f.Path)
	if err != nil {
		return err
	}
	return nil
}

// New 创建一个新文件 or 路径
func (f *File) New() error {
	// 创建文件
	if f.IsDir {
		err := os.MkdirAll(f.Path, f.Mode)
		if err != nil {
			return err
		}
		return nil
	} else {
		file, err := os.Create(f.Path)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)
	}
	// 设置权限

	err := os.Chmod(f.Path, f.Mode)
	if err != nil {
		return err
	}
	return nil
}

// setName 设置文件名 同步修改文件路径
func (f *File) setName(name string) error {
	ok := CheckFileName(name)
	if !ok {
		return errors.New("文件名不合法")
	}
	// 兼容中文
	RunePath := []rune(f.Path)
	RuneName := []rune(f.Name)
	RunePath = RunePath[:len(RunePath)-len(RuneName)]
	f.Name = name
	f.Path = string(RunePath) + f.Name

	return nil
}

// Copy 复制文件
func (f *File) Copy(path string) error {
	if f.IsDir {
		// Todo: 复制文件夹
	} else {
		// 读取文件
		file, err := os.Open(f.Path)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)

		// 创建文件
		newFileObj := File{
			Name:  f.Name,
			Path:  path,
			Mode:  f.Mode,
			IsDir: f.IsDir,
		}
		err = newFileObj.New()
		if err != nil {
			return err
		}
		newFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func(newFile *os.File) {
			err := newFile.Close()
			if err != nil {
				return
			}
		}(newFile)

		// 复制文件
		_, err = io.Copy(newFile, file)
		if err != nil {
			return err
		}
	}
	return nil
}

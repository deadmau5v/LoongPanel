package Files

import (
	"LoongPanel/Panel/System"
	"errors"
	"github.com/spf13/afero"
	"os"
	"os/exec"
	"time"
)

// Delete 删除当前文件
func (f *File) Delete() error {
	fs := afero.NewOsFs()
	if f.IsDir {
		err := fs.RemoveAll(f.Path)
		if err != nil {
			return err
		}
	} else {
		err := fs.Remove(f.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create 创建一个新文件 or 路径
func (f *File) Create() error {
	fs := afero.NewOsFs()
	// 创建文件
	if f.IsDir {
		err := fs.MkdirAll(f.Path, f.Mode)
		if err != nil {
			return err
		}
	} else {
		_, err := fs.Create(f.Path)
		if err != nil {
			return err
		}
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

// NewFile 返回默认对象
func NewFile() *File {
	return &File{
		Name:  "NewFile",
		Size:  0,
		Time:  time.Now(),
		IsDir: false,
		Mode:  os.ModePerm,
		Path:  "/NewFile",
		Ext:   "",
	}
}

// Copy 复制文件
func (f *File) Copy(dst *File) error {
	if f.Path == dst.Path {
		return errors.New("源文件与目标文件路径相同")
	}
	fs := afero.NewOsFs()
	src, err := fs.Stat(f.Path)
	if err != nil {
		return err
	}

	if f.IsDir {
		// 复制目录
		err := fs.MkdirAll(dst.Path+string(os.PathSeparator), src.Mode())
		if err != nil {
			return err
		}
		switch System.Data.OSName {

		case "linux":
			{
				err := exec.Command("cp", "-r", f.Path, dst.Path+string(os.PathSeparator)).Run()
				if err != nil {
					return err
				}
			}
		case "windows":
			{
				err := exec.Command("xcopy", f.Path, dst.Path, "/s", "/e").Run()
				if err != nil {
					return err
				}
			}

		}
	} else {
		// 复制文件
		switch System.Data.OSName {
		case "windows":
			{
				err := exec.Command("copy", f.Path, dst.Path).Run()
				if err != nil {
					return err
				}
			}
		case "linux":
			{
				err := exec.Command("cp", "-f", f.Path, dst.Path).Run()
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}

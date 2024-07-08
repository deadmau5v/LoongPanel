/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：统一标准化操作Files对象的方法
 */

package Files

import (
	"LoongPanel/Panel/Service/System"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/afero"
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

// SetName 设置文件名 同步修改文件路径
func (f *File) SetName(name string) error {
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

// SetPath 设置文件路径
func (f *File) SetPath(path string) *File {
	f_ := f.CopyObj()
	if path[len(path)-1] == os.PathSeparator {
		f_.Path = path + f_.Name
	} else {
		f_.Path = path + string(os.PathSeparator) + f_.Name
	}
	return f_
}

// GetDir 获取文件的纯目录 (不包含文件名)
func (f *File) GetDir() string {
	name := []rune(f.Name)
	dir := []rune(f.Path)
	dir = dir[:len(dir)-len(name)]
	return string(dir)
}

// NewObj 返回默认对象
func NewObj() *File {
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

// Rename 重命名文件
func (f *File) Rename(dst *File) error {
	if f.Path == dst.Path {
		return errors.New("源文件与目标文件路径相同")
	} else if strings.Contains(dst.Path, f.Path) {
		return fmt.Errorf("目标文件路径包含源文件路径 %s -> %s", f.Path, dst.Path)
	}
	fs := afero.NewOsFs()
	err := fs.Rename(f.Path, dst.Path)
	if err != nil {
		return err
	}
	return nil
}

// CopyObj 复制对象
func (f *File) CopyObj() *File {
	return &File{
		Name:  f.Name,
		Size:  f.Size,
		Time:  f.Time,
		IsDir: f.IsDir,
		Mode:  f.Mode,
		Path:  f.Path,
		Ext:   f.Ext,
	}
}

// Move 移动文件
func (f *File) Move(dst *File) error {
	return f.Rename(dst)
}

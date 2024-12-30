//go:build linux

/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：linux 独有的文件操作支持
 */

package Files

import (
	"os"
	"syscall"
)

func getUidGid(fileStat os.FileInfo) (uint32, uint32) {
	return fileStat.Sys().(*syscall.Stat_t).Uid, fileStat.Sys().(*syscall.Stat_t).Gid
}

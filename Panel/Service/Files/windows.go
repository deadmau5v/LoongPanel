//go:build windows

/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：windows 独有的文件操作
 */

package Files

import (
	"fmt"
	"syscall"
)

func getUidGid(fileStat any) (uint32, uint32) {
	if false { // 调试
		attributes := fileStat.(syscall.Win32finddata)
		fmt.Println(attributes)
	}
	return 0, 0
}

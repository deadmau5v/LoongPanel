//go:build windows

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

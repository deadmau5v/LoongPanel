//go:build linux

package Files

import (
	"syscall"
)

func getUidGid(fileStat any) (uint32, uint32) {
	return fileStat.(*syscall.Stat_t).Uid, fileStat.(*syscall.Stat_t).Gid
}

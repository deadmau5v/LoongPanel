//go:build windows

/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：windows 独有的文件操作
 */

package Files

func getUidGid() (uint32, uint32) {
	// Windows 平台暂不支持
	return 0, 0
}

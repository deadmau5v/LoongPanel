/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：文件操作工具类
 */

package Files

import (
	"strings"
)

func CheckFileName(path string) bool {
	words := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, word := range words {
		if strings.Contains(path, word) {
			return false
		}
	}

	return true
}

func RuneIndex(runes []rune, str string) int {
	for i, v := range runes {
		if string(v) == str {
			return i
		}
	}
	return -1
}

func RuneIndexR(runes []rune, str string) int {
	for i := len(runes) - 1; i >= 0; i-- {
		if string(runes[i]) == str {
			return i
		}
	}
	return -1
}

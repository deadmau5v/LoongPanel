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

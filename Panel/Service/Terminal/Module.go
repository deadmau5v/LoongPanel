/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：定义了Terminal模块的数据结构
 */

package Terminal

import (
	"bytes"
	"os"
	"sync"
	"time"
)

type Screen struct {
	Name        string    `json:"name"`
	Id          uint32    `json:"id"`
	Time        time.Time `json:"time"`
	Tmx         *os.File
	Output      *bytes.Buffer
	subscribers []chan []byte
	outputLen   int
	Connected   bool
}

type ScreenManager struct {
	Screens map[uint32]*Screen
	Mu      sync.RWMutex
}

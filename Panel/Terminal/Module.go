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
}

type ScreenManager struct {
	Screens map[uint32]*Screen
	Mu      sync.RWMutex
}

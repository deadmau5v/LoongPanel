package Terminal

import "sync"

func init() {
	MainScreenManager = &ScreenManager{
		Screens: make(map[uint32]*Screen),
		Mu:      sync.RWMutex{},
	}
}

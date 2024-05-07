/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：初始化一个全局变量
 */

package Terminal

import "sync"

func init() {
	MainScreenManager = &ScreenManager{
		Screens: make(map[uint32]*Screen),
		Mu:      sync.RWMutex{},
	}
}

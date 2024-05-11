/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：Terminal模块的方法 操作Screen的输入输出等
 */

package Terminal

import (
	"LoongPanel/Panel/Service/System"
	"bytes"
	"errors"
	"io"
	"log"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

func (s *Screen) Subscribe() chan []byte {
	c := make(chan []byte, 1024)
	s.subscribers = append(s.subscribers, c)
	return c
}

func (s *Screen) Unsubscribe(c chan []byte) {
	for i, subscriber := range s.subscribers {
		if subscriber == c {
			s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
			close(c)
			break
		}
	}
}

func (s *Screen) Publish() {
	newOutput := s.Output.Bytes()[s.outputLen:]
	if string(newOutput) != "" {
		s.outputLen = len(s.Output.Bytes())
		for _, subscriber := range s.subscribers {
			select {
			case subscriber <- newOutput:
			default:
				// 如果订阅者的通道已满就忽略
			}
		}
	}
}

func (s *Screen) InputByte(b []byte) {
	_, _ = s.Tmx.Write(b)
	s.Publish()
}

func (sm *ScreenManager) Create(name string, id uint32) error {

	flag := false
	for _, v := range sm.Screens {
		if v.Id == id {
			flag = true
		}
	}
	if flag {
		return errors.New("已存在ID")
	}

	sm.Mu.Lock()
	defer sm.Mu.Unlock()
	var buf bytes.Buffer

	c := getInitShell()
	if c == nil {
		return errors.New("无法获取Shell")
	}
	tty, err := pty.Start(c)
	if err != nil {
		log.Fatal(err)
	}

	err = pty.Setsize(tty, &pty.Winsize{
		Rows: 16,
		Cols: 200,
	})
	screen := &Screen{
		Name:        name,
		Id:          id,
		Time:        time.Now(),
		Tmx:         tty,
		Output:      &buf,
		subscribers: make([]chan []byte, 0),
		outputLen:   0,
	}

	output := io.MultiWriter(&buf)

	go func() {
		if _, err := io.Copy(output, screen.Tmx); err != nil {
			log.Printf("Error reading from pty: %v", err)
		}
	}()

	sm.Screens[id] = screen
	return nil
}

func (sm *ScreenManager) GetScreen(id int) *Screen {
	for _, v := range sm.Screens {
		if v.Id == uint32(id) {
			return v
		}
	}
	return nil
}

func (s *Screen) GetOutput() string {
	return s.Output.String()
}

func (s *Screen) Input(cmd string) {
	_, _ = s.Tmx.Write([]byte(cmd + "\n"))
}

func (s *Screen) Close() {
	_ = s.Tmx.Close()
}

func (sm *ScreenManager) Close(id int) {
	sm.Mu.Unlock()
	screens := make(map[uint32]*Screen)
	for i, v := range sm.Screens {
		if v.Id != uint32(id) {
			screens[i] = v
		}
	}
	sm.Screens = screens
	sm.Mu.Unlock()
}

func getInitShell() *exec.Cmd {
	switch System.Data.OSName {
	case "linux":
		return exec.Command("bash")
	}
	return nil
}

func GetNextId() int {
	return int(time.Now().Unix())
}

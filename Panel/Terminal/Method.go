package Terminal

import (
	"LoongPanel/Panel/System"
	"bytes"
	"errors"
	"io"
	"log"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

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
	tty, err := pty.Start(c)
	if err != nil {
		log.Fatal(err)
	}

	err = pty.Setsize(tty, &pty.Winsize{
		Rows: 16,
		Cols: 200,
	})

	screen := &Screen{
		Name:   name,
		Id:     id,
		Time:   time.Now(),
		Tmx:    tty,
		Output: &buf,
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

func (s *Screen) InputByte(b []byte) {
	_, _ = s.Tmx.Write(b)
}

func (s *Screen) GetOutputChannel() chan string {
	c := make(chan string, 1024)
	idx := 0
	go func() {
		for {
			c <- s.Output.String()[idx:]
			idx = len(s.Output.String())
		}
	}()
	return c
}

func (sm *ScreenManager) Input(id uint32, cmd string) {
	sm.Mu.Lock()

	screen, ok := sm.Screens[id]
	if !ok {
		return
	}

	_, _ = screen.Tmx.Write([]byte(cmd + "\n"))
	sm.Mu.Unlock()
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
	case "windows":
		return exec.Command("cmd")
	}
	return nil
}

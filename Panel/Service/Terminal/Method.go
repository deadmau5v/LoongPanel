/*
 * 创建人： deadmau5v
 * 创建时间： 2024-5-7
 * 文件作用：Terminal 网页终端实现
 */

package Terminal

import (
	"LoongPanel/Panel/Service/PanelLog"
	"io"

	"github.com/gorilla/websocket"
	"github.com/helloyi/go-sshclient"
	"golang.org/x/crypto/ssh"
)

type WsReaderWriter struct {
	*websocket.Conn
}

func (w *WsReaderWriter) Write(p []byte) (n int, err error) {
	writer, err := w.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}
	defer writer.Close()
	return writer.Write(p)
}

func (w *WsReaderWriter) Read(p []byte) (n int, err error) {
	var msgType int
	var reader io.Reader
	for {
		msgType, reader, err = w.Conn.NextReader()
		if err != nil {
			return 0, err
		}
		if msgType != websocket.TextMessage {
			continue
		}
		return reader.Read(p)
	}
}

func Shell(c *websocket.Conn, host, port, user, password string) error {

	config := &sshclient.TerminalConfig{
		Term:   "xterm",
		Height: 40,
		Weight: 80,
		Modes: ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		},
	}
	rw := &WsReaderWriter{c}
	client, err := sshclient.DialWithPasswd(host+":"+port, user, password)
	if err != nil {
		return err
	}
	terminal := client.Terminal(config)
	terminal.SetStdio(rw, rw, rw)
	err = terminal.Start()
	if err != nil {
		PanelLog.DEBUG("[网页终端]", "链接关闭", err.Error())
		return err
	}

	defer client.Close()
	return nil
}

package main

import (
	"fmt"

	"github.com/helloyangqi/gonet"
	"github.com/helloyangqi/gonet/buffer"
)


type SessionHandler struct{}

func (s SessionHandler) OnConnected(sess *gonet.TCPSession) {
	sess.SetBufferSize(16, 4)
	fmt.Println("OnConnected:", sess.Conn.RemoteAddr().String())
}

func (s SessionHandler) OnDisconnected(sess *gonet.TCPSession, err error) {
	fmt.Println("OnDisconncted", sess.Conn.RemoteAddr().String())
}


func (sh *SessionHandler) OnMessage(sess *gonet.TCPSession, param interface{}) {
	buffer := param.(*buffer.Buffer)
	fmt.Printf("session[%s] peek:%s\n", sess.Conn.RemoteAddr().String(), string(buffer.Peek(20)))
	if buffer.Size() >= 20 {
		fmt.Printf("session[%s] OnMessage:%s\n", sess.Conn.RemoteAddr().String(), string(buffer.Read(20)))
	}
}

func (sh *SessionHandler) OnError(sess *gonet.TCPSession, err error) {
	fmt.Printf("session[%s] error:%v\n", sess.Conn.RemoteAddr().String(), err)
}

func main() {
	cfg := gonet.TCPServerConfig{
		ListenAddr:     "127.0.0.1:5566",
		//ServerHandler:  &ServerHandler{},
		SessionHandler: &SessionHandler{},
	}
	svr, err := gonet.NewTCPServer(cfg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ch := svr.Start()
	err = <-ch
	fmt.Println("server exit:", err.Error())
}

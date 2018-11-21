package gonet

import (
	"bufio"
	"fmt"
	"time"
)

type ServerHandler struct{}

func (s ServerHandler) OnConnected(sess *TCPSession) {
	fmt.Println("OnConnected:", sess.Conn.RemoteAddr().String())
}

func (s ServerHandler) OnDisconnected(sess *TCPSession, err error) {
	fmt.Println("OnDisconncted", sess.Conn.RemoteAddr().String())
}

type SessionHandler struct{}

func (sh *SessionHandler) OnMessage(sess *TCPSession, param interface{}) error {
	fmt.Println("OnMessage")
	reader := param.(*bufio.Reader)
	data, err := reader.Discard(10)
	fmt.Printf("session[%s] message:%s, err:%v\n", sess.Conn.RemoteAddr().String(), string(data), err)
	return nil
}

func (sh *SessionHandler) OnError(sess *TCPSession, err error) {
	fmt.Printf("session[%s] error:%v\n", sess.Conn.RemoteAddr().String(), err)
}

func main() {
	cfg := TCPServerConfig{
		ListenAddr:     "127.0.0.1:5566",
		ServerHandler:  &ServerHandler{},
		SessionHandler: &SessionHandler{},
	}
	svr, err := NewTCPServer(cfg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ch := svr.Start()
	go func() {
		time.Sleep(5 * time.Second)
		//svr.Stop()
	}()
	err = <-ch
	fmt.Println("server exit:", err.Error())
}

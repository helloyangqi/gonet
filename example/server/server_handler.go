package main

import (
	"fmt"

	"github.com/helloyangqi/gonet"
	"github.com/helloyangqi/gonet/buffer"
)

type ServerHandler struct{}

func (s ServerHandler) OnConnected(sess *gonet.TCPSession) {
	//sess.SetBufferSize(16, 4)
	fmt.Println("OnConnected:", sess.Conn.RemoteAddr().String())
}

func (s ServerHandler) OnDisconnected(sess *gonet.TCPSession, err error) {
	fmt.Println("OnDisconncted", sess.Conn.RemoteAddr().String())
}

type SessionHandler struct{}

func (sh *SessionHandler) OnMessage(sess *gonet.TCPSession, param interface{}) {
	buffer := param.(*buffer.Buffer)
	fmt.Printf("session[%s] peek:%s\n", sess.Conn.RemoteAddr().String(), string(buffer.Peek(20)))
	//if buffer.Size() >= 20 {
	fmt.Printf("session[%s] OnMessage:%s\n", sess.Conn.RemoteAddr().String(), string(buffer.Read(20)))
	//}
}

func (sh *SessionHandler) OnError(sess *gonet.TCPSession, err error) {
	fmt.Printf("session[%s] error:%v\n", sess.Conn.RemoteAddr().String(), err)
}

package gonet

import (
	"bufio"
	"fmt"
	"github.com/helloyangqi/gonet/decoder"
	"github.com/helloyangqi/gonet/encoder"
	"time"
)

type SessionHandler struct{}

func (s *SessionHandler) OnConnected(sess *TCPSession) {
	fmt.Println("OnConnected:", sess.Conn.RemoteAddr().String())
	sess.decoderList.PushBack(decoder.NeweLengthDecoder(true, true))
	sess.encoderList.PushBack(encoder.NewLengthEncoder(true, true))
}

func (s *SessionHandler) OnDisconnected(sess *TCPSession, err error) {
	fmt.Println("OnDisconncted", sess.Conn.RemoteAddr().String())
}

func (sh *SessionHandler) OnMessage(sess *TCPSession, param interface{})  {
	fmt.Println("OnMessage")
	reader := param.(*bufio.Reader)
	data, err := reader.Discard(10)
	fmt.Printf("session[%s] message:%s, err:%v\n", sess.Conn.RemoteAddr().String(), string(data), err)
}

func (sh *SessionHandler) OnError(sess *TCPSession, err error) {
	fmt.Printf("session[%s] error:%v\n", sess.Conn.RemoteAddr().String(), err)
}

func main() {
	cfg := TCPServerConfig{
		ListenAddr:     "127.0.0.1:5566",
		//ServerHandler:  &ServerHandler{},
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

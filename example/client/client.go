package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/helloyangqi/gonet"
	"github.com/helloyangqi/gonet/buffer"
)

type ClientHandler struct{}

func (h *ClientHandler) OnMessage(sess *gonet.TCPSession, msg interface{}) {
	buff := msg.(buffer.Buffer)
	fmt.Printf("client[%s] recv message:%s\n", sess.Conn.LocalAddr().String(), string(buff.ReadAll()))
}

func (h *ClientHandler) OnError(sess *gonet.TCPSession, err error) {
	fmt.Printf("client[%s] error:%v\n", sess.Conn.LocalAddr().String(), err)
}

func (h *ClientHandler) OnDisconnect(sess *gonet.TCPSession, err error) {
	fmt.Printf("client[%s] disconnect, error:%v\n", sess.Conn.LocalAddr().String(), err)
}

func main() {
	count, interval := 1, 1000
	flag.IntVar(&count, "count", 1, "user count")
	flag.IntVar(&interval, "interval", 1000, "send message interval. ms")
	flag.Parse()

	config := gonet.TCPClientConfig{
		Host:        "localhost:5566",
		Handler:     &ClientHandler{},
		DialTimeout: time.Second * 5,
	}

	clients := make([]*gonet.TCPClient, 0, count)
	for i := 0; i < count; i++ {
		func() {
			client := gonet.NewTCPClient(config)
			if err := client.Start(); err != nil {
				fmt.Printf("client start failed:%v", err)
				return
			}
			clients = append(clients, client)
		}()
	}

	ticker := time.NewTicker(time.Millisecond * time.Duration(interval))
	for {
		select {
		case <-ticker.C:
			data := strings.Repeat("helloworld", 1)
			for _, client := range clients {
				client.Write(data)
			}
		}
	}
	return
}

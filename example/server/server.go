package main

import (
	"fmt"

	"github.com/helloyangqi/gonet"
)

func main() {
	cfg := gonet.TCPServerConfig{
		ListenAddr:     "127.0.0.1:5566",
		ServerHandler:  &ServerHandler{},
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

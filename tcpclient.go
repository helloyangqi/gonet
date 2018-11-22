package gonet

import (
	"net"
	"time"
)

type TCPClientHandler interface {
	OnMessage(*TCPSession, interface{})
	OnError(*TCPSession, error)
	OnDisconnect(*TCPSession, error)
}

type TCPClientConfig struct {
	Host        string
	DialTimeout time.Duration
	Handler     TCPClientHandler
}

type TCPClient struct {
	*TCPSession
	Host        string
	handler     TCPClientHandler
	dialTimeout time.Duration
}

func NewTCPClient(config TCPClientConfig) *TCPClient {
	client :=  &TCPClient{Host: config.Host, handler: config.Handler, dialTimeout: config.DialTimeout}
	client.TCPSession = newTCPSession(nil, client.handler.OnDisconnect, client.handler)
	return client
}

func (c *TCPClient) Start() error {
	conn, err := net.DialTimeout("tcp", c.Host, c.dialTimeout)
	if err != nil {
		return err
	}
	c.TCPSession.Conn = conn.(*net.TCPConn)
	c.start()
	return nil
}

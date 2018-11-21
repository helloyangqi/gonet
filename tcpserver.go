package gonet

import (
	"net"
	"sync"
	"time"
)

type TCPServerHandler interface {
	OnConnected(*TCPSession)
	OnDisconnected(*TCPSession, error)
}

type TCPServerConfig struct {
	ListenAddr     string
	ServerHandler  TCPServerHandler
	SessionHandler TCPSessionHandler
}

type TCPServer struct {
	ListenAddr  *net.TCPAddr
	listener    *net.TCPListener
	closeCh     chan interface{}
	exitCh      chan error
	sessions    sync.Map
	svrHandler  TCPServerHandler
	sessHandler TCPSessionHandler
}

func NewTCPServer(config TCPServerConfig) (*TCPServer, error) {
	addr, err := net.ResolveTCPAddr("tcp", config.ListenAddr)
	if err != nil {
		return nil, err
	}
	s := &TCPServer{
		ListenAddr:  addr,
		closeCh:     make(chan interface{}, 1),
		exitCh:      make(chan error, 1),
		svrHandler:  config.ServerHandler,
		sessHandler: config.SessionHandler,
	}
	return s, nil
}

func (s *TCPServer) Stop() {
	s.listener.Close()
	s.sessions.Range(func(k interface{}, v interface{}) bool {
		sess := v.(*TCPSession)
		sess.Close()
		return true
	})
}

func (s *TCPServer) Start() (ch <-chan error) {
	ch = s.exitCh
	var err error
	s.listener, err = net.ListenTCP("tcp", s.ListenAddr)
	if err != nil {
		s.exitCh <- err
		return
	}

	go func() {
		var tempDelay time.Duration
		for {
			conn, err := s.listener.AcceptTCP()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					nlog.Error("TCPServer[%s] Accept error:%v; retrying in %v", s.ListenAddr.String(), err, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				s.exitCh <- err
				return
			} else {
				tempDelay = 0
				sess := newTCPSession(conn, s.onDisconnected, s.sessHandler)
				s.onConnected(sess)
				sess.start()
			}
		}
	}()

	return
}

func (s *TCPServer) onConnected(sess *TCPSession) {
	s.sessions.Store(sess.Conn.RemoteAddr().String(), sess)
	s.svrHandler.OnConnected(sess)
}

func (s *TCPServer) onDisconnected(sess *TCPSession, err error) {
	s.sessions.Delete(sess.Conn.RemoteAddr().String())
	s.svrHandler.OnDisconnected(sess, err)
}

func (s *TCPServer) Broadcast(msg interface{}) {
	s.sessions.Range(func(k, v interface{}) bool {
		sess, ok := v.(*TCPSession)
		if ok {
			sess.Write(msg)
		}
		return true
	})
}

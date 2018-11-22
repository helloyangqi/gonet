package gonet

import (
	"container/list"
	"fmt"
	"net"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/helloyangqi/gonet/buffer"

	"github.com/helloyangqi/gonet/encoder"

	"github.com/helloyangqi/gonet/decoder"
)

type TCPSessionHandler interface {
	OnMessage(*TCPSession, interface{})
	OnError(*TCPSession, error)
}

type TCPSessionCloseFunc func(*TCPSession, error)

type SessionState int32

const (
	Connected               SessionState = 0x01
	Connecting              SessionState = 0x02
	Closed                  SessionState = 0xff
	defaultWriteChannelSize int          = 4096
)

type TCPSession struct {
	Conn         *net.TCPConn
	closeFn      TCPSessionCloseFunc
	handler      TCPSessionHandler
	state        SessionState
	keepAlive    time.Duration
	decoderList  *list.List
	encoderList  *list.List
	writeChannel chan interface{}
	buffer       *buffer.Buffer
	closeChannel chan interface{}
	Context      interface{}
}

func newTCPSession(conn *net.TCPConn, closef TCPSessionCloseFunc, h TCPSessionHandler) *TCPSession {
	return &TCPSession{
		Conn:         conn,
		closeFn:      closef,
		handler:      h,
		state:        Connected,
		decoderList:  list.New(),
		encoderList:  list.New(),
		writeChannel: make(chan interface{}, defaultWriteChannelSize),
		buffer:       buffer.NewBuffer(),
		closeChannel: make(chan interface{}, 1)}
}

func (sess *TCPSession) start() {
	atomic.StoreInt32((*int32)(&sess.state), int32(Connected))
	go sess.reader()
	go sess.writer()
}

func (sess *TCPSession) reader() {
	var err error
	defer func() {
		nlog.Debug("TCPSession[%s] reader exit, error:%v", sess.Conn.RemoteAddr().String(), err)
		sess.Conn.CloseRead()
		sess.setState(Closed)
		sess.closeFn(sess, err)
		sess.closeChannel <- 1
	}()

	for atomic.LoadInt32((*int32)(&sess.state)) == int32(Connected) {
		if sess.keepAlive.Nanoseconds() > 0 {
			if err = sess.Conn.SetReadDeadline(time.Now().Add(sess.keepAlive)); err != nil {
				sess.handler.OnError(sess, err)
				sess.closeFn(sess, err)
				return
			}
		}

		n, err := sess.buffer.ReadFrom(sess.Conn)
		if err != nil {
			sess.handler.OnError(sess, err)
			sess.closeFn(sess, err)
			return
		}

		if n <= 0 {
			sess.closeFn(sess, nil)
			return
		}

		var out interface{} = nil
		if sess.decoderList.Len() == 0 {
			out = sess.buffer
		} else {
			var in interface{} = sess.buffer
			for i := sess.decoderList.Front(); i != nil && in != nil; i = i.Next() {
				dc := i.Value.(decoder.Decoder)
				if out, err = dc.Decode(in); err != nil {
					sess.handler.OnError(sess, err)
					break
				} else {
					if out == nil {
						break
					}
				}
				in = out
			}
		}
		if out != nil {
			sess.handler.OnMessage(sess, out)
		}
	}
}

func (sess *TCPSession) writer() {
	var err error
	defer func() {
		nlog.Debug("TCPSession[%s] writer exit, error:%v", sess.Conn.RemoteAddr().String(), err)
		sess.Conn.CloseWrite()
		sess.setState(Closed)
		//sess.closeFn(sess, err)
	}()

	for atomic.LoadInt32((*int32)(&sess.state)) == int32(Connected) {
		select {
		case <-sess.closeChannel:
			return
		case msg, ok := <-sess.writeChannel:
			{
				if !ok {
					return
				}

				var in interface{} = msg
				var out interface{} = in
				for i := sess.encoderList.Front(); i != nil && in != nil; i = i.Next() {
					ec := i.Value.(encoder.Encoder)
					if out, err = ec.Encode(in); err != nil {
						sess.handler.OnError(sess, err)
						break
					}
					in = out
				}

				if out != nil {
					var data []byte
					switch v := out.(type) {
					case []byte:
						data = v
					case string:
						data = []byte(v)
					default:
						sess.handler.OnError(sess, fmt.Errorf("TCPSession encoder return unknow type:%s", reflect.TypeOf(out).String()))
						continue
					}

					if err = sess.send(data); err != nil {
						sess.handler.OnError(sess, err)
						return
					}
				}
			}
		}
	}
}

func (sess *TCPSession) send(data []byte) error {
	if !sess.IsConnected() {
		return fmt.Errorf("TCPSession[%s] already disconnected on send, state[%d]", sess.Conn.RemoteAddr().String(), sess.state)
	}

	w := 0
	bl := len(data)
	for w < bl {
		wlen, err := sess.Conn.Write(data[w:])
		if err != nil {
			return err
		}
		w += wlen
	}
	return nil
}

func (sess *TCPSession) Close() {
	if atomic.LoadInt32((*int32)(&sess.state)) != int32(Closed) {
		atomic.StoreInt32((*int32)(&sess.state), int32(Closed))
		sess.Conn.Close()
	}
}

func (sess *TCPSession) Write(msg interface{}) error {
	select {
	case sess.writeChannel <- msg:
		return nil
	default:
		return fmt.Errorf("TCPSession write channel is full")
	}
}

func (sess *TCPSession) SetWriteChannelSize(size int) {
	close(sess.writeChannel)
	sess.writeChannel = make(chan interface{}, size)
}

func (sess *TCPSession) SetBufferSize(size int, lowestCap int) {
	sess.buffer = nil
	sess.buffer = buffer.NewBufferSize(size, lowestCap)
}

func (sess *TCPSession) setState(state SessionState) {
	atomic.StoreInt32((*int32)(&sess.state), int32(state))
}

func (sess *TCPSession) IsConnected() bool {
	return atomic.LoadInt32((*int32)(&sess.state)) == int32(Connected)
}

func (sess *TCPSession) AddEncoder(ec encoder.Encoder) {
	sess.encoderList.PushBack(ec)
}

func (sess *TCPSession) AddDecoder(dc decoder.Decoder) {
	sess.decoderList.PushBack(dc)
}

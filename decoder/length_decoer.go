package decoder

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/helloyangqi/gonet/buffer"
)

const LengthFieldLength int = 4

type LengthDecoder struct {
	IsLittleEndian  bool
	IsContainLength bool
}

func NeweLengthDecoder(isLittleEndian bool, isContainLength bool) *LengthDecoder {
	return &LengthDecoder{isLittleEndian, isContainLength}
}

func (dc LengthDecoder) Decode(in interface{}) (interface{}, error) {
	var buf *buffer.Buffer
	switch v := in.(type) {
	case *buffer.Buffer:
		buf = v
	case []byte:
		buf = buffer.NewBufferSize(len(v), len(v))
		buf.Write(v)
	case string:
		buf = buffer.NewBufferSize(len(v), len(v))
		buf.Write([]byte(v))
	default:
		return nil, fmt.Errorf("LengthDecoder decode has unknow type:%s", reflect.TypeOf(in).String())
	}

	if buf.Size() < LengthFieldLength {
		return nil, nil
	}

	var frameLength uint32 = 0
	if dc.IsLittleEndian {
		frameLength = binary.LittleEndian.Uint32(buf.Peek(LengthFieldLength))
	} else {
		frameLength = binary.BigEndian.Uint32(buf.Peek(LengthFieldLength))
	}

	if !dc.IsContainLength {
		frameLength += uint32(LengthFieldLength)
	}
	if uint32(buf.Size()) < frameLength {
		return nil, nil
	}
	buf.Discard(LengthFieldLength)
	messageLength := frameLength
	if dc.IsContainLength {
		messageLength -= uint32(LengthFieldLength)
	}
	return buf.Read(int(messageLength)), nil
}

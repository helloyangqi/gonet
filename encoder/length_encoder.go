package encoder

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/helloyangqi/gonet/buffer"
)

const LengthFieldLength int = 4

type LengthEncoder struct {
	IsLittleEndian  bool
	IsContainLength bool
}

func (dc LengthEncoder) Encode(in interface{}) (interface{}, error) {
	var data []byte
	switch v := in.(type) {
	case *buffer.Buffer:
		data = v.ReadAll()
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return nil, fmt.Errorf("LengthEncoder encode has unknow type:%s", reflect.TypeOf(in).String())
	}

	var frameLength uint32 = uint32(len(data))
	if dc.IsContainLength {
		frameLength += uint32(LengthFieldLength)
	}

	buffer := buffer.NewBuffer()
	lengthData := make([]byte, LengthFieldLength)
	if dc.IsLittleEndian {
		binary.LittleEndian.PutUint32(lengthData, frameLength)
	} else {
		binary.BigEndian.PutUint32(lengthData, frameLength)
	}

	buffer.Write(lengthData)
	buffer.Write(data)
	return buffer, nil
}

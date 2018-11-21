package decoder

import (
	"fmt"
	"reflect"

	"github.com/helloyangqi/gonet/buffer"

	"github.com/golang/protobuf/proto"
)

type ProtobufDecoder struct {
	msgType reflect.Type
}

func NewProtobufDecoder(msg proto.Message) *ProtobufDecoder {
	rv := reflect.Indirect(reflect.ValueOf(msg))
	return &ProtobufDecoder{msgType: rv.Type()}
}

func (dc ProtobufDecoder) Decode(in interface{}) (interface{}, error) {
	var data []byte
	switch v := in.(type) {
	case *buffer.Buffer:
		data = v.ReadAll()
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return nil, fmt.Errorf("ProtobufDecoder decode has unknow type:%s", reflect.TypeOf(in).String())
	}

	pb := reflect.New(dc.msgType).Interface().(proto.Message)
	if err := proto.Unmarshal(data, pb); err != nil {
		return nil, fmt.Errorf("ProtobufDecoder decode unmarshal error:%s", err.Error())
	}
	return pb, nil
}

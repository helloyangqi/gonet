package encoder

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
)

type ProtobufEncoder struct {
}

func (ec ProtobufEncoder) Encode(in interface{}) (interface{}, error) {
	pb, ok := in.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("ProtobufEncoder Encode has error type:%s", reflect.TypeOf(in).String())
	}
	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, fmt.Errorf("ProtobufEncoder Encode marshal error:%v", err)
	}
	return data, nil
}

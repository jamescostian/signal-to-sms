// Package protobuf provides a frameio.Encoder that encodes to protobufs
package protobuf

import (
	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type protobufEncoder struct{}

var ProtobufEncoder frameio.Encoder = protobufEncoder{}

func (protobufEncoder) Marshal(m proto.Message) ([]byte, error) {
	return proto.Marshal(m)
}
func (protobufEncoder) Unmarshal(b []byte, m protoreflect.ProtoMessage) error {
	return proto.Unmarshal(b, m)
}

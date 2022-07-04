// Package prototext provides a frameio.Encoder that encodes to prototext, and does so deterministically if you use the deterministic_output build tag
package prototext

import (
	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type prototextEncoder struct{}

var PrototextEncoder frameio.Encoder = prototextEncoder{}

func (prototextEncoder) Marshal(m proto.Message) ([]byte, error) {
	return []byte(prototext.Format(m)), nil
}
func (prototextEncoder) Unmarshal(b []byte, m protoreflect.ProtoMessage) error {
	return prototext.Unmarshal(b, m)
}

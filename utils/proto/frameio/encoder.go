package frameio

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

//go:generate mockery --name Encoder --inpackage

type Encoder interface {
	Unmarshal(b []byte, m protoreflect.ProtoMessage) error
	Marshal(m proto.Message) ([]byte, error)
}

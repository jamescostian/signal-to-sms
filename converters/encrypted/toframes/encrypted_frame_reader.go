package toframes

import (
	"io"

	"github.com/jamescostian/signal-to-sms/internal/fromsignalback/decrypt"
	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/protobuf"
)

// NewFrameReader provides a FrameReader that reads from an io.Reader that has an encrypted Signal for Android backup
func NewFrameReader(reader io.Reader, password string) (frameio.FrameReader, error) {
	return decrypt.NewBackupFile(reader, password, protobuf.ProtobufEncoder)
}

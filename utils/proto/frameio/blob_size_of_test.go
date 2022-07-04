package frameio_test

import (
	"testing"

	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestBlobSizeOf(t *testing.T) {
	assert.Equal(t, uint32(3), frameio.BlobSizeOf(&signal.BackupFrame{Attachment: &signal.Attachment{Length: proto.Uint32(3)}}))
	assert.Equal(t, uint32(5), frameio.BlobSizeOf(&signal.BackupFrame{Avatar: &signal.Avatar{Length: proto.Uint32(5)}}))
	assert.Equal(t, uint32(7), frameio.BlobSizeOf(&signal.BackupFrame{Sticker: &signal.Sticker{Length: proto.Uint32(7)}}))
	assert.Equal(t, uint32(0), frameio.BlobSizeOf(&signal.BackupFrame{Header: &signal.Header{}}))
}

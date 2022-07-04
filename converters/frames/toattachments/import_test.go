package toattachments_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jamescostian/signal-to-sms/converters/frames/toattachments"
	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestNewFrameWriter(t *testing.T) {
	store := &attachments.MockStore{}
	writer := toattachments.NewFrameWriter(context.Background(), store)

	// Ignores SQL
	assert.NoError(t, writer.Write(&signal.BackupFrame{Statement: &signal.SqlStatement{
		Statement: proto.String("SELECT 1"),
	}}, nil))

	// Ignores stickers
	assert.NoError(t, writer.Write(&signal.BackupFrame{Sticker: &signal.Sticker{
		RowId:  proto.Uint64(123),
		Length: proto.Uint32(4),
	}}, []byte("test")))

	// Ignores avatars
	assert.NoError(t, writer.Write(&signal.BackupFrame{Avatar: &signal.Avatar{
		Name:        proto.String("1"),
		RecipientId: proto.String("2"),
		Length:      proto.Uint32(3),
	}}, []byte("hi!")))

	// Pays attention to attachments
	attachmentData := []byte("finally")
	store.On("SetAttachment", context.Background(), int64(123), attachmentData).Once().Return(nil)
	assert.NoError(t, writer.Write(&signal.BackupFrame{Attachment: &signal.Attachment{
		RowId:        proto.Uint64(123),
		AttachmentId: proto.Uint64(456),
		Length:       proto.Uint32(uint32(len(attachmentData))),
	}}, attachmentData))

	// Bubbles up errors
	err := fmt.Errorf("oh no")
	store.On("SetAttachment", context.Background(), int64(888), attachmentData).Once().Return(err)
	assert.Equal(t, err, writer.Write(&signal.BackupFrame{Attachment: &signal.Attachment{
		RowId:        proto.Uint64(888),
		AttachmentId: proto.Uint64(456),
		Length:       proto.Uint32(uint32(len(attachmentData))),
	}}, attachmentData))
}

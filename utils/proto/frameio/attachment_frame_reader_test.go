package frameio_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestAttachmentFrameReader(t *testing.T) {
	t.Run("eofs_immediately_if_no_attachments", func(t *testing.T) {
		store := &attachments.MockStore{}
		store.On("ListAttachments", context.Background()).Return([]int64{}, nil)
		reader, err := frameio.NewAttachmentFrameReader(context.Background(), store)
		if assert.NoError(t, err) {
			frame := &signal.BackupFrame{}
			blob := new([]byte)
			assert.Equal(t, io.EOF, reader.Read(frame, blob))
		}
	})
	t.Run("can_read_a_few_attachments_and_reach_eof", func(t *testing.T) {
		store := &attachments.MockStore{}
		blobs := [][]byte{{'f', 'o', 'o'}, {'h', 'e', 'l', 'l', 'o'}}
		ids := []int64{42, 256}
		store.On("ListAttachments", context.Background()).Return(ids, nil)
		store.On("GetAttachment", context.Background(), ids[0]).Once().Return(blobs[0], nil)
		store.On("GetAttachment", context.Background(), ids[1]).Once().Return(blobs[1], nil)
		reader, err := frameio.NewAttachmentFrameReader(context.Background(), store)
		if assert.NoError(t, err) {
			frame := &signal.BackupFrame{}
			blob := new([]byte)

			assert.NoError(t, reader.Read(frame, blob))
			assert.Equal(t, blobs[0], *blob)
			assert.True(t, proto.Equal(frame, &signal.BackupFrame{
				Attachment: &signal.Attachment{
					RowId:        proto.Uint64(uint64(ids[0])),
					AttachmentId: proto.Uint64(uint64(ids[0])),
					Length:       proto.Uint32(uint32(len(blobs[0]))),
				},
			}))

			assert.NoError(t, reader.Read(frame, blob))
			assert.Equal(t, blobs[1], *blob)
			assert.True(t, proto.Equal(frame, &signal.BackupFrame{
				Attachment: &signal.Attachment{
					RowId:        proto.Uint64(uint64(ids[1])),
					AttachmentId: proto.Uint64(uint64(ids[1])),
					Length:       proto.Uint32(uint32(len(blobs[1]))),
				},
			}))

			assert.Equal(t, io.EOF, reader.Read(frame, blob))
		}
	})
	t.Run("forwards_get_attachment_errors", func(t *testing.T) {
		store := &attachments.MockStore{}
		store.On("ListAttachments", context.Background()).Return([]int64{1}, nil)
		store.On("GetAttachment", context.Background(), int64(1)).Once().Return(nil, errors.New("error"))
		reader, err := frameio.NewAttachmentFrameReader(context.Background(), store)
		if assert.NoError(t, err) {
			frame := &signal.BackupFrame{}
			blob := new([]byte)
			assert.Error(t, reader.Read(frame, blob))
		}
	})
}

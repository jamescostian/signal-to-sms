package frameio_test

import (
	"context"
	"io"
	"testing"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func TestMergedReader(t *testing.T) {
	setup := func(frameToReturn *signal.BackupFrame, errReading error, attachments attachments.Store) (frame *signal.BackupFrame, blob *[]byte, err error) {
		frame = &signal.BackupFrame{}
		blob = new([]byte)
		reader := &frameio.MockFrameReadCloser{}
		reader.On("Read", frame, blob).Once().Run(func(args mock.Arguments) {
			f := args.Get(0).(*signal.BackupFrame)
			proto.Reset(f)
			proto.Merge(f, frameToReturn)
		}).Return(errReading)
		err = frameio.NewMergedReader(context.Background(), reader, attachments).Read(frame, blob)
		return
	}
	t.Run("success_regular_sql_statement", func(t *testing.T) {
		myFrame := &signal.BackupFrame{
			Statement: &signal.SqlStatement{
				Statement:  proto.String("SELECT ?"),
				Parameters: []*signal.SqlStatement_SqlParameter{{IntegerParameter: new(uint64)}},
			},
		}
		frame, blob, err := setup(myFrame, nil, &attachments.MockStore{})
		assert.NoError(t, err)
		assert.True(t, proto.Equal(myFrame, frame))
		assert.Len(t, *blob, 0)
	})
	t.Run("success_attachment", func(t *testing.T) {
		myFrame := &signal.BackupFrame{
			Attachment: &signal.Attachment{
				RowId:        proto.Uint64(123),
				AttachmentId: proto.Uint64(456),
				Length:       proto.Uint32(1),
			},
		}
		attachments := &attachments.MockStore{}
		attachments.On("GetAttachment", context.Background(), int64(123)).Once().Return([]byte{'!'}, nil)
		frame, blob, err := setup(myFrame, nil, attachments)
		assert.NoError(t, err)
		assert.True(t, proto.Equal(myFrame, frame))
		assert.Equal(t, *blob, []byte{'!'})
	})
	t.Run("success_avatar", func(t *testing.T) {
		// Test out a avatar - this is needed because avatars have a length, so you need a blob of the right size for a avatar
		myFrame := &signal.BackupFrame{
			Avatar: &signal.Avatar{
				Name:        proto.String("unknown"),
				RecipientId: proto.String("1"),
				Length:      proto.Uint32(10),
			},
		}
		frame, blob, err := setup(myFrame, nil, &attachments.MockStore{})
		assert.NoError(t, err)
		assert.True(t, proto.Equal(myFrame, frame))
		assert.Len(t, *blob, 10)
	})
	t.Run("success_sticker", func(t *testing.T) {
		// Test out a sticker - this is needed because stickers have a length, so you need a blob of the right size for a sticker
		myFrame := &signal.BackupFrame{
			Sticker: &signal.Sticker{
				RowId:  proto.Uint64(0),
				Length: proto.Uint32(10),
			},
		}
		frame, blob, err := setup(myFrame, nil, &attachments.MockStore{})
		assert.NoError(t, err)
		assert.True(t, proto.Equal(myFrame, frame))
		assert.Len(t, *blob, 10)
	})
	t.Run("passes_eof_through", func(t *testing.T) {
		frame := &signal.BackupFrame{}
		blob := new([]byte)
		reader := &frameio.MockFrameReadCloser{}
		reader.On("Read", frame, blob).Once().Return(io.EOF)
		merged := frameio.NewMergedReader(context.Background(), reader, &attachments.MockStore{})
		assert.Equal(t, io.EOF, merged.Read(frame, blob))
	})
}

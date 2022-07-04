package decrypt_test

import (
	"bytes"
	"testing"

	"github.com/jamescostian/signal-to-sms/internal/fromsignalback/decrypt"
	"github.com/jamescostian/signal-to-sms/internal/testdata"
	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

var (
	salt = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

	// First 4 bytes are the length of the header, second 2 bytes are the header itself (made up in this case, they will be decoded by a proto encoder)
	initialHeader       = []byte{0x0, 0x0, 0x0, 0x2, 0x12, 0x34}
	frameCleartext      = []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}
	frame               = []byte{0x0, 0x0, 0x0, 0x1a, 0xca, 0x8e, 0xa2, 0x15, 0xa6, 0xdf, 0x7c, 0x88, 0xe6, 0x12, 0xfe, 0xf2, 0xd2, 0xf0, 0x39, 0x62, 0x3e, 0x3, 0x4d, 0xb7, 0x57, 0x56, 0xf3, 0x16, 0x43, 0x87}
	attachmentCleartext = []byte{0x31, 0x41, 0x59, 0x26}
	attachment          = []byte{0x0, 0x0, 0x0, 0xe, 0xc3, 0xd7, 0x72, 0x43, 0x3d, 0xd6, 0x3e, 0xc0, 0x23, 0xa3, 0x3b, 0x44, 0x14, 0x77}
)

// The IV slice passed in will have its contents modified (because normally that's 100% safe). Use this function to always end up with a fresh, never-seen-before iv slice
func iv() []byte {
	return []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
}

func TestNewBackupFile(t *testing.T) {
	// panic(fmt.Sprintf("%#v\n", frame))
	t.Run("empty_file_fails", func(t *testing.T) {
		_, err := decrypt.NewBackupFile(bytes.NewBuffer(nil), testdata.TestDataPassword, &frameio.MockEncoder{})
		assert.Error(t, err)
	})
	t.Run("too_large_header_length_causes_fail", func(t *testing.T) {
		_, err := decrypt.NewBackupFile(bytes.NewBuffer(initialHeader[:len(initialHeader)-1]), testdata.TestDataPassword, &frameio.MockEncoder{})
		assert.Error(t, err)
	})
	t.Run("error_unmarshaling_proto", func(t *testing.T) {
		enc := &frameio.MockEncoder{}
		enc.On("Unmarshal", initialHeader[4:], mock.Anything).Once().Return(errors.New("idk"))
		_, err := decrypt.NewBackupFile(bytes.NewBuffer(initialHeader), testdata.TestDataPassword, enc)
		assert.Error(t, err)
	})
	t.Run("bad_iv", func(t *testing.T) {
		enc := &frameio.MockEncoder{}
		enc.On("Unmarshal", initialHeader[4:], mock.Anything).Once().Run(func(args mock.Arguments) {
			args.Get(1).(*signal.BackupFrame).Header = &signal.Header{Iv: []byte{0, 1, 2, 3}, Salt: salt}
		}).Return(nil)
		_, err := decrypt.NewBackupFile(bytes.NewBuffer(initialHeader), testdata.TestDataPassword, enc)
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		enc := &frameio.MockEncoder{}
		enc.On("Unmarshal", initialHeader[4:], mock.Anything).Once().Run(func(args mock.Arguments) {
			args.Get(1).(*signal.BackupFrame).Header = &signal.Header{Iv: iv(), Salt: salt}
		}).Return(nil)
		_, err := decrypt.NewBackupFile(bytes.NewBuffer(initialHeader), testdata.TestDataPassword, enc)
		require.NoError(t, err)
	})
}

func TestRead(t *testing.T) {
	t.Run("missing_frame_length", func(t *testing.T) {
		enc := &frameio.MockEncoder{}
		enc.On("Unmarshal", initialHeader[4:], mock.Anything).Once().Run(func(args mock.Arguments) {
			args.Get(1).(*signal.BackupFrame).Header = &signal.Header{Iv: iv(), Salt: salt}
		}).Return(nil)
		bf, err := decrypt.NewBackupFile(bytes.NewBuffer(initialHeader), testdata.TestDataPassword, enc)
		require.NoError(t, err)
		assert.Error(t, bf.Read(&signal.BackupFrame{}, &[]byte{}))
	})
	t.Run("missing_frame", func(t *testing.T) {
		enc := &frameio.MockEncoder{}
		enc.On("Unmarshal", initialHeader[4:], mock.Anything).Once().Run(func(args mock.Arguments) {
			args.Get(1).(*signal.BackupFrame).Header = &signal.Header{Iv: iv(), Salt: salt}
		}).Return(nil)
		bf, err := decrypt.NewBackupFile(bytes.NewBuffer(append(initialHeader, frame[:4]...)), testdata.TestDataPassword, enc)
		require.NoError(t, err)
		assert.Error(t, bf.Read(&signal.BackupFrame{}, &[]byte{}))
	})
	t.Run("mac_mismatch", func(t *testing.T) {
		enc := &frameio.MockEncoder{}
		enc.On("Unmarshal", initialHeader[4:], mock.Anything).Once().Run(func(args mock.Arguments) {
			args.Get(1).(*signal.BackupFrame).Header = &signal.Header{Iv: iv(), Salt: salt}
		}).Return(nil)
		data := append(initialHeader, frame...)
		data[len(data)-1]++ // Mess up the mac by 1
		bf, err := decrypt.NewBackupFile(bytes.NewBuffer(data), testdata.TestDataPassword, enc)
		require.NoError(t, err)
		assert.Error(t, bf.Read(&signal.BackupFrame{}, &[]byte{}))
	})
	t.Run("error_unmarshaling", func(t *testing.T) {
		enc := &frameio.MockEncoder{}
		enc.On("Unmarshal", initialHeader[4:], mock.Anything).Once().Run(func(args mock.Arguments) {
			args.Get(1).(*signal.BackupFrame).Header = &signal.Header{Iv: iv(), Salt: salt}
		}).Return(nil)
		enc.On("Unmarshal", frameCleartext, mock.Anything).Once().Return(errors.New("idk"))
		bf, err := decrypt.NewBackupFile(bytes.NewBuffer(append(initialHeader, frame...)), testdata.TestDataPassword, enc)
		require.NoError(t, err)
		assert.Error(t, bf.Read(&signal.BackupFrame{}, &[]byte{}))
	})
	t.Run("success_with_last_frame", func(t *testing.T) {
		enc := &frameio.MockEncoder{}
		enc.On("Unmarshal", initialHeader[4:], mock.Anything).Once().Run(func(args mock.Arguments) {
			args.Get(1).(*signal.BackupFrame).Header = &signal.Header{Iv: iv(), Salt: salt}
		}).Return(nil)
		enc.On("Unmarshal", frameCleartext, mock.Anything).Once().Run(func(args mock.Arguments) {
			args.Get(1).(*signal.BackupFrame).End = proto.Bool(true)
		}).Return(nil)
		bf, err := decrypt.NewBackupFile(bytes.NewBuffer(append(initialHeader, frame...)), testdata.TestDataPassword, enc)
		require.NoError(t, err)
		frame := &signal.BackupFrame{}
		assert.NoError(t, bf.Read(frame, &[]byte{}))
		assert.True(t, frame.GetEnd())
		assert.Nil(t, frame.Attachment)
		assert.Nil(t, frame.Avatar)
		assert.Nil(t, frame.Header)
		assert.Nil(t, frame.KeyValue)
		assert.Nil(t, frame.Preference)
		assert.Nil(t, frame.Statement)
		assert.Nil(t, frame.Sticker)
		assert.Nil(t, frame.Version)
	})
	t.Run("success_decrypting_blob", func(t *testing.T) {
		enc := &frameio.MockEncoder{}
		enc.On("Unmarshal", initialHeader[4:], mock.Anything).Once().Run(func(args mock.Arguments) {
			args.Get(1).(*signal.BackupFrame).Header = &signal.Header{Iv: iv(), Salt: salt}
		}).Return(nil)
		enc.On("Unmarshal", frameCleartext, mock.Anything).Once().Run(func(args mock.Arguments) {
			rowID := uint64(42)
			attachmentID := uint64(123)
			length := uint32(len(attachmentCleartext))
			args.Get(1).(*signal.BackupFrame).Attachment = &signal.Attachment{RowId: &rowID, AttachmentId: &attachmentID, Length: &length}
		}).Return(nil)
		bf, err := decrypt.NewBackupFile(bytes.NewBuffer(append(initialHeader, append(frame, attachment...)...)), testdata.TestDataPassword, enc)
		require.NoError(t, err)
		assert.NoError(t, bf.Read(nil, &[]byte{}))
	})
}

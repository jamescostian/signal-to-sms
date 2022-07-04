package frameio_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestBackupFramesWriteCloser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		buf := &closableBuffer{}
		encoder := &frameio.MockEncoder{}
		f1 := &signal.BackupFrame{Header: &signal.Header{Iv: []byte("hi")}}
		f2 := &signal.BackupFrame{End: proto.Bool(true)}
		encoder.On("Marshal", mock.Anything).Once().Return([]byte("output!"), nil)
		w := frameio.NewBackupFramesWriteCloser(buf, encoder)
		require.NoError(t, w.Write(f1, nil))
		require.NoError(t, w.Write(f2, nil))
		assert.NoError(t, w.Close())
		framesWritten := encoder.Calls[0].Arguments[0].(*signal.BackupFrames).Frames
		assert.True(t, proto.Equal(framesWritten[0], f1))
		assert.True(t, proto.Equal(framesWritten[1], f2))
		assert.Equal(t, []byte("output!"), buf.Bytes())
	})
	t.Run("error_marshalling", func(t *testing.T) {
		encoder := &frameio.MockEncoder{}
		encoder.On("Marshal", mock.Anything).Once().Return([]byte(nil), errors.New("error"))
		w := frameio.NewBackupFramesWriteCloser(&closableBuffer{}, encoder)
		require.NoError(t, w.Write(&signal.BackupFrame{}, nil))
		assert.Error(t, w.Close())
	})
	t.Run("error_writing_output", func(t *testing.T) {
		encoder := &frameio.MockEncoder{}
		encoder.On("Marshal", mock.Anything).Once().Return([]byte("output"), nil)
		w := frameio.NewBackupFramesWriteCloser(brokenWriter{}, encoder)
		require.NoError(t, w.Write(&signal.BackupFrame{}, nil))
		assert.Error(t, w.Close())
	})
}

func TestBackupFramesReadCloser(t *testing.T) {
	f1 := &signal.BackupFrame{Statement: &signal.SqlStatement{}}
	f2 := &signal.BackupFrame{Attachment: &signal.Attachment{}}
	t.Run("success", func(t *testing.T) {
		sourceData := []byte("hello")
		buf := &closableBuffer{*bytes.NewBuffer(sourceData)}
		encoder := &frameio.MockEncoder{}
		encoder.On("Unmarshal", sourceData, mock.Anything).Once().Run(func(args mock.Arguments) {
			frames := args.Get(1).(*signal.BackupFrames)
			*frames = signal.BackupFrames{Frames: []*signal.BackupFrame{f1, f2}}
		}).Return(nil)
		r, err := frameio.NewBackupFramesReadCloser(buf, encoder)
		require.NoError(t, err)
		frameRead := &signal.BackupFrame{}
		assert.NoError(t, r.Read(frameRead, nil))
		assert.True(t, proto.Equal(frameRead, f1))
		assert.NoError(t, r.Read(frameRead, nil))
		assert.True(t, proto.Equal(frameRead, f2))
		assert.Equal(t, io.EOF, r.Read(frameRead, nil))
		assert.NoError(t, r.Close())
	})
	t.Run("error_reading_input", func(t *testing.T) {
		_, err := frameio.NewBackupFramesReadCloser(brokenReader{}, &frameio.MockEncoder{})
		assert.Error(t, err)
	})
}

type closableBuffer struct {
	bytes.Buffer
}

func (c *closableBuffer) Close() error {
	return nil
}

type brokenWriter struct{}

func (brokenWriter) Write([]byte) (int, error) {
	return 0, errors.New("error")
}
func (brokenWriter) Close() error {
	return errors.New("error")
}

type brokenReader struct{}

func (brokenReader) Read([]byte) (int, error) {
	return 0, errors.New("error")
}
func (brokenReader) Close() error {
	return errors.New("error")
}

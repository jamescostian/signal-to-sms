package frameio_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestCopyFrames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dest, src := &frameio.MockFrameWriteCloser{}, &frameio.MockFrameReadCloser{}
		src.On("Read", mock.Anything, mock.Anything).Once().Run(func(args mock.Arguments) {}).Return(nil)
		dest.On("Write", mock.Anything, mock.Anything).Once().Run(func(args mock.Arguments) {}).Return(nil)
		src.On("Read", mock.Anything, mock.Anything).Once().Run(func(args mock.Arguments) {}).Return(nil)
		dest.On("Write", mock.Anything, mock.Anything).Once().Run(func(args mock.Arguments) {}).Return(nil)
		src.On("Read", mock.Anything, mock.Anything).Once().Return(io.EOF)
		assert.NoError(t, frameio.CopyFrames(context.Background(), dest, src))
	})
	t.Run("success_with_zero_frames_is_fine", func(t *testing.T) {
		dest, src := &frameio.MockFrameWriteCloser{}, &frameio.MockFrameReadCloser{}
		src.On("Read", mock.Anything, mock.Anything).Once().Return(io.EOF)
		assert.NoError(t, frameio.CopyFrames(context.Background(), dest, src))
	})
	t.Run("cannot_copy_if_cannot_read", func(t *testing.T) {
		dest, src := &frameio.MockFrameWriteCloser{}, &frameio.MockFrameReadCloser{}
		src.On("Read", mock.Anything, mock.Anything).Once().Return(errors.New("error"))
		assert.Error(t, frameio.CopyFrames(context.Background(), dest, src))
	})
	t.Run("cannot_copy_if_cannot_write", func(t *testing.T) {
		dest, src := &frameio.MockFrameWriteCloser{}, &frameio.MockFrameReadCloser{}
		src.On("Read", mock.Anything, mock.Anything).Once().Return(nil)
		dest.On("Write", mock.Anything, mock.Anything).Once().Return(errors.New("error"))
		assert.Error(t, frameio.CopyFrames(context.Background(), dest, src))
	})
}

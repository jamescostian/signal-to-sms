package frameio_test

import (
	"errors"
	"io"
	"testing"

	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func TestMultiWriter(t *testing.T) {
	f1 := &signal.BackupFrame{Statement: &signal.SqlStatement{}}
	f2 := &signal.BackupFrame{Attachment: &signal.Attachment{}}
	blob1, blob2 := []byte{}, []byte{}
	t.Run("success_0_writers", func(t *testing.T) {
		mw := frameio.NewMultiWriter()
		assert.NoError(t, mw.Write(f1, blob1))
		assert.NoError(t, mw.Write(f2, blob2))
	})
	t.Run("success_1_writer", func(t *testing.T) {
		w := &frameio.MockFrameWriter{}
		w.On("Write", f1, blob1).Once().Return(nil)
		w.On("Write", f2, blob2).Once().Return(nil)
		mw := frameio.NewMultiWriter(w)
		assert.NoError(t, mw.Write(f1, blob1))
		assert.NoError(t, mw.Write(f2, blob2))
	})

	t.Run("success_3_writers", func(t *testing.T) {
		w1 := &frameio.MockFrameWriter{}
		w2 := &frameio.MockFrameWriter{}
		w3 := &frameio.MockFrameWriter{}
		w1.On("Write", f1, blob1).Once().Return(nil)
		w1.On("Write", f2, blob2).Once().Return(nil)
		w2.On("Write", f1, blob1).Once().Return(nil)
		w2.On("Write", f2, blob2).Once().Return(nil)
		w3.On("Write", f1, blob1).Once().Return(nil)
		w3.On("Write", f2, blob2).Once().Return(nil)
		mw := frameio.NewMultiWriter(w1, w2, w3)
		assert.NoError(t, mw.Write(f1, blob1))
		assert.NoError(t, mw.Write(f2, blob2))
	})
	t.Run("error_from_one_of_3_writers", func(t *testing.T) {
		w1 := &frameio.MockFrameWriter{}
		w2 := &frameio.MockFrameWriter{}
		w3 := &frameio.MockFrameWriter{}
		w1.On("Write", f1, blob1).Once().Return(nil)
		w1.On("Write", f2, blob2).Once().Return(nil)
		w2.On("Write", f1, blob1).Once().Return(nil)
		w2.On("Write", f2, blob2).Once().Return(errors.New("error"))
		w3.On("Write", f1, blob1).Once().Return(nil)
		mw := frameio.NewMultiWriter(w1, w2, w3)
		assert.NoError(t, mw.Write(f1, blob1))
		assert.Error(t, mw.Write(f2, blob2))
	})
}

func TestMultiReader(t *testing.T) {
	f1 := &signal.BackupFrame{Statement: &signal.SqlStatement{}}
	f2 := &signal.BackupFrame{Attachment: &signal.Attachment{}}
	f3 := &signal.BackupFrame{Sticker: &signal.Sticker{}}
	f4 := &signal.BackupFrame{Avatar: &signal.Avatar{}}
	blob1, blob2, blob3 := &[]byte{}, &[]byte{}, &[]byte{}
	t.Run("success_1_reader_with_end_frame_to_skip", func(t *testing.T) {
		endFrame := &signal.BackupFrame{}
		r := &frameio.MockFrameReader{}
		mr := frameio.NewMultiReader(r)

		// First, a normal read
		r.On("Read", f1, blob1).Once().Return(nil)
		assert.NoError(t, mr.Read(f1, blob1))

		// Next, reading returns an end frame.
		// This causes a skip to the next frame (because multi reader ignores end frames, it only emits one itself when all readers EOF).
		// The next frame doesn't exist though, the reader gives an EOF.
		// Because it's the EOF of the last reader, it's the end of everything, so an end frame is provided.
		r.On("Read", endFrame, blob2).Once().Run(func(args mock.Arguments) {
			frame := args.Get(0).(*signal.BackupFrame)
			frame.Reset()
			frame.End = proto.Bool(true)
		}).Return(nil)
		r.On("Read", endFrame, blob2).Once().Return(io.EOF)
		assert.NoError(t, mr.Read(endFrame, blob2))
		assert.True(t, *endFrame.End)

		// A final read will now provide an EOF, because all of the readers have finished all of their reading.
		assert.Equal(t, io.EOF, mr.Read(endFrame, blob2))
	})
	t.Run("success_3_readers", func(t *testing.T) {
		endFrame := &signal.BackupFrame{}
		r1 := &frameio.MockFrameReader{}
		r2 := &frameio.MockFrameReader{}
		r3 := &frameio.MockFrameReader{}
		mr := frameio.NewMultiReader(r1, r2, r3)

		// Read succeeds from first reader
		r1.On("Read", f1, blob1).Once().Return(nil)
		assert.NoError(t, mr.Read(f1, blob1))

		// First reader, it just so happens, is done. So the next read must go on to the next reader
		r1.On("Read", f2, blob2).Once().Return(io.EOF)
		r2.On("Read", f2, blob2).Once().Return(nil)
		assert.NoError(t, mr.Read(f2, blob2))

		// Another successful read
		r2.On("Read", f3, blob3).Once().Return(nil)
		assert.NoError(t, mr.Read(f3, blob3))

		// 2 read failures in a row (both readers are done).
		// Now that all readers are done, an end frame can be provided
		r2.On("Read", endFrame, blob3).Once().Return(io.EOF)
		r3.On("Read", endFrame, blob3).Once().Return(io.EOF)
		assert.NoError(t, mr.Read(endFrame, blob3))
		assert.True(t, *endFrame.End)

		// Finally, EOF is provided after the end frame
		assert.Equal(t, io.EOF, mr.Read(endFrame, blob3))
	})

	t.Run("end_immediately_if_0_readers", func(t *testing.T) {
		endFrame := &signal.BackupFrame{}
		mr := frameio.NewMultiReader()
		assert.NoError(t, mr.Read(endFrame, blob1))
		assert.True(t, *endFrame.End)
		assert.Equal(t, io.EOF, mr.Read(endFrame, blob1))
	})
	t.Run("error_from_one_of_3_readers", func(t *testing.T) {
		r1 := &frameio.MockFrameReader{}
		r2 := &frameio.MockFrameReader{}
		r3 := &frameio.MockFrameReader{}
		mr := frameio.NewMultiReader(r1, r2, r3)
		r1.On("Read", f1, blob1).Once().Return(nil)
		assert.NoError(t, mr.Read(f1, blob1))
		r1.On("Read", f2, blob2).Once().Return(io.EOF)
		r2.On("Read", f2, blob2).Once().Return(nil)
		assert.NoError(t, mr.Read(f2, blob2))
		r2.On("Read", f3, blob3).Once().Return(nil)
		assert.NoError(t, mr.Read(f3, blob3))
		r2.On("Read", f4, blob3).Once().Return(errors.New("error"))
		assert.Error(t, mr.Read(f4, blob3))
	})
}

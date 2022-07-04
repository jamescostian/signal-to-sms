package toencrypted

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	encryptedData "github.com/jamescostian/signal-to-sms/internal/testdata/encrypted"
	prototextData "github.com/jamescostian/signal-to-sms/internal/testdata/prototext"
	"github.com/jamescostian/signal-to-sms/internal/testutil"
	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/prototext"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func init() {
	fillRandBytes = fillRandBytesDeterministically
}

func TestEncrypt(t *testing.T) {
	t.Run("success_with_testdata", func(t *testing.T) {
		output, err := ioutil.TempFile("", "")
		require.NoError(t, err)
		defer os.Remove(output.Name())

		// Read in a real example of prototext + attachments for frames
		var frames frameio.FrameReader
		frames, err = frameio.NewBackupFramesReadCloser(prototextData.Example, prototext.PrototextEncoder)
		require.NoError(t, err)
		store := testutil.SetupSeededAttachmentStore(t)
		frames = frameio.NewMergedReader(context.Background(), frames, store)

		assert.NoError(t, Encrypt(context.Background(), frames, output, encryptedData.ExamplePassword))

		// Now compare the output
		outBytes, err := os.ReadFile(output.Name())
		require.NoError(t, err)
		expectedBytes, err := os.ReadFile(encryptedData.ExamplePath)
		require.NoError(t, err)
		// Not doing assert.Equal because the output is CRAZY LONG. These two are helpful though
		assert.Equal(t, len(expectedBytes), len(outBytes), "output should have the expected length")
		assert.True(t, bytes.Equal(expectedBytes, outBytes), "output should match")
	})
	t.Run("can_be_stopped_by_context", func(t *testing.T) {
		output, err := ioutil.TempFile("", "")
		require.NoError(t, err)
		defer os.Remove(output.Name())

		ctx, cancel := context.WithCancel(context.Background())
		frames := &frameio.MockFrameReader{}
		frames.On("Read", mock.Anything, mock.Anything).Once().Run(func(args mock.Arguments) {
			// Cancel, which should cause an error when reading the next frame
			cancel()
			// Set this frame properly anyway, so that it doesn't cause an error
			frame := args.Get(0).(*signal.BackupFrame)
			frame.Reset()
			frame.Statement = &signal.SqlStatement{
				Statement:  proto.String("SELECT ?"),
				Parameters: []*signal.SqlStatement_SqlParameter{{IntegerParameter: proto.Uint64(1)}},
			}
		}).Return(nil)

		assert.Equal(t, context.Canceled, Encrypt(ctx, frames, output, encryptedData.ExamplePassword))
	})
	t.Run("errors_if_file_is_already_closed", func(t *testing.T) {
		output, err := ioutil.TempFile("", "")
		require.NoError(t, err)
		require.NoError(t, output.Close())
		os.Remove(output.Name())

		frames := &frameio.MockFrameReader{}
		assert.Error(t, Encrypt(context.Background(), frames, output, encryptedData.ExamplePassword))
	})
	t.Run("errors_if_frame_reader_errors", func(t *testing.T) {
		output, err := ioutil.TempFile("", "")
		require.NoError(t, err)
		defer os.Remove(output.Name())

		frames := &frameio.MockFrameReader{}
		frames.On("Read", mock.Anything, mock.Anything).Once().Return(fmt.Errorf("ohno"))
		assert.Error(t, Encrypt(context.Background(), frames, output, encryptedData.ExamplePassword))
	})
}

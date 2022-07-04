package conversion

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jamescostian/signal-to-sms/formats"
	sqlite "github.com/jamescostian/signal-to-sms/internal/testdata/attachments/sqlite"
	"github.com/jamescostian/signal-to-sms/internal/testdata/prototext"
	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialExtAttachmentsSetup(t *testing.T) {
	t.Run("has_attachment_input", func(t *testing.T) {
		t.Run("uses_attachmentIn_if_there_is_no_attachmentOut_so_that_attachments_can_be_read_from_attachmentIn", func(t *testing.T) {
			attachIn := &attachments.MockStore{}
			c := &Conversion{
				conversionPath: []MsgConverter{
					{InputFormat: formats.PrototextMsg, OutputFormat: formats.BackupFrame},
					{InputFormat: formats.BackupFrame, OutputFormat: formats.Encrypted},
				},
			}
			require.NoError(t, c.initialExtAttachmentsSetup(context.Background(), attachIn, nil))
			assert.Equal(t, attachIn, c.ExternalAttachments)
		})
		t.Run("converts_from_attachmentIn_to_attachmentOut_if_attachments_are_always_external", func(t *testing.T) {
			attachIn, attachOut := &attachments.MockStore{}, &attachments.MockStore{}
			// Attachments will be converted. Here's the mock calls that will happen to perform that conversion:
			ids := []int64{1, 23, 456}
			data := [][]byte{{'a'}, {'b'}, {'c'}}
			attachIn.On("ListAttachments", context.Background()).Return(ids, nil)
			attachIn.On("GetAttachment", context.Background(), ids[0]).Return(data[0], nil)
			attachIn.On("GetAttachment", context.Background(), ids[1]).Return(data[1], nil)
			attachIn.On("GetAttachment", context.Background(), ids[2]).Return(data[2], nil)
			attachOut.On("SetAttachment", context.Background(), ids[0], data[0]).Return(nil)
			attachOut.On("SetAttachment", context.Background(), ids[1], data[1]).Return(nil)
			attachOut.On("SetAttachment", context.Background(), ids[2], data[2]).Return(nil)
			// Finally, the attachmentIn will be closed because the normal Convert process won't ever see attachmentIn, so it won't be able to close it.
			attachIn.On("Close").Return(nil)
			c := &Conversion{
				conversionPath: []MsgConverter{
					{InputFormat: formats.PrototextMsg, OutputFormat: formats.SQLiteMsg},
					{InputFormat: formats.SQLiteMsg, OutputFormat: formats.PrototextMsg},
				},
				attachmentOut: attachOut,
			}
			require.NoError(t, c.initialExtAttachmentsSetup(context.Background(), attachIn, attachOut))
			assert.Equal(t, attachOut, c.ExternalAttachments)
		})
	})
	t.Run("no_attachment_input", func(t *testing.T) {
		t.Run("uses_attachmentOut_if_converting_directly_to_it_rather_than_creating_a_temporary_attachment_for_no_reason", func(t *testing.T) {
			attachOut := &attachments.MockStore{}
			c := &Conversion{
				conversionPath: []MsgConverter{
					{InputFormat: formats.Encrypted, OutputFormat: formats.BackupFrame},
					{InputFormat: formats.BackupFrame, OutputFormat: formats.PrototextMsg},
				},
				attachmentOut: attachOut,
			}
			require.NoError(t, c.initialExtAttachmentsSetup(context.Background(), nil, attachOut))
			assert.Equal(t, attachOut, c.ExternalAttachments)
		})
		t.Run("creates_temporary_ExternalAttachments_if_needed_so_that_attachmentOut_is_only_written_to_at_the_end,_instead_of_having_to_be_cleaned_out_when_switching_between_external_and_included_attachments", func(t *testing.T) {
			attachOut := &attachments.MockStore{}
			c := &Conversion{
				conversionPath: []MsgConverter{
					{InputFormat: formats.Encrypted, OutputFormat: formats.BackupFrame},
					{InputFormat: formats.BackupFrame, OutputFormat: formats.PrototextMsg},
					{InputFormat: formats.PrototextMsg, OutputFormat: formats.BackupFrame},
					{InputFormat: formats.BackupFrame, OutputFormat: formats.PrototextMsg},
				},
				attachmentOut: attachOut}
			require.NoError(t, c.initialExtAttachmentsSetup(context.Background(), nil, attachOut))
			assert.NotEqual(t, attachOut, c.ExternalAttachments)
			assert.NoError(t, c.ExternalAttachments.Close())
		})
	})
}

func TestConvert(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		convPath := []MsgConverter{
			{
				InputFormat:  formats.PrototextMsg,
				OutputFormat: formats.BackupFrame,
				Convert: func(ctx context.Context, c *Conversion) (intermediateOutput interface{}, err error) {
					return
				},
			},
			{
				InputFormat:  formats.BackupFrame,
				OutputFormat: formats.Encrypted,
				Convert: func(ctx context.Context, c *Conversion) (intermediateOutput interface{}, err error) {
					return
				},
			},
		}
		attachmentsIn, err := formats.SQLiteAttachments.Open(sqlite.ExamplePath)
		require.NoError(t, err)
		defer attachmentsIn.Close()
		attachmentsOut, err := attachments.NewTempStore()
		require.NoError(t, err)
		defer attachmentsOut.Close()
		msgOut, err := ioutil.TempFile("", "")
		require.NoError(t, err)
		defer msgOut.Close()

		assert.NoError(t, Convert(context.Background(), convPath, prototext.ExamplePath, msgOut.Name(), attachmentsIn, attachmentsOut, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600))
	})
	t.Run("one_converter_error_halts_future_conversions", func(t *testing.T) {
		convPath := []MsgConverter{
			{
				InputFormat:  formats.PrototextMsg,
				OutputFormat: formats.BackupFrame,
				Convert: func(ctx context.Context, c *Conversion) (intermediateOutput interface{}, err error) {
					return nil, errors.New("oh no")
				},
			},
			{
				InputFormat:  formats.BackupFrame,
				OutputFormat: formats.Encrypted,
				Convert: func(ctx context.Context, c *Conversion) (intermediateOutput interface{}, err error) {
					require.Fail(t, "should not run the second conversion if the first failed")
					return
				},
			},
		}
		attachmentsIn, err := formats.SQLiteAttachments.Open(sqlite.ExamplePath)
		require.NoError(t, err)
		defer attachmentsIn.Close()
		attachmentsOut, err := attachments.NewTempStore()
		require.NoError(t, err)
		defer attachmentsOut.Close()
		msgOut, err := ioutil.TempFile("", "")
		require.NoError(t, err)
		defer msgOut.Close()

		assert.Error(t, Convert(context.Background(), convPath, prototext.ExamplePath, msgOut.Name(), attachmentsIn, attachmentsOut, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600))
	})
	t.Run("context_errors_bubble_up", func(t *testing.T) {
		convPath := []MsgConverter{
			{
				InputFormat:  formats.PrototextMsg,
				OutputFormat: formats.BackupFrame,
				Convert: func(ctx context.Context, c *Conversion) (intermediateOutput interface{}, err error) {
					return
				},
			},
		}
		attachmentsIn, err := formats.SQLiteAttachments.Open(sqlite.ExamplePath)
		require.NoError(t, err)
		defer attachmentsIn.Close()
		attachmentsOut, err := attachments.NewTempStore()
		require.NoError(t, err)
		defer attachmentsOut.Close()
		msgOut, err := ioutil.TempFile("", "")
		require.NoError(t, err)
		defer msgOut.Close()

		expiredCtx, expireCtx := context.WithCancel(context.Background())
		expireCtx()
		assert.Error(t, Convert(expiredCtx, convPath, prototext.ExamplePath, msgOut.Name(), attachmentsIn, attachmentsOut, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600))
	})
	t.Run("error_on_broken_attachments", func(t *testing.T) {
		convPath := []MsgConverter{
			{
				InputFormat:  formats.PrototextMsg,
				OutputFormat: formats.SQLiteMsg,
				Convert: func(ctx context.Context, c *Conversion) (intermediateOutput interface{}, err error) {
					return
				},
			},
		}
		msgOut, err := ioutil.TempFile("", "")
		require.NoError(t, err)
		defer msgOut.Close()

		assert.Error(t, Convert(context.Background(), convPath, prototext.ExamplePath, msgOut.Name(), attachments.NoReads(nil), attachments.NoWrites(nil), 0, 0))
	})
}

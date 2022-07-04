package conversion

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jamescostian/signal-to-sms/formats"
	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Writing tests for this is proving to take forever, so here's 1 large test that goes through virtually all of the bugs I can imagine

func TestComplexConversion(t *testing.T) {
	fmt.Println("Initial setup")
	inFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	defer os.Remove(inFile.Name())
	outFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	defer os.Remove(outFile.Name())
	attachmentOut := &attachments.MockStore{}
	attachmentOut.On("Close").Once().Return(nil)
	attachmentData := []byte("test")
	attachmentOut.On("SetAttachment", context.Background(), int64(3), attachmentData).Once().Return(nil)
	converters := []MsgConverter{
		{
			InputFormat:  formats.Encrypted,
			OutputFormat: formats.BackupFrame,
			Convert: func(ctx context.Context, c *Conversion) (interface{}, error) {
				_, ok := c.MsgIn.(*os.File)
				assert.True(t, ok, "MsgIn should be the result of opening the input file")
				assert.Nil(t, c.MsgOut, "there shouldn't be any output initially; calling Convert() the first time should result in the output being filled")
				assert.Error(t, c.ExternalAttachments.Close(), "you shouldn't be able to close ExternalAttachments; they weren't opened by you, and they may need to be used by other converters in this process")
				assert.Error(t, c.ExternalAttachments.SetAttachment(context.Background(), 0, []byte("test")), "you shouldn't be able to write to ExternalAttachments because the output format has attachments included in the messages (MsgOut)")
				_, err := c.ExternalAttachments.ListAttachments(context.Background())
				assert.Error(t, err, "you shouldn't be able to read from ExternalAttachments because the input format has attachments included in the messages (MsgIn)")
				return "output #1", nil
			},
		},
		{
			InputFormat:  formats.BackupFrame,
			OutputFormat: formats.PrototextMsg,
			Convert: func(ctx context.Context, c *Conversion) (interface{}, error) {
				assert.Equal(t, "output #1", c.MsgIn.(string), "the output of the non-writeable format that came prior should persist into this run as the input; otherwise, how could that output be used since it couldn't have been written?")
				_, ok := c.MsgOut.(*frameio.BackupFramesWriteCloser)
				assert.True(t, ok, "MsgOut should be the result of opening the input file")
				assert.Error(t, c.ExternalAttachments.Close(), "you shouldn't be able to close ExternalAttachments; they weren't opened by you, and they may need to be used by other converters in this process")
				assert.NoError(t, c.ExternalAttachments.SetAttachment(context.Background(), 0, []byte("test")), "you should be able to write to ExternalAttachments because the output format doesn't have attachments included in the messages (MsgOut)")
				_, err := c.ExternalAttachments.ListAttachments(context.Background())
				assert.Error(t, err, "you shouldn't be able to read from ExternalAttachments because the input format has attachments included in the messages (MsgIn)")
				return "output #2 (discarded, because prototextmsg should be opened for reading separately, since now it's opened for writing specifically)", nil
			},
		},
		{
			InputFormat:  formats.PrototextMsg,
			OutputFormat: formats.SQLiteMsg,
			Convert: func(ctx context.Context, c *Conversion) (interface{}, error) {
				_, ok := c.MsgIn.(*frameio.BackupFramesReadCloser)
				assert.True(t, ok, "MsgIn should be the result of opening the input file")
				assert.Error(t, c.ExternalAttachments.Close(), "you shouldn't be able to close ExternalAttachments; they weren't opened by you, and they may need to be used by other converters in this process")
				assert.NoError(t, c.ExternalAttachments.SetAttachment(context.Background(), 1, []byte("test")), "you should be able to write to ExternalAttachments because the output format doesn't have attachments included in the messages (MsgOut)")
				attachment0, err := c.ExternalAttachments.GetAttachment(context.Background(), 0)
				assert.Equal(t, []byte("test"), attachment0, "the ExternalAttachments from last run shouldn't disappear, since they weren't persisted anywhere - in the previous step, the output format didn't have attachments included in its messages, so if the ExternalAttachments were thrown away, then the attachments would just be gone in general")
				assert.NoError(t, err, "you should be able to read from ExternalAttachments because the input format doesn't have attachments included in the messages (MsgIn)")
				return "output #3 (discarded, because sqlitemsg should be opened for reading separately, since now it's opened for writing specifically)", nil
			},
		},
		{
			InputFormat:  formats.SQLiteMsg,
			OutputFormat: formats.SQLiteMsgAndAttachment,
			Convert: func(ctx context.Context, c *Conversion) (interface{}, error) {
				_, ok := c.MsgIn.(*sqlite.ClosableTx)
				assert.True(t, ok, "MsgIn should be the result of opening the input file")
				assert.Error(t, c.ExternalAttachments.Close(), "you shouldn't be able to close ExternalAttachments; they weren't opened by you, and they may need to be used by other converters in this process")
				assert.Error(t, c.ExternalAttachments.SetAttachment(context.Background(), 2, []byte("test")), "you shouldn't be able to write to ExternalAttachments because the output format has attachments included in the messages (MsgOut)")
				attachment1, err := c.ExternalAttachments.GetAttachment(context.Background(), 1)
				assert.Equal(t, []byte("test"), attachment1, "the ExternalAttachments from last run shouldn't disappear, since they weren't persisted anywhere - in the previous step, the output format didn't have attachments included in its messages, so if the ExternalAttachments were thrown away, then the attachments would just be gone in general")
				assert.NoError(t, err, "you should be able to read from ExternalAttachments because the input format doesn't have attachments included in the messages (MsgIn)")
				return "output #4 (discarded, because this OutputFormat has a OpenForWrites)", nil
			},
		},
		{
			InputFormat:  formats.SQLiteMsgAndAttachment,
			OutputFormat: formats.SQLiteMsg,
			Convert: func(ctx context.Context, c *Conversion) (interface{}, error) {
				_, ok := c.MsgIn.(*sqlite.ClosableTx)
				assert.True(t, ok, "MsgIn should be the result of opening the input file")
				assert.Error(t, c.ExternalAttachments.Close(), "you shouldn't be able to close ExternalAttachments; they weren't opened by you, and they may need to be used by other converters in this process")
				_, err := c.ExternalAttachments.GetAttachment(context.Background(), 0)
				assert.Error(t, err, "you shouldn't be able to read attachments from a previous run where attachments should have all been converted to being internal")
				assert.NoError(t, c.ExternalAttachments.SetAttachment(context.Background(), 2, []byte("test")), "you should be able to write to ExternalAttachments because the output format doesn't have attachments included in the messages (MsgOut)")
				_, err = c.ExternalAttachments.GetAttachment(context.Background(), 1)
				assert.Error(t, err, "the ExternalAttachments shouldn't be readable because the input format has attachments included within itself")
				return "output #5 (discarded, because this OutputFormat has a OpenForWrites)", nil
			},
		},
		{
			InputFormat:  formats.SQLiteMsg,
			OutputFormat: formats.XML,
			Convert: func(ctx context.Context, c *Conversion) (interface{}, error) {
				_, ok := c.MsgIn.(*sqlite.ClosableTx)
				assert.True(t, ok, "MsgIn should be the result of opening the input file")
				assert.Error(t, c.ExternalAttachments.Close(), "you shouldn't be able to close ExternalAttachments; they weren't opened by you, and they may need to be used by other converters in this process")
				assert.Error(t, c.ExternalAttachments.SetAttachment(context.Background(), 3, []byte("test")), "you shouldn't be able to write to ExternalAttachments because the output format has attachments included in the messages (MsgOut)")
				_, err := c.ExternalAttachments.GetAttachment(context.Background(), 1)
				assert.Error(t, err, "you shouldn't be able to read attachments from a previous run where attachments should have all been converted to being internal")
				attachment2, err := c.ExternalAttachments.GetAttachment(context.Background(), 2)
				assert.Equal(t, []byte("test"), attachment2, "the ExternalAttachments from last run shouldn't disappear, since they weren't persisted anywhere - in the previous step, the output format didn't have attachments included in its messages, so if the ExternalAttachments were thrown away, then the attachments would just be gone in general")
				assert.NoError(t, err, "you should be able to read from ExternalAttachments because the input format doesn't have attachments included in the messages (MsgIn)")
				return "output #6", nil
			},
		},
		{
			InputFormat:  formats.XML,
			OutputFormat: formats.SQLiteMsg,
			Convert: func(ctx context.Context, c *Conversion) (interface{}, error) {
				assert.Equal(t, "output #6", c.MsgIn.(string))
				assert.Error(t, c.ExternalAttachments.Close(), "you shouldn't be able to close ExternalAttachments; they weren't opened by you, and they may need to be used by other converters in this process")
				_, err := c.ExternalAttachments.GetAttachment(context.Background(), 2)
				assert.Error(t, err, "you shouldn't be able to read attachments from a previous run where attachments should have all been converted to being internal")
				assert.NoError(t, c.ExternalAttachments.SetAttachment(context.Background(), 3, attachmentData), "you should be able to write to ExternalAttachments because the output format doesn't have attachments included in the messages (MsgOut)")
				_, err = c.ExternalAttachments.GetAttachment(context.Background(), 3)
				assert.Error(t, err, "the ExternalAttachments shouldn't be readable because the input format has attachments included within itself")
				return "output #7 (discarded, because this OutputFormat has a OpenForWrites)", nil
			},
		},
	}
	c := &Conversion{
		conversionPath:     converters,
		msgInPath:          inFile.Name(),
		msgOutPath:         outFile.Name(),
		attachmentOut:      attachmentOut,
		MsgFileFlags:       os.O_CREATE | os.O_WRONLY,
		MsgFilePermissions: 0600,
	}
	require.NoError(t, c.initialExtAttachmentsSetup(context.Background(), nil, attachmentOut))

	fmt.Println("Trying out conversion step 1")
	require.False(t, c.Finished())
	require.False(t, c.IsFinalConversion())
	require.Equal(t, converters[0].InputFormat.Name, c.CurrentConverter().InputFormat.Name)
	require.Equal(t, converters[0].OutputFormat.Name, c.CurrentConverter().OutputFormat.Name)
	textSavedFromCleanUpFuncs := ""
	c.AddCleanUpFn(func(conversionSuccessful bool) error {
		textSavedFromCleanUpFuncs += fmt.Sprintln("First clean up function! Success:", conversionSuccessful)
		return nil
	})
	c.AddCleanUpFn(func(conversionSuccessful bool) error {
		textSavedFromCleanUpFuncs += fmt.Sprintln("Second clean up function! Success:", conversionSuccessful)
		return nil
	})
	c.AddFinalCleanUpFn(func(conversionSuccessful bool) error {
		textSavedFromCleanUpFuncs += fmt.Sprintln("First final clean up function! Success:", conversionSuccessful)
		return nil
	})
	c.AddFinalCleanUpFn(func(conversionSuccessful bool) error {
		textSavedFromCleanUpFuncs += fmt.Sprintln("Second final clean up function! Success:", conversionSuccessful)
		return nil
	})
	require.NoError(t, ConvertToNextFormat(context.Background(), c))
	assert.Equal(t, "First clean up function! Success: true\nSecond clean up function! Success: true\n", textSavedFromCleanUpFuncs)

	fmt.Println("Trying out conversion step 2")
	require.False(t, c.Finished())
	require.False(t, c.IsFinalConversion())
	require.Equal(t, converters[1].InputFormat.Name, c.CurrentConverter().InputFormat.Name)
	require.Equal(t, converters[1].OutputFormat.Name, c.CurrentConverter().OutputFormat.Name)
	require.NoError(t, ConvertToNextFormat(context.Background(), c))

	fmt.Println("Trying out conversion step 3")
	require.False(t, c.Finished())
	require.False(t, c.IsFinalConversion())
	require.Equal(t, converters[2].InputFormat.Name, c.CurrentConverter().InputFormat.Name)
	require.Equal(t, converters[2].OutputFormat.Name, c.CurrentConverter().OutputFormat.Name)
	require.NoError(t, ConvertToNextFormat(context.Background(), c))

	fmt.Println("Trying out conversion step 4")
	require.False(t, c.Finished())
	require.False(t, c.IsFinalConversion())
	require.Equal(t, converters[3].InputFormat.Name, c.CurrentConverter().InputFormat.Name)
	require.Equal(t, converters[3].OutputFormat.Name, c.CurrentConverter().OutputFormat.Name)
	require.NoError(t, ConvertToNextFormat(context.Background(), c))

	fmt.Println("Trying out conversion step 5")
	require.False(t, c.Finished())
	require.False(t, c.IsFinalConversion())
	require.Equal(t, converters[4].InputFormat.Name, c.CurrentConverter().InputFormat.Name)
	require.Equal(t, converters[4].OutputFormat.Name, c.CurrentConverter().OutputFormat.Name)
	require.NoError(t, ConvertToNextFormat(context.Background(), c))

	fmt.Println("Trying out conversion step 6")
	require.False(t, c.Finished())
	require.False(t, c.IsFinalConversion())
	require.Equal(t, converters[5].InputFormat.Name, c.CurrentConverter().InputFormat.Name)
	require.Equal(t, converters[5].OutputFormat.Name, c.CurrentConverter().OutputFormat.Name)
	require.NoError(t, ConvertToNextFormat(context.Background(), c))

	fmt.Println("Trying out conversion step 7")
	require.False(t, c.Finished())
	require.True(t, c.IsFinalConversion())
	require.Equal(t, converters[6].InputFormat.Name, c.CurrentConverter().InputFormat.Name)
	require.Equal(t, converters[6].OutputFormat.Name, c.CurrentConverter().OutputFormat.Name)
	require.NoError(t, ConvertToNextFormat(context.Background(), c))
	assert.True(t, c.Finished())
	assert.Equal(t, "First clean up function! Success: true\nSecond clean up function! Success: true\nFirst final clean up function! Success: true\nSecond final clean up function! Success: true\n", textSavedFromCleanUpFuncs)

	assert.Error(t, ConvertToNextFormat(context.Background(), c))
}

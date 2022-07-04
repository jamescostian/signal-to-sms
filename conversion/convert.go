package conversion

import (
	"context"
	"io/fs"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"go.uber.org/multierr"
)

type ctxKey string

var ConversionProgressCtxKey ctxKey = "ConversionProgress"

// Convert between MsgFormats and AttachmentFormats.
// Does not handle opening/closing input/output files, but will handle opening/closing any temporary files needed.
// If you want to get updates on the progress of the conversion, pass a chan *MsgConverter via the context key ConversionProgressCtxKey - it should be buffered, and of length len(conversionPath).
// If you don't know what to do for the last 2 parameters, I'd suggest using os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600
func Convert(
	ctx context.Context,
	conversionPath []MsgConverter,
	msgInputPath string,
	msgOutputPath string,
	attachmentIn attachments.Store,
	attachmentOut attachments.Store,
	msgFileFlags int,
	msgFilePermissions fs.FileMode,
) error {
	c := &Conversion{
		conversionPath:     conversionPath,
		attachmentOut:      attachmentOut,
		msgOutPath:         msgOutputPath,
		msgInPath:          msgInputPath,
		MsgFileFlags:       msgFileFlags,
		MsgFilePermissions: msgFilePermissions,
	}
	err := c.initialExtAttachmentsSetup(ctx, attachmentIn, attachmentOut)
	if err != nil {
		return err
	}

	progress, ok := ctx.Value(ConversionProgressCtxKey).(chan *MsgConverter)
	if !ok {
		// Buffered so no consumer is needed to read it. There will be at most exactly 1 message per MstConverter in the path.
		progress = make(chan *MsgConverter, len(conversionPath))
	}
	defer close(progress)
	for !c.Finished() && ctx.Err() == nil {
		converter := c.CurrentConverter()
		progress <- converter
		if err = c.convertToNextFormat(ctx); err != nil {
			break
		}
	}
	if ctx.Err() == nil && err == nil {
		return nil
	}

	c.unrecoverableError = true
	multierr.AppendInto(&err, c.cleanUp(false))
	multierr.AppendInto(&err, ctx.Err())
	return err
}

func (c *Conversion) initialExtAttachmentsSetup(ctx context.Context, attachmentIn attachments.Store, attachmentOut attachments.Store) (err error) {
	c.AddFinalCleanUpFn(func(conversionSuccessful bool) error {
		err := c.ExternalAttachments.Close()
		c.ExternalAttachments = nil
		return err
	})
	// In general, the conversion code doesn't really care about converting attachments, because different formats and converters can do wild things with attachments.
	// As a result, there's this one edge-case where if the whole conversion path involves never touching attachments, they will never actually get converted.
	// The following "if" statement takes care of that.
	// The other branches ensure that ExternalAttachments is set up properly in non-edge-cases
	if attachmentOut != nil && attachmentIn != nil && c.shouldUseFinalAttachmentOut() {
		if err = attachments.Copy(ctx, attachmentOut, attachmentIn); err != nil {
			return
		}
		c.ExternalAttachments = attachmentOut
		// At this point, attachmentIn will never be used again, and will be forgotten about. So it should be closed now, since it won't get closed later.
		err = attachmentIn.Close()
	} else if attachmentIn != nil {
		c.ExternalAttachments = attachmentIn
	} else if c.shouldUseFinalAttachmentOut() {
		c.ExternalAttachments = c.attachmentOut
	} else {
		// When this is closed, it will also delete any file that it creates
		c.ExternalAttachments, err = attachments.NewTempStore()
	}
	return
}

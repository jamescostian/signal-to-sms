package toattachments

import (
	"context"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
)

type importer struct {
	Store attachments.Store
	Err   error
	ctx   context.Context
}

// NewFrameWriter provide an efficient way to import ONLY attachments from Signal's backup protobufs into an attachment store (and ignore everything else)
func NewFrameWriter(ctx context.Context, store attachments.Store) frameio.FrameWriter {
	return &importer{ctx: ctx, Store: store}
}

func (i *importer) Write(frame *signal.BackupFrame, blob []byte) error {
	if frame.Attachment != nil {
		if i.Err = i.Store.SetAttachment(i.ctx, int64(*frame.Attachment.RowId), blob); i.Err != nil {
			return i.Err
		}
	}
	return nil
}

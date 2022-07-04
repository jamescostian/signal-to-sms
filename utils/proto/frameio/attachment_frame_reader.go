package frameio

import (
	"context"
	"io"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
)

func NewAttachmentFrameReader(ctx context.Context, attachments attachments.Store) (FrameReader, error) {
	ids, err := attachments.ListAttachments(ctx)
	return &attachmentFrameReader{ctx: ctx, Store: attachments, IDs: ids}, err
}

type attachmentFrameReader struct {
	Store attachments.Store
	IDs   []int64
	ctx   context.Context
}

func (attachments *attachmentFrameReader) Read(backupFrameDest *signal.BackupFrame, blobDest *[]byte) (err error) {
	if len(attachments.IDs) == 0 {
		return io.EOF
	}
	id := attachments.IDs[0]
	*blobDest, err = attachments.Store.GetAttachment(attachments.ctx, id)
	if err != nil {
		return
	}
	attachments.IDs = attachments.IDs[1:]

	u64ID := uint64(id)
	u32Len := uint32(len(*blobDest))
	backupFrameDest.Reset()
	backupFrameDest.Attachment = &signal.Attachment{RowId: &u64ID, AttachmentId: &u64ID, Length: &u32Len}
	return
}

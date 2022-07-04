package frameio

import (
	"bytes"
	"context"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
)

func NewMergedReader(ctx context.Context, reader FrameReader, attachments attachments.Store) FrameReader {
	return mergedReader{ctx, reader, attachments, bytes.NewBuffer(nil)}
}

type mergedReader struct {
	ctx         context.Context
	Reader      FrameReader
	Attachments attachments.Store
	emptyBytes  *bytes.Buffer
}

func (m mergedReader) Read(frame *signal.BackupFrame, blob *[]byte) (err error) {
	if err = m.Reader.Read(frame, blob); err != nil {
		return
	}
	*blob = nil
	if frame.Attachment != nil {
		*blob, err = m.Attachments.GetAttachment(m.ctx, int64(*frame.Attachment.RowId))
	}
	if frame.Avatar != nil {
		len := int(frame.Avatar.GetLength())
		m.emptyBytes.Grow(len)
		*blob = m.emptyBytes.Bytes()[:len]
	}
	if frame.Sticker != nil {
		len := int(frame.Sticker.GetLength())
		m.emptyBytes.Grow(len)
		*blob = m.emptyBytes.Bytes()[:len]
	}
	return
}

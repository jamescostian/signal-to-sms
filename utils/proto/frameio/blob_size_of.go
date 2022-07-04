package frameio

import "github.com/jamescostian/signal-to-sms/utils/proto/signal"

// BlobSizeOf tells you how big the blob for a given frame is
func BlobSizeOf(frame *signal.BackupFrame) uint32 {
	if frame.Attachment != nil {
		return frame.Attachment.GetLength()
	}
	if frame.Avatar != nil {
		return frame.Avatar.GetLength()
	}
	if frame.Sticker != nil {
		return frame.Sticker.GetLength()
	}
	return 0
}

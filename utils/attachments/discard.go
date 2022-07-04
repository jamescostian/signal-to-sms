package attachments

import (
	"context"
	"fmt"
)

type discard struct{}

// Discard never actually sets or deletes attachments. Never able to get attachments either.
func Discard() Store {
	return discard{}
}

// SetAttachment does nothing
func (discard) SetAttachment(ctx context.Context, id int64, blob []byte) error {
	return nil
}

// GetAttachment always fails because no attachments are stored
func (discard) GetAttachment(ctx context.Context, id int64) (blob []byte, err error) {
	return nil, fmt.Errorf("unable to retrieve attachment - this attachments.Store discards all attachments")
}

// ListAttachments always provides a nil slice of IDs
func (discard) ListAttachments(ctx context.Context) (ids []int64, err error) {
	return nil, nil
}

// DeleteAttachment does nothing
func (discard) DeleteAttachment(ctx context.Context, id int64) (err error) {
	return nil
}

// Close does nothing
func (discard) Close() (err error) {
	return nil
}

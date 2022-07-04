package attachments

import "context"

//go:generate mockery --name Store --inpackage

// Store is a way to store and retrieve binary blobs by IDs.
type Store interface {
	// SetAttachment saves an attachment, so that later on GetAttachment can retrieve it.
	SetAttachment(ctx context.Context, id int64, blob []byte) error
	// GetAttachment retrieves an attachment previously stored by SetAttachment. If the ID was never used with SetAttachment, or was deleted with DeleteAttachment, it returns an error.
	GetAttachment(ctx context.Context, id int64) (blob []byte, err error)
	// ListAttachments gives a list of attachment IDs that have been set and haven't been deleted.
	ListAttachments(ctx context.Context) (ids []int64, err error)
	// DeleteAttachment frees up resources used to store an attachment, and causes GetAttachment to error going forward when asked for the deleted attachment.
	DeleteAttachment(ctx context.Context, id int64) (err error)
	// Close allows for cleaning up
	Close() (err error)
}

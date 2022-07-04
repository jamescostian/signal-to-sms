package attachments

import (
	"context"

	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/pkg/errors"
)

// NewTempStore creates a Store in a temporary file, and when you Close the store, it deletes that file.
// It's a variable so it can be overwritten if you import this package.
// The built-in implementation uses SQLite, but any other implementation will work so long as its Close method removes all traces of this store
var NewTempStore = func() (Store, error) {
	db, delete, err := sqlite.CreateTemp()
	if err != nil {
		return nil, errors.Wrap(err, "unable to make temporary attachment store")
	}
	store, err := NewSQLStore(db)
	return &tempStore{
		store: store,
		close: func() error {
			delete()
			return nil
		},
	}, err
}

type tempStore struct {
	store Store
	close func() error
}

func (t *tempStore) SetAttachment(ctx context.Context, id int64, blob []byte) error {
	return t.store.SetAttachment(ctx, id, blob)
}
func (t *tempStore) GetAttachment(ctx context.Context, id int64) (blob []byte, err error) {
	return t.store.GetAttachment(ctx, id)
}
func (t *tempStore) ListAttachments(ctx context.Context) (ids []int64, err error) {
	return t.store.ListAttachments(ctx)
}
func (t *tempStore) DeleteAttachment(ctx context.Context, id int64) (err error) {
	return t.store.DeleteAttachment(ctx, id)
}
func (t *tempStore) Close() error {
	return t.close()
}

package attachments

import (
	"context"

	"github.com/jamescostian/signal-to-sms/utils/sqlite"
)

// sqlStore is a Store that stores attachments using database/sql.
// You can have it use transactions (*sql.Tx) or just a normal *sql.DB.
type sqlStore struct {
	// Queries should be run on this
	DBOrTx sqlite.DBOrTx
}

// NewSQLStore makes a Store that uses an SQL database to store attachments.
// You can supply either a *sql.DB or a *sql.Tx for the queries to run under.
func NewSQLStore(dbOrTx sqlite.DBOrTx) (Store, error) {
	// Incompatible with MSSQL 🤷‍♂️
	_, err := dbOrTx.Exec("CREATE TABLE IF NOT EXISTS _attachment_blobs (_id INTEGER PRIMARY KEY, data BLOB)")
	store := sqlStore{DBOrTx: dbOrTx}
	return &store, err
}

// SetAttachment saves an attachment to the DB
func (store *sqlStore) SetAttachment(ctx context.Context, id int64, blob []byte) error {
	_, err := store.DBOrTx.ExecContext(ctx, "INSERT INTO _attachment_blobs VALUES (?, ?)", id, blob)
	return err
}

// GetAttachment retrieves an attachment from the DB, or errors if no attachment could be found
func (store *sqlStore) GetAttachment(ctx context.Context, id int64) ([]byte, error) {
	var bytes []byte
	err := store.DBOrTx.QueryRowContext(ctx, "SELECT data FROM _attachment_blobs WHERE _id = ?", id).Scan(&bytes)
	return bytes, err
}

// ListAttachments gets a list of attachments in the DB
func (store *sqlStore) ListAttachments(ctx context.Context) ([]int64, error) {
	rows, err := store.DBOrTx.QueryContext(ctx, "SELECT _id FROM _attachment_blobs")
	if err != nil {
		return nil, err
	}
	var ids []int64
	var id int64
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// DeleteAttachment deletes an attachment from the DB
func (store *sqlStore) DeleteAttachment(ctx context.Context, id int64) error {
	_, err := store.DBOrTx.ExecContext(ctx, "DELETE FROM _attachment_blobs WHERE _id = ?", id)
	return err
}

func (store *sqlStore) Close() error {
	return sqlite.CommitOrClose(store.DBOrTx)
}

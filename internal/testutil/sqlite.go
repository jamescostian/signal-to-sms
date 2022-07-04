package testutil

import (
	"database/sql"
	"testing"

	sqliteAttachmentTestData "github.com/jamescostian/signal-to-sms/internal/testdata/attachments/sqlite"
	sqliteTestData "github.com/jamescostian/signal-to-sms/internal/testdata/sqlite"
	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/stretchr/testify/require"
)

func ExecSQL(t *testing.T, db *sql.DB, sql string, args ...interface{}) {
	_, err := db.Exec(sql, args...)
	require.NoError(t, err, "error running some SQL as part of the test")
}

func SetupDB(t *testing.T, basedOn string) *sql.DB {
	db, _, delete, err := sqlite.CreateTempCopyOf(basedOn)
	require.NoError(t, err)
	t.Cleanup(delete)
	return db
}
func SetupSeededMsgDB(t *testing.T) *sql.DB {
	return SetupDB(t, sqliteTestData.ExamplePath)
}
func SetupSeededAttachmentDB(t *testing.T) *sql.DB {
	return SetupDB(t, sqliteAttachmentTestData.ExamplePath)
}
func SetupSeededAttachmentStore(t *testing.T) attachments.Store {
	store, err := attachments.NewSQLStore(SetupDB(t, sqliteAttachmentTestData.ExamplePath))
	require.NoError(t, err)
	return store
}

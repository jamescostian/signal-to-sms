package sqlite_test

import (
	"testing"

	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTemp(t *testing.T) {
	db, delete, err := sqlite.CreateTemp()
	require.NoError(t, err, "Should be able to create a temporary SQLite DB, no problem")

	var one int
	err = db.QueryRow("SELECT 1").Scan(&one)
	assert.NoError(t, err, "Should be able to query from the temporary DB")

	delete()
	err = db.QueryRow("SELECT 1").Scan(&one)
	assert.Error(t, err, "Should not be able to query from a temporary DB after the file is deleted")
}

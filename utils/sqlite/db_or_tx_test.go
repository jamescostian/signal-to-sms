package sqlite

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommitOrClose(t *testing.T) {
	file := "deleteme-TestCommitOrClose-tx.sqlite"
	db, err := Open(file)
	defer os.Remove(file)
	require.NoError(t, err)
	tx, err := db.Begin()
	require.NoError(t, err)
	assert.NoError(t, CommitOrClose(tx), "Can commit a *sql.Tx")
	assert.NoError(t, CommitOrClose(db), "Can close a *sql.DB")
}

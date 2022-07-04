package sqlite

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenToTxNormalMethods(t *testing.T) {
	file := "deleteme-TestOpenToTxNormalMethods.sqlite"
	tx, err := OpenToTx(file)
	defer os.Remove(file)
	require.NoError(t, err)

	// All the normal methods should work fine
	_, err = tx.Exec("SELECT 1")
	assert.NoError(t, err)
	_, err = tx.ExecContext(context.Background(), "SELECT 1")
	assert.NoError(t, err)
	_, err = tx.Prepare("SELECT 1")
	assert.NoError(t, err)
	_, err = tx.PrepareContext(context.Background(), "SELECT 1")
	assert.NoError(t, err)
	_, err = tx.Query("SELECT 1")
	assert.NoError(t, err)
	_, err = tx.QueryContext(context.Background(), "SELECT 1")
	assert.NoError(t, err)

	var result interface{}
	err = tx.QueryRow("SELECT 1").Scan(&result)
	assert.NoError(t, err)
	err = tx.QueryRowContext(context.Background(), "SELECT 1").Scan(&result)
	assert.NoError(t, err)

	assert.NoError(t, tx.Close())
}

func TestOpenToTxClose(t *testing.T) {
	file := "deleteme-TestOpenToTxClose.sqlite"
	tx, err := OpenToTx(file)
	defer os.Remove(file)
	require.NoError(t, err)
	assert.NoError(t, tx.Close())
}

func TestOpenToTxCommit(t *testing.T) {
	file := "deleteme-TestOpenToTxCommit.sqlite"
	tx, err := OpenToTx(file)
	defer os.Remove(file)
	require.NoError(t, err)
	assert.NoError(t, tx.Commit())
}

func TestOpenToTxRollback(t *testing.T) {
	file := "deleteme-TestOpenToTxRollback.sqlite"
	tx, err := OpenToTx(file)
	defer os.Remove(file)
	require.NoError(t, err)
	assert.NoError(t, tx.Rollback())
}

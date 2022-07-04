package attachments_test

import (
	"context"
	"sort"
	"testing"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLStore(t *testing.T) {
	// First, set up: create a database and the sqlStore
	db, delete, err := sqlite.CreateTemp()
	require.NoError(t, err)
	defer delete()
	store, err := attachments.NewSQLStore(db)
	require.NoError(t, err)

	// Now add blobs to it
	err = store.SetAttachment(context.Background(), 123, []byte("hi"))
	require.NoError(t, err)
	err = store.SetAttachment(context.Background(), 42, []byte("\x00\x01\x02"))
	require.NoError(t, err)

	// Next, verify that the attachments are set
	blob, err := store.GetAttachment(context.Background(), 123)
	assert.NoError(t, err)
	assert.Equal(t, []byte("hi"), blob)
	blob2, err := store.GetAttachment(context.Background(), 42)
	assert.NoError(t, err)
	assert.Equal(t, []byte("\x00\x01\x02"), blob2)
	assert.NotEqual(t, blob2, blob, "Calling GetAttachment should return a different buffer every time, not reuse buffers and end up changing old values!")

	// The attachments should also be listed properly
	list, err := store.ListAttachments(context.Background())
	assert.NoError(t, err)
	sort.Slice(list, func(a int, b int) bool {
		return list[a] < list[b]
	})
	assert.Equal(t, []int64{42, 123}, list)

	// Verify that reading an attachment not set results in an error
	_, err = store.GetAttachment(context.Background(), 99)
	assert.Error(t, err)

	// Verify that deleting attachments works fine
	assert.NoError(t, store.DeleteAttachment(context.Background(), 42))
	_, err = store.GetAttachment(context.Background(), 42)
	assert.Error(t, err)
	list, err = store.ListAttachments(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, []int64{123}, list)

	assert.NoError(t, store.Close(), "Should be able to close, no problem")
}

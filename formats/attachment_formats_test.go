package formats_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jamescostian/signal-to-sms/formats"
	sqliteData "github.com/jamescostian/signal-to-sms/internal/testdata/sqlite"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscard(t *testing.T) {
	d, err := formats.DiscardAttachments.Open("???? invalid path, but it shouldn't matter")
	assert.NoError(t, err)

	assert.NoError(t, d.SetAttachment(context.Background(), 0, []byte("hi")))

	ids, err := d.ListAttachments(context.Background())
	assert.NoError(t, err, "Not having any attachments is fine")
	assert.Len(t, ids, 0, "There should not be any attachments listed")

	assert.NoError(t, d.DeleteAttachment(context.Background(), 0), "There's nothing stored, so it's already gone")

	_, err = d.GetAttachment(context.Background(), 0)
	assert.Error(t, err, "You can't get an attachment that hasn't been actually stored")

	assert.NoError(t, d.Close(), "Should not have any issues closing")

	assert.NoError(t, formats.DiscardAttachments.Initialize("", 0, 0))
}

func TestSQLite(t *testing.T) {
	db, path, delete, err := sqlite.CreateTempCopyOf(sqliteData.ExamplePath)
	assert.NoError(t, err)
	db.Close()
	defer delete()
	d, err := formats.SQLiteAttachments.Open(path)
	require.NoError(t, err)

	attachmentID := int64(823895)
	// In case this test got stopped earlier, ensure that there is no attachment from this test here
	if err := d.DeleteAttachment(context.Background(), attachmentID); err == nil {
		fmt.Println("Deleted attachment made for testing")
	}

	require.NoError(t, d.SetAttachment(context.Background(), attachmentID, []byte("hi")))

	ids, err := d.ListAttachments(context.Background())
	assert.NoError(t, err)
	assert.Contains(t, ids, attachmentID)

	_, err = d.GetAttachment(context.Background(), attachmentID)
	assert.NoError(t, err)

	assert.NoError(t, d.DeleteAttachment(context.Background(), attachmentID))
	_, err = d.GetAttachment(context.Background(), attachmentID)
	assert.Error(t, err)

	assert.NoError(t, d.Close())

	tempFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	assert.NoError(t, formats.SQLiteAttachments.Initialize(tempFile.Name(), 0, 0))
	assert.Error(t, formats.SQLiteAttachments.Initialize("", 0, 0))
}

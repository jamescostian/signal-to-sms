package attachments_test

import (
	"context"
	"testing"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/stretchr/testify/assert"
)

func TestDiscard(t *testing.T) {
	d := attachments.Discard()

	assert.NoError(t, d.SetAttachment(context.Background(), 0, []byte("hi")))

	ids, err := d.ListAttachments(context.Background())
	assert.NoError(t, err, "Not having any attachments is fine")
	assert.Len(t, ids, 0, "There should not be any attachments listed")

	assert.NoError(t, d.DeleteAttachment(context.Background(), 0), "There's nothing stored, so it's already gone")

	_, err = d.GetAttachment(context.Background(), 0)
	assert.Error(t, err, "You can't get an attachment that hasn't been actually stored")

	assert.NoError(t, d.Close(), "Should not have any issues closing")
}

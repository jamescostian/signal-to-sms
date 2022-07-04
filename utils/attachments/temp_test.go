package attachments_test

import (
	"context"
	"testing"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTempStore(t *testing.T) {
	s, err := attachments.NewTempStore()
	require.NoError(t, err, "Should be able to create a temp store")

	ids, err := s.ListAttachments(context.Background())
	assert.NoError(t, err, "Should be able to get a list of attachments immediately")
	assert.Len(t, ids, 0, "There shouldn't be any attachments initially")
	_, err = s.GetAttachment(context.Background(), 0)
	assert.Error(t, err, "Should not be able to get an attachment that doesn't exist")

	err = s.SetAttachment(context.Background(), 0, []byte("it works!"))
	assert.NoError(t, err, "Should be able to set an attachment")
	data, err := s.GetAttachment(context.Background(), 0)
	assert.NoError(t, err, "Should be able to get an attachment that was previously set")
	assert.Equal(t, "it works!", string(data))

	err = s.DeleteAttachment(context.Background(), 0)
	assert.NoError(t, err, "Should be able to delete an attachment")
	_, err = s.GetAttachment(context.Background(), 0)
	assert.Error(t, err, "Should not be able to get an attachment that was deleted after being set")

	assert.NoError(t, s.Close(), "Should be able to close the temporary attachment store")
}

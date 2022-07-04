package attachments_test

import (
	"context"
	"testing"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/stretchr/testify/assert"
)

func TestNoReads(t *testing.T) {
	mock := &attachments.MockStore{}
	mock.On("SetAttachment", context.Background(), int64(0), []byte(nil)).Return(nil)
	mock.On("DeleteAttachment", context.Background(), int64(0)).Return(nil)
	mock.On("Close").Return(nil)
	s := attachments.NoReads(mock)

	assert.NoError(t, s.SetAttachment(context.Background(), 0, nil), "Should not have an issue writing (SetAttachment)")
	assert.NoError(t, s.Close(), "Should not have any issues closing")
	assert.NoError(t, s.DeleteAttachment(context.Background(), 0), "Should not have an issue writing (DeletAttachment)")

	_, err := s.ListAttachments(context.Background())
	assert.Error(t, err, "Should not be able to read (ListAttachments)")
	_, err = s.GetAttachment(context.Background(), 0)
	assert.Error(t, err, "Should not be able to read (GetAttachment)")
}

func TestNoWrites(t *testing.T) {
	mock := &attachments.MockStore{}
	mock.On("GetAttachment", context.Background(), int64(0)).Return(nil, nil)
	mock.On("ListAttachments", context.Background()).Return(nil, nil)
	mock.On("Close").Return(nil)
	s := attachments.NoWrites(mock)

	_, err := s.ListAttachments(context.Background())
	assert.NoError(t, err, "Should be able to read (ListAttachments)")
	_, err = s.GetAttachment(context.Background(), 0)
	assert.NoError(t, err, "Should be able to read (GetAttachment)")
	assert.NoError(t, s.Close(), "Should not have any issues closing")

	assert.Error(t, s.SetAttachment(context.Background(), 0, nil), "Should not be able to write (SetAttachment)")
	assert.Error(t, s.DeleteAttachment(context.Background(), 0), "Should not be able to write (DeletAttachment)")
}

func TestUnclosable(t *testing.T) {
	mock := &attachments.MockStore{}
	mock.On("GetAttachment", context.Background(), int64(0)).Return(nil, nil)
	mock.On("ListAttachments", context.Background()).Return(nil, nil)
	mock.On("SetAttachment", context.Background(), int64(0), []byte(nil)).Return(nil)
	mock.On("DeleteAttachment", context.Background(), int64(0)).Return(nil)
	s := attachments.Unclosable(mock)

	assert.NoError(t, s.SetAttachment(context.Background(), 0, nil), "Should not have an issue writing (SetAttachment)")
	assert.NoError(t, s.DeleteAttachment(context.Background(), 0), "Should not have an issue writing (DeletAttachment)")
	_, err := s.ListAttachments(context.Background())
	assert.NoError(t, err, "Should be able to read (ListAttachments)")
	_, err = s.GetAttachment(context.Background(), 0)
	assert.NoError(t, err, "Should be able to read (GetAttachment)")

	assert.Error(t, s.Close(), "Should not be able to close")
}

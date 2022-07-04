package attachments_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	ids := []int64{1, 23, 456}
	data := [][]byte{{'a'}, {'b'}, {'c'}}
	t.Run("success", func(t *testing.T) {
		src, dest := &attachments.MockStore{}, &attachments.MockStore{}

		src.On("ListAttachments", context.Background()).Return(ids, nil)
		src.On("GetAttachment", context.Background(), ids[0]).Return(data[0], nil)
		src.On("GetAttachment", context.Background(), ids[1]).Return(data[1], nil)
		src.On("GetAttachment", context.Background(), ids[2]).Return(data[2], nil)

		dest.On("SetAttachment", context.Background(), ids[0], data[0]).Return(nil)
		dest.On("SetAttachment", context.Background(), ids[1], data[1]).Return(nil)
		dest.On("SetAttachment", context.Background(), ids[2], data[2]).Return(nil)

		assert.NoError(t, attachments.Copy(context.Background(), dest, src))
	})
	t.Run("cannot_list_attachments", func(t *testing.T) {
		src, dest := &attachments.MockStore{}, &attachments.MockStore{}

		src.On("ListAttachments", context.Background()).Return(nil, errors.New("error"))
		assert.Error(t, attachments.Copy(context.Background(), dest, src))
	})
	t.Run("cannot_get_attachments", func(t *testing.T) {
		src, dest := &attachments.MockStore{}, &attachments.MockStore{}

		src.On("ListAttachments", context.Background()).Return(ids, nil)
		src.On("GetAttachment", context.Background(), ids[0]).Return(nil, errors.New("error"))

		assert.Error(t, attachments.Copy(context.Background(), dest, src))
	})
	t.Run("cannot_set_attachments", func(t *testing.T) {
		src, dest := &attachments.MockStore{}, &attachments.MockStore{}

		src.On("ListAttachments", context.Background()).Return(ids, nil)
		src.On("GetAttachment", context.Background(), ids[0]).Return(data[0], nil)
		src.On("GetAttachment", context.Background(), ids[1]).Return(data[1], nil)
		src.On("GetAttachment", context.Background(), ids[2]).Return(data[2], nil)

		dest.On("SetAttachment", context.Background(), ids[0], data[0]).Return(errors.New("error"))

		assert.Error(t, attachments.Copy(context.Background(), dest, src))
	})
}

package mimetoext_test

import (
	"testing"

	"github.com/jamescostian/signal-to-sms/internal/fromsignalback/mimetoext"
	"github.com/stretchr/testify/assert"
)

func TestMIMEToExt(t *testing.T) {
	ext, err := mimetoext.Convert("image/jpeg")
	assert.Equal(t, "jpg", ext)
	assert.NoError(t, err)

	ext, err = mimetoext.Convert("video/mp4")
	assert.Equal(t, "mp4", ext)
	assert.NoError(t, err)

	ext, err = mimetoext.Convert("")
	assert.Equal(t, "", ext)
	assert.Error(t, err)

	ext, err = mimetoext.Convert("idk!!!!!!!!!")
	assert.Equal(t, "", ext)
	assert.Error(t, err)
}

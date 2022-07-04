package formats_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/jamescostian/signal-to-sms/formats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenFileForReadsAndWrites(t *testing.T) {
	tempF, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	name := tempF.Name()
	tempF.Close()
	defer os.Remove(name)

	f, err := formats.OpenFileForWrites(name, os.O_WRONLY, 0600)
	require.NoError(t, err)
	_, err = f.(*os.File).Write([]byte("hi"))
	require.NoError(t, err)
	f.Close()

	f, err = formats.OpenFileForReads(name)
	require.NoError(t, err)
	dataRead, err := ioutil.ReadAll(f.(*os.File))
	assert.Equal(t, "hi", string(dataRead))
	assert.NoError(t, err)
	f.Close()
}

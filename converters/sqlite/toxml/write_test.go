package toxml_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jamescostian/signal-to-sms/converters/sqlite/toxml"
	xmlData "github.com/jamescostian/signal-to-sms/internal/testdata/xml"
	"github.com/jamescostian/signal-to-sms/internal/testutil"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	t.Run("success_with_testdata", func(t *testing.T) {
		db := testutil.SetupSeededMsgDB(t)
		store := testutil.SetupSeededAttachmentStore(t)
		out, err := ioutil.TempFile("", "")
		require.NoError(t, err)
		defer os.Remove(out.Name())
		assert.NoError(t, toxml.Write(context.Background(), db, store, out, xmlData.ExamplePhone, xmlData.ExampleIndent))

		// Now compare the output
		outBytes, err := os.ReadFile(out.Name())
		require.NoError(t, err)
		expectedBytes, err := os.ReadFile(xmlData.ExamplePath)
		require.NoError(t, err)
		assert.Equal(t, string(expectedBytes), string(outBytes))
	})
	t.Run("error_if_file_is_closed", func(t *testing.T) {
		db := testutil.SetupSeededMsgDB(t)
		store := testutil.SetupSeededAttachmentStore(t)
		out, err := ioutil.TempFile("", "")
		require.NoError(t, err)
		require.NoError(t, out.Close())
		os.Remove(out.Name())
		assert.Error(t, toxml.Write(context.Background(), db, store, out, xmlData.ExamplePhone, xmlData.ExampleIndent))
	})
	t.Run("error_if_msgs_is_messed_up", func(t *testing.T) {
		db, delete, err := sqlite.CreateTemp()
		require.NoError(t, err)
		defer delete()
		store := testutil.SetupSeededAttachmentStore(t)
		out, err := ioutil.TempFile("", "")
		require.NoError(t, err)
		defer os.Remove(out.Name())
		assert.Error(t, toxml.Write(context.Background(), db, store, out, xmlData.ExamplePhone, xmlData.ExampleIndent))
	})
}

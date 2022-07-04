package toxml_test

import (
	"context"
	"testing"

	"github.com/jamescostian/signal-to-sms/converters/sqlite/toxml"
	"github.com/jamescostian/signal-to-sms/internal/testutil"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRecipients(t *testing.T) {
	t.Run("works_with_testdata", func(t *testing.T) {
		db := testutil.SetupSeededMsgDB(t)
		out, err := toxml.GetRecipients(context.Background(), db)
		if assert.NoError(t, err) {
			assert.Equal(t, map[int]string{
				1: "+16054756968",
				2: "+16186258313",
				3: "40404",
			}, map[int]string(out))
		}
	})
	t.Run("error_querying_breaks_everything", func(t *testing.T) {
		db, delete, err := sqlite.CreateTemp()
		require.NoError(t, err)
		defer delete()
		_, err = toxml.GetRecipients(context.Background(), db)
		assert.Error(t, err)
	})
}

package sqlite_test

import (
	"context"
	"testing"

	"github.com/jamescostian/signal-to-sms/internal/testutil"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteEmptyMessages(t *testing.T) {
	db := testutil.SetupSeededMsgDB(t)
	assert.NoError(t, sqlite.DeleteEmptyMessages(context.Background(), db))

	queryAllInts := func(sql string) (ints []int) {
		rows, err := db.Query(sql)
		require.NoError(t, err, "error querying the DB as part of the test")
		var num int
		for rows.Next() {
			require.NoError(t, rows.Scan(&num))
			ints = append(ints, num)
		}
		return
	}

	assert.Equal(t, []int{1, 4}, queryAllInts("SELECT _id FROM sms ORDER BY _id"))
	assert.Equal(t, []int{1, 3}, queryAllInts("SELECT _id FROM mms ORDER BY _id"))
}

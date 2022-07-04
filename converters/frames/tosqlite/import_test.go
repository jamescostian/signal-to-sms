package tosqlite_test

import (
	"context"
	"testing"

	"github.com/jamescostian/signal-to-sms/converters/frames/tosqlite"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrameWriter(t *testing.T) {
	db, delete, err := sqlite.CreateTemp()
	require.NoError(t, err)
	defer delete()

	importer := tosqlite.NewFrameWriter(context.Background(), db, true)
	checkedIfShouldIgnore := []string{}
	tosqlite.ChangeShouldIgnoreImport(importer, func(statement string) bool {
		checkedIfShouldIgnore = append(checkedIfShouldIgnore, statement)
		return statement == "ignore me"
	})

	// In go you can't write stuff like &true, so the following allow working around that easily
	tru, str, intgr, flt, blob := true, "hi", uint64(42), float64(3.14), []byte("hi")

	t.Run("creating_a_table_and_inserting_data_work", func(t *testing.T) {
		checkedIfShouldIgnore = nil
		createTbl := "CREATE TABLE test (_id INTEGER PRIMARY KEY, num INTEGER, flt REAL, str TEXT, blb BLOB)"
		assert.NoError(t, importer.Write(&signal.BackupFrame{Statement: &signal.SqlStatement{Statement: &createTbl}}, nil))
		insert := "INSERT INTO test VALUES (?, ?, ?, ?, ?)"
		assert.NoError(t, importer.Write(&signal.BackupFrame{Statement: &signal.SqlStatement{
			Statement: &insert,
			Parameters: []*signal.SqlStatement_SqlParameter{
				{IntegerParameter: &intgr},
				{Nullparameter: &tru},
				{DoubleParameter: &flt},
				{StringParamter: &str},
				{BlobParameter: blob},
			},
		}}, nil))
		assert.Equal(t, []string{createTbl, insert}, checkedIfShouldIgnore)
	})

	t.Run("respects_if_told_to_ignore_statements", func(t *testing.T) {
		checkedIfShouldIgnore = nil
		ignoreMe := "ignore me"
		assert.NoError(t, importer.Write(&signal.BackupFrame{Statement: &signal.SqlStatement{Statement: &ignoreMe}}, nil))
		assert.NoError(t, importer.Write(&signal.BackupFrame{Statement: &signal.SqlStatement{
			Statement:  &ignoreMe,
			Parameters: []*signal.SqlStatement_SqlParameter{{Nullparameter: &tru}},
		}}, nil))
		assert.Equal(t, []string{ignoreMe, ignoreMe}, checkedIfShouldIgnore)
	})

	t.Run("finds_invalid_sql", func(t *testing.T) {
		checkedIfShouldIgnore = nil
		sql := "idk??"
		assert.Error(t, importer.Write(&signal.BackupFrame{Statement: &signal.SqlStatement{Statement: &sql}}, nil))
		assert.Error(t, importer.Write(&signal.BackupFrame{Statement: &signal.SqlStatement{
			Statement:  &sql,
			Parameters: []*signal.SqlStatement_SqlParameter{{Nullparameter: &tru}},
		}}, nil))
		assert.Equal(t, []string{sql, sql}, checkedIfShouldIgnore)
	})

	t.Run("finds_invalid_parameters", func(t *testing.T) {
		insert := "INSERT INTO test VALUES (?, ?, ?, ?, ?)"
		assert.Error(t, importer.Write(&signal.BackupFrame{Statement: &signal.SqlStatement{
			Statement:  &insert,
			Parameters: []*signal.SqlStatement_SqlParameter{{}, {}, {}, {}, {}},
		}}, nil))
	})
}

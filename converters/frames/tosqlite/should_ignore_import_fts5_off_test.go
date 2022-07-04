//go:build !sqlite_fts5
// +build !sqlite_fts5

package tosqlite_test

import (
	"testing"

	"github.com/jamescostian/signal-to-sms/converters/frames/tosqlite"
	"github.com/stretchr/testify/assert"
)

func TestShouldIgnoreImport(t *testing.T) {
	assert.False(t, tosqlite.ShouldIgnoreImport("CREATE TABLE sms (_id INTEGER PRIMARY KEY AUTOINCREMENT)"))
	assert.False(t, tosqlite.ShouldIgnoreImport("INSERT INTO sms (3)"))

	assert.True(t, tosqlite.ShouldIgnoreImport("CREATE TABLE sqlite_sequence(name,seq)"))
	assert.True(t, tosqlite.ShouldIgnoreImport("CREATE VIRTUAL TABLE sms_fts USING fts5(body, thread_id UNINDEXED, content=sms, content_rowid=_id)"))
	assert.True(t, tosqlite.ShouldIgnoreImport("CREATE TRIGGER mms_ai AFTER INSERT ON mms BEGIN INSERT INTO mms_fts(rowid, body, thread_id) VALUES (new._id, new.body, new.thread_id); END"))

	assert.False(t, tosqlite.ShouldIgnoreImport("CREATE INDEX part_mms_id_index ON part (mid)"))
}

func TestFastShouldIgnoreImport(t *testing.T) {
	assert.False(t, tosqlite.FastShouldIgnoreImport("CREATE TABLE sms (_id INTEGER PRIMARY KEY AUTOINCREMENT)"))
	assert.False(t, tosqlite.FastShouldIgnoreImport("INSERT INTO sms (3)"))

	assert.True(t, tosqlite.FastShouldIgnoreImport("CREATE TABLE sqlite_sequence(name,seq)"))
	assert.True(t, tosqlite.FastShouldIgnoreImport("CREATE VIRTUAL TABLE sms_fts USING fts5(body, thread_id UNINDEXED, content=sms, content_rowid=_id)"))
	assert.True(t, tosqlite.FastShouldIgnoreImport("CREATE TRIGGER mms_ai AFTER INSERT ON mms BEGIN INSERT INTO mms_fts(rowid, body, thread_id) VALUES (new._id, new.body, new.thread_id); END"))

	assert.True(t, tosqlite.FastShouldIgnoreImport("CREATE INDEX part_mms_id_index ON part (mid)"))
}

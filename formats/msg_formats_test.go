package formats_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/jamescostian/signal-to-sms/formats"
	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestPrototextMsg(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	name := file.Name()
	os.Remove(name)

	writeable, err := formats.PrototextMsg.OpenForWrites(name, os.O_CREATE|os.O_WRONLY, 0600)
	require.NoError(t, err)
	defer os.Remove(name)
	assert.NoError(t, writeable.(frameio.FrameWriteCloser).Write(&signal.BackupFrame{Attachment: &signal.Attachment{RowId: proto.Uint64(123)}}, nil))
	assert.NoError(t, writeable.Close())

	readable, err := formats.PrototextMsg.OpenForReads(name)
	require.NoError(t, err)
	var frame signal.BackupFrame
	assert.NoError(t, readable.(frameio.FrameReadCloser).Read(&frame, nil))
	assert.NoError(t, readable.Close())
	assert.Equal(t, uint64(123), frame.Attachment.GetRowId())
}

func TestSQLiteMsg(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	name := file.Name()
	os.Remove(name)

	writeable, err := formats.SQLiteMsg.OpenForWrites(name, os.O_CREATE|os.O_WRONLY, 0600)
	require.NoError(t, err)
	defer os.Remove(name)
	_, err = writeable.(*sqlite.ClosableTx).Exec("CREATE TABLE foo (bar INTEGER)")
	require.NoError(t, err)
	_, err = writeable.(*sqlite.ClosableTx).Exec("INSERT INTO foo VALUES (123)")
	require.NoError(t, err)
	assert.NoError(t, writeable.Close())

	readable, err := formats.SQLiteMsg.OpenForReads(name)
	require.NoError(t, err)
	row := readable.(*sqlite.ClosableTx).QueryRow("SELECT * FROM foo")
	var valueRead int
	assert.NoError(t, row.Scan(&valueRead))
	assert.Equal(t, 123, valueRead)
	assert.NoError(t, readable.Close())
}

package conversion_test

import (
	"testing"

	"github.com/jamescostian/signal-to-sms/conversion"
	"github.com/jamescostian/signal-to-sms/formats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindPath(t *testing.T) {
	// Set up: build the graph
	g, err := conversion.NewGraph([]conversion.MsgConverter{
		{InputFormat: formats.Encrypted, OutputFormat: formats.BackupFrame},
		{InputFormat: formats.BackupFrame, OutputFormat: formats.Encrypted},
		{InputFormat: formats.BackupFrame, OutputFormat: formats.PrototextMsg},
		{InputFormat: formats.PrototextMsg, OutputFormat: formats.BackupFrame},
		{InputFormat: formats.BackupFrame, OutputFormat: formats.SQLiteMsgAndAttachment},
		{InputFormat: formats.BackupFrame, OutputFormat: formats.SQLiteMsg},
		{InputFormat: formats.SQLiteMsgAndAttachment, OutputFormat: formats.XML},
		{InputFormat: formats.SQLiteMsg, OutputFormat: formats.XML},
	}, []formats.MsgFormat{
		formats.BackupFrame,
		formats.Encrypted,
		formats.PrototextMsg,
		formats.SQLiteMsg,
		formats.SQLiteMsgAndAttachment,
		formats.XML,
	})
	require.NoError(t, err)

	// Now actually find the path
	path, err := g.FindPath(formats.Encrypted, formats.SQLiteMsgAndAttachment)
	require.NoError(t, err)
	require.Len(t, path, 2)
	assert.Equal(t, formats.Encrypted.Name, path[0].InputFormat.Name)
	assert.Equal(t, formats.BackupFrame.Name, path[0].OutputFormat.Name)
	assert.Equal(t, formats.BackupFrame.Name, path[1].InputFormat.Name)
	assert.Equal(t, formats.SQLiteMsgAndAttachment.Name, path[1].OutputFormat.Name)

	// Try to find a path that doesn't exist, make sure it can't be found
	_, err = g.FindPath(formats.XML, formats.Encrypted)
	assert.Error(t, err)
}

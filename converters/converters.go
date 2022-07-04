package converters

import (
	"context"
	"os"

	"github.com/jamescostian/signal-to-sms/conversion"
	"github.com/jamescostian/signal-to-sms/converters/encrypted/toframes"
	"github.com/jamescostian/signal-to-sms/converters/frames/toattachments"
	"github.com/jamescostian/signal-to-sms/converters/frames/toencrypted"
	"github.com/jamescostian/signal-to-sms/converters/frames/tosqlite"
	"github.com/jamescostian/signal-to-sms/converters/sqlite/toxml"
	"github.com/jamescostian/signal-to-sms/formats"
	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
)

var MsgConverter = []conversion.MsgConverter{
	{
		InputFormat:  formats.Encrypted,
		OutputFormat: formats.BackupFrame,
		Convert: func(ctx context.Context, c *conversion.Conversion) (interface{}, error) {
			return toframes.NewFrameReader(c.MsgIn.(*os.File), ctx.Value(formats.SignalBackupDBPasswordCtxKey).(string))
		},
	},
	{
		InputFormat:  formats.BackupFrame,
		OutputFormat: formats.Encrypted,
		Convert: func(ctx context.Context, c *conversion.Conversion) (interface{}, error) {
			return nil, toencrypted.Encrypt(ctx, c.MsgIn.(frameio.FrameReader), c.MsgOut.(*os.File), ctx.Value(formats.SignalBackupDBPasswordCtxKey).(string))
		},
	},
	{
		InputFormat:  formats.BackupFrame,
		OutputFormat: formats.PrototextMsg,
		Convert: func(ctx context.Context, c *conversion.Conversion) (interface{}, error) {
			writeFramesToAttachments := toattachments.NewFrameWriter(ctx, c.ExternalAttachments)
			writeToMsgOutAndAttachments := frameio.NewMultiWriter(c.MsgOut.(frameio.FrameWriteCloser), writeFramesToAttachments)
			return nil, frameio.CopyFrames(ctx, writeToMsgOutAndAttachments, c.MsgIn.(frameio.FrameReader))
		},
	},
	{
		InputFormat:  formats.PrototextMsg,
		OutputFormat: formats.BackupFrame,
		Convert: func(ctx context.Context, c *conversion.Conversion) (interface{}, error) {
			return frameio.NewMergedReader(ctx, c.MsgIn.(frameio.FrameReader), c.ExternalAttachments), nil
		},
	},
	{
		InputFormat:  formats.BackupFrame,
		OutputFormat: formats.SQLiteMsgAndAttachment,
		Convert: func(ctx context.Context, c *conversion.Conversion) (interface{}, error) {
			dbOrTx := c.MsgOut.(sqlite.DBOrTx)
			store, err := attachments.NewSQLStore(dbOrTx)
			if err != nil {
				return nil, err
			}
			return backupFrameToSQLite(ctx, c, store, dbOrTx)
		},
	},
	{
		InputFormat:  formats.BackupFrame,
		OutputFormat: formats.SQLiteMsg,
		Convert: func(ctx context.Context, c *conversion.Conversion) (interface{}, error) {
			return backupFrameToSQLite(ctx, c, c.ExternalAttachments, c.MsgOut.(sqlite.DBOrTx))
		},
	},
	{
		InputFormat:  formats.SQLiteMsgAndAttachment,
		OutputFormat: formats.XML,
		Convert: func(ctx context.Context, c *conversion.Conversion) (interface{}, error) {
			dbOrTx := c.MsgIn.(sqlite.DBOrTx)
			store, err := attachments.NewSQLStore(dbOrTx)
			if err != nil {
				return nil, err
			}
			return sqliteToXML(ctx, c, store, dbOrTx)
		},
	},
	{
		InputFormat:  formats.SQLiteMsg,
		OutputFormat: formats.XML,
		Convert: func(ctx context.Context, c *conversion.Conversion) (interface{}, error) {
			return sqliteToXML(ctx, c, c.ExternalAttachments, c.MsgIn.(sqlite.DBOrTx))
		},
	},
}

func backupFrameToSQLite(ctx context.Context, c *conversion.Conversion, store attachments.Store, dbOrTx sqlite.DBOrTx) (interface{}, error) {
	frameReader := c.MsgIn.(frameio.FrameReader)
	writeFramesToSQLite := tosqlite.NewFrameWriter(ctx, dbOrTx, c.IsFinalConversion())
	writeFramesToAttachments := toattachments.NewFrameWriter(ctx, store)
	return nil, frameio.CopyFrames(ctx, frameio.NewMultiWriter(writeFramesToSQLite, writeFramesToAttachments), frameReader)
}

const xmlIndent = "\t"

func sqliteToXML(ctx context.Context, c *conversion.Conversion, store attachments.Store, dbOrTx sqlite.DBOrTx) (interface{}, error) {
	if err := sqlite.DeleteEmptyMessages(ctx, dbOrTx); err != nil {
		return nil, err
	}
	return nil, toxml.Write(ctx, dbOrTx, store, c.MsgOut.(*os.File), ctx.Value(formats.MyPhoneNumberCtxKey).(string), xmlIndent)
}

package formats

import (
	"io"
	"os"

	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/prototext"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
)

type MsgFormat struct {
	IncludesAttachments bool
	// Should be unique across *every* MsgFormat, regardless of which map it's in, will be in logs
	Name string

	// Given the path to a file in this format, provide whatever you need to read from it (e.g. an *os.File), as long as it has a Close method.
	// Not applicable for some formats; those formats simply can't be used as the starting input of a conversion, but can be used during in-between steps.
	OpenForReads func(path string) (io.Closer, error)

	// Given the parameters one might use with os., provide whatever you need to write to it (e.g. an *os.File), as long as it has a Close method.
	// Not applicable for some formats; those formats simply can't be used as the final output of a conversion, but can be used during in-between steps.
	OpenForWrites func(path string, flag int, perm os.FileMode) (io.Closer, error)

	// E.g. []interface{}{SignalBackupDBPasswordCtxKey, MyPhoneNumberCtxKey}
	RequiredCtxKeys []interface{}
}

type ctxKey string

var (
	SignalBackupDBPasswordCtxKey ctxKey = "SignalBackupDBPasswordCtxKey"
	MyPhoneNumberCtxKey          ctxKey = "MyPhoneNumberCtxKey"
)

var PrototextMsg = MsgFormat{
	Name: "prototext",
	OpenForReads: func(path string) (io.Closer, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return frameio.NewBackupFramesReadCloser(file, prototext.PrototextEncoder)
	},
	OpenForWrites: func(path string, flag int, perm os.FileMode) (io.Closer, error) {
		file, err := os.OpenFile(path, flag, perm)
		return frameio.NewBackupFramesWriteCloser(file, prototext.PrototextEncoder), err
	},
}
var SQLiteMsg = MsgFormat{
	Name: "sqlite (attachments separate)",
	OpenForReads: func(path string) (io.Closer, error) {
		return sqlite.OpenToTx(path)
	},
	OpenForWrites: func(path string, flag int, perm os.FileMode) (io.Closer, error) {
		db, err := sqlite.Create(path, flag, perm)
		if err != nil {
			return nil, err
		}
		if err := db.Close(); err != nil {
			return nil, err
		}
		return sqlite.OpenToTx(path)
	},
}

var MsgOnlyFormats = map[string]MsgFormat{
	"prototext": PrototextMsg,
	"sqlite":    SQLiteMsg,
}

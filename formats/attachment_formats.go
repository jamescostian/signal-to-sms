package formats

import (
	"os"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
)

// AttachmentFormat provides instructions on how to set up a new attachment store in this format (Initialize), as well as how to open an existing or recently initialized store in this format (Open)
type AttachmentFormat struct {
	Open       func(path string) (attachments.Store, error)
	Initialize func(path string, flag int, perm os.FileMode) error
}

var DiscardAttachments = AttachmentFormat{
	Open: func(_ string) (attachments.Store, error) {
		return attachments.Discard(), nil
	},
	Initialize: func(path string, flag int, perm os.FileMode) error {
		return nil
	},
}
var SQLiteAttachments = AttachmentFormat{
	Open: func(path string) (attachments.Store, error) {
		db, err := sqlite.Open(path)
		if err != nil {
			return nil, err
		}
		return attachments.NewSQLStore(db)
	},
	Initialize: func(path string, flag int, perm os.FileMode) error {
		file, err := os.OpenFile(path, flag, perm)
		if err != nil {
			return err
		}
		return file.Close()
	},
}

var AttachmentFormats = map[string]AttachmentFormat{
	"discard": DiscardAttachments,
	"sqlite":  SQLiteAttachments,
}

package toxml

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/pkg/errors"
)

// Write converts SMSes and MMSes in a *sql.DB/Tx and AttachmentStore to XML, and streams the XML it to an io.Writer.
func Write(ctx context.Context, dbOrTx sqlite.DBOrTx, attachments attachments.Store, out io.Writer, myPhoneNumber string, indent string) error {
	w, err := newWriter(ctx, dbOrTx, attachments, out, myPhoneNumber, indent)
	if err != nil {
		return err
	}

	if err := w.WriteHeader(ctx); err != nil {
		return err
	}

	if err := w.MarshalSMSes(ctx); err != nil {
		return errors.Wrap(err, "error converting SMSes to XML")
	}
	if err := w.MarshalMMSes(ctx); err != nil {
		return errors.Wrap(err, "error converting MMSes to XML")
	}

	return w.WriteFooter()
}

func newWriter(ctx context.Context, dbOrTx sqlite.DBOrTx, attachments attachments.Store, out io.Writer, myPhoneNumber string, indent string) (*writer, error) {
	w := writer{
		dbOrTx:        dbOrTx,
		attachments:   attachments,
		out:           out,
		encoder:       xml.NewEncoder(out),
		myPhoneNumber: myPhoneNumber,
	}
	w.encoder.Indent(indent, indent)

	// Cache a recipient lookup map for usage later on - the DB is giving some recipient IDs that need to be converted to phone numbers, and you can't JOIN on something like "123,456"
	recipients, err := getRecipients(ctx, w.dbOrTx)
	if err != nil {
		return nil, err
	}
	w.recipients = recipients
	return &w, nil
}

type writer struct {
	dbOrTx        sqlite.DBOrTx
	attachments   attachments.Store
	out           io.Writer
	encoder       *xml.Encoder
	myPhoneNumber string
	recipients    recipients

	// For avoiding allocations :)
	mmsRecipientIDs []int
	mmsPhoneNumbers []string
}

func (w *writer) WriteHeader(ctx context.Context) error {
	// The header includes the number of messages, so first grab the messages
	var numMessages int64
	row := w.dbOrTx.QueryRowContext(ctx, "SELECT (SELECT COUNT(_id) FROM mms) + (SELECT COUNT(_id) FROM sms)")
	if err := row.Scan(&numMessages); err != nil {
		return err
	}
	// Write the head that SMS Backup & Restore uses in its backups, to help match them, and then add the opening <sms> tag
	// A newline is needed at the end because encoding/xml won't add a newline before encoding the very first SMS/MMS
	_, err := fmt.Fprintf(w.out, "<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>\n<smses count=\"%v\">\n", numMessages)
	return err
}

func (w *writer) WriteFooter() error {
	// A preceding newline is needed because encoding/xml will end without a newline
	_, err := w.out.Write([]byte("\n</smses>"))
	return err
}

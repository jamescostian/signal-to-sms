package toxml

import (
	"context"

	"github.com/jamescostian/signal-to-sms/internal/fromsignalback/smstype"
)

// MarshalSMSes will loop through each SMS in the DB and let you run a function on each SMS.
// Like go's for loops, the variable you get here will be reused, so make a copy of it if you need it to last after your function returns!
// If there's an error getting something from the DB, or if your function returns an error, then iteration will halt, and the error will be returned by MarshalSMSes.
func (w *writer) MarshalSMSes(ctx context.Context) error {
	rows, err := w.dbOrTx.QueryContext(ctx, `
		SELECT
			COALESCE(sms.protocol, 0),
			recipient.phone,
			sms.date,
			COALESCE(sms.type, 0),
			sms.subject,
			sms.body,
			sms.service_center,
			sms.read,
			sms.status,
			sms.date_sent
		FROM
			sms
			INNER JOIN recipient ON sms.address = recipient._id
		WHERE
			body != '' AND body IS NOT NULL
		ORDER BY
			sms._id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var sms sms
	for rows.Next() {
		if err = rows.Scan(
			&sms.Protocol,
			&sms.PhoneNumber,
			&sms.DateReceived,
			&sms.Type,
			&sms.Subject,
			&sms.Body,
			&sms.ServiceCenter,
			&sms.Read,
			&sms.Status,
			&sms.DateSent,
		); err != nil {
			return err
		}
		if sms.Protocol == signalSMSProtocol {
			sms.Protocol = 0
		}
		if sms.Type, err = smstype.TranslateSMSType(int64(sms.Type)); err != nil {
			return err
		}
		if err = w.encoder.Encode(&sms); err != nil {
			return err
		}
	}
	return nil
}

package sqlite

import "context"

// DeleteEmptyMessages deletes SMSes/MMSes that have no body and no attachments
func DeleteEmptyMessages(ctx context.Context, dbOrTx DBOrTx) error {
	_, err := dbOrTx.ExecContext(ctx, "DELETE FROM sms WHERE body = '' OR body IS NULL")
	if err != nil {
		return err
	}
	_, err = dbOrTx.ExecContext(ctx, "DELETE FROM mms WHERE (body = '' OR body IS NULL) AND _id NOT IN (SELECT mid FROM part GROUP BY mid)")
	return err
}

package toxml

import (
	"context"

	"github.com/jamescostian/signal-to-sms/utils/sqlite"
)

// Maps recipient IDs to phone numbers
type recipients map[int]string

// This is basically your contact list, so there's going to be a good amount of rows here.
// It's faster to start out with a healthy sized map, rather than slowly work your way up via many allocations.
const initialRecipientsSize = 1024

// getRecipients returns a map from recipient IDs to their phone numbers
func getRecipients(ctx context.Context, dbOrTx sqlite.DBOrTx) (recipients, error) {
	rcpnts := make(recipients, initialRecipientsSize)
	rows, err := dbOrTx.QueryContext(ctx, "SELECT _id, phone FROM recipient WHERE phone IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		id    int
		phone string
	)
	for rows.Next() {
		if err := rows.Scan(&id, &phone); err != nil {
			return nil, err
		}
		rcpnts[id] = phone
	}
	return rcpnts, nil
}

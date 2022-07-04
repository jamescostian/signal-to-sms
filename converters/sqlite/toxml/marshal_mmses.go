package toxml

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// MarshalMMSes will loop through each MMS in the DB and marshal it into the writer
func (w *writer) MarshalMMSes(ctx context.Context) error {
	rows, err := w.dbOrTx.QueryContext(ctx, `
		SELECT
			mms.date_received,
			mms.read,
			mms.ct_l,
			mms.tr_id,
			mms._id,
			mms.date / 1000,
			mms.read,
			COALESCE(mms.m_type, ""),
			IIF(
				mms.exp IS NOT NULL,
				cast(mms.exp as text),
				mms.exp
			),
			COALESCE(SUM(part.data_size), 0),
			mms.body,
			mms.address,
			thread.thread_recipient_id,
			groups.members
		FROM
			mms
			LEFT JOIN thread ON mms.thread_id = thread._id
			LEFT JOIN groups ON thread.thread_recipient_id = groups.recipient_id
			LEFT JOIN part ON mms._id = part.mid
		WHERE
			(body != '' AND body IS NOT NULL)
			OR mms._id IN (
				SELECT
					mid
				FROM
					part
				GROUP BY
					mid
			)
		GROUP BY
			mms._id
		ORDER BY
			mms.part_count DESC,
			mms._id ASC`)
	if err != nil {
		return err
	}
	defer rows.Close()

	// The following variables are all just reused throughout iterations of the loop to avoid allocations per MMS.
	// Group chats count as MMSes (can't have a group SMS), so MMSes are super common and need to be optimized.
	var (
		mms = mms{
			SubjectCharset: utf8,
			MessageClass:   "personal",
			ContentType:    wap,
			DeliveryReport: mmsNo,
			Priority:       mmsPriorityNormal,
			ReadReport:     mmsNo,
			ResponseStatus: mmsResponseOk,
			MMSVersion:     normalMMSVersion,
		}

		id   int64
		body *string

		primaryRecipientID int
		members            *string
		threadRecipientID  *int
	)

	for rows.Next() {
		if err = rows.Scan(
			&mms.DateReceived,
			&mms.Read,
			&mms.ContentLocation,
			&mms.TransactionID,
			&id,
			&mms.DateSent,
			&mms.Seen,
			&mms.Type,
			&mms.Expiry,
			&mms.Size,
			&body,
			&primaryRecipientID,
			&threadRecipientID,
			&members,
		); err != nil {
			return err
		}

		if mms.Type == mmsSendReq {
			mms.MsgBox = mmsSentBox
		} else {
			mms.MsgBox = mmsInbox
		}

		if err = w.setupParts(ctx, id, &mms, body); err != nil {
			return err
		}

		if err = w.setupPhoneNumbers(primaryRecipientID, threadRecipientID, members, &mms); err != nil {
			return err
		}

		if err = w.encoder.Encode(&mms); err != nil {
			return err
		}
	}
	return nil
}

func (w *writer) setupParts(ctx context.Context, id int64, mms *mms, body *string) error {
	// Reset the parts slice to zero (keeping the underlying slice in place to avoid allocations), and add parts to it from the DB
	mms.Parts = mms.Parts[:0]
	err := w.ForEachPartOf(ctx, id, func(part *mmsPart) error {
		mms.Parts = append(mms.Parts, *part)
		return nil
	})
	if err != nil {
		return err
	}

	// Now look into making the body be a part, and set TextOnly appropriately
	mms.TextOnly = 0
	if body == nil || len(*body) == 0 {
		return nil
	}
	if len(mms.Parts) == 0 {
		mms.TextOnly = 1
	}
	mms.Size += len(*body)
	mms.Parts = append(mms.Parts, mmsPart{
		Order:           len(mms.Parts),
		ContentType:     "text/plain",
		Name:            "null",
		Charset:         utf8,
		ContentLocation: strOrNull{String: fmt.Sprintf("txt%06d.txt", id), Valid: true},
		Text:            strOrNull{String: *body, Valid: true},
	})
	return nil
}
func (w *writer) setupPhoneNumbers(primaryRecipientID int, threadRecipientID *int, members *string, mms *mms) error {
	err := w.findMMSPhoneNumsAndRecipientIDs(primaryRecipientID, threadRecipientID, members)
	if err != nil {
		return err
	}

	// Reset the Addrs slice to zero (keeping the underlying slice in place to avoid allocations)
	mms.Addrs = mms.Addrs[:0]
	// Build up mms.Addrs and mms.PhoneNumber. mms.PhoneNumber should become a ~ delimited list of phone numbers that aren't w.myPhoneNumber
	mms.PhoneNumber = ""
	foundMyNumber := false
	for i, recipientID := range w.mmsRecipientIDs {
		phoneNumber := w.mmsPhoneNumbers[i]
		// Since dialing and country codes exist, ignore them.
		// This might cause bugs, but I've seen data so messed up that I don't know how to fix them.
		// And this is a passion project :) if you come across this bug and happen to know what to do, have at it
		if strings.HasSuffix(phoneNumber, w.myPhoneNumber) {
			if mms.Type == mmsSendReq {
				mms.Addrs = append(mms.Addrs, mmsAddr{PhoneNumber: phoneNumber, Type: mmsAddrFrom, Charset: utf8})
			} else {
				mms.Addrs = append(mms.Addrs, mmsAddr{PhoneNumber: phoneNumber, Type: mmsAddrCC, Charset: utf8})
			}
			foundMyNumber = true
			continue
		}

		// This number is not w.myPhoneNumber, so it can be added to the list of phone numbers. At the end, the first ~ will be trimmed off
		mms.PhoneNumber += "~" + phoneNumber
		if mms.Type == mmsSendReq {
			mms.Addrs = append(mms.Addrs, mmsAddr{PhoneNumber: phoneNumber, Type: mmsAddrTo, Charset: utf8})
		} else if primaryRecipientID == recipientID {
			mms.Addrs = append(mms.Addrs, mmsAddr{PhoneNumber: phoneNumber, Type: mmsAddrFrom, Charset: utf8})
		} else {
			mms.Addrs = append(mms.Addrs, mmsAddr{PhoneNumber: phoneNumber, Type: mmsAddrCC, Charset: utf8})
		}
	}

	// mms.PhoneNumber has been added to with a ~ at the beginning each time, so time to remove the first one
	mms.PhoneNumber = strings.TrimPrefix(mms.PhoneNumber, "~")

	// Sometimes an MMS in signal doesn't include w.myPhoneNumber. Add it!
	if !foundMyNumber {
		if mms.Type == mmsSendReq {
			mms.Addrs = append(mms.Addrs, mmsAddr{PhoneNumber: w.myPhoneNumber, Type: mmsAddrFrom, Charset: utf8})
			return nil
		}
		mms.Addrs = append(mms.Addrs, mmsAddr{PhoneNumber: w.myPhoneNumber, Type: mmsAddrTo, Charset: utf8})
	}

	return nil
}

func (w *writer) findMMSPhoneNumsAndRecipientIDs(primaryRecipientID int, threadRecipientID *int, members *string) error {
	// Sometimes, there's no list of members to go through :'(
	if members == nil {
		if threadRecipientID == nil {
			w.mmsPhoneNumbers = []string{w.recipients[primaryRecipientID]}
			w.mmsRecipientIDs = []int{primaryRecipientID}
			return nil
		}
		w.mmsPhoneNumbers = []string{w.recipients[*threadRecipientID]}
		w.mmsRecipientIDs = []int{*threadRecipientID}
		return nil
	}

	// Reset these slices to zero (keeping the underlying slice in place to avoid allocations)
	w.mmsRecipientIDs = w.mmsRecipientIDs[:0]
	w.mmsPhoneNumbers = w.mmsPhoneNumbers[:0]
	recipientIDsStr := strings.Split(*members, ",")
	for i, recipientIDStr := range recipientIDsStr {
		parsed, err := strconv.ParseInt(recipientIDStr, 10, 64)
		if err != nil {
			return nil
		}
		w.mmsRecipientIDs = append(w.mmsRecipientIDs, int(parsed))
		phoneNumber, found := w.recipients[int(w.mmsRecipientIDs[i])]
		if !found {
			return fmt.Errorf("unknown recipient %v in members %v", recipientIDStr, members)
		}
		w.mmsPhoneNumbers = append(w.mmsPhoneNumbers, phoneNumber)
	}
	return nil
}

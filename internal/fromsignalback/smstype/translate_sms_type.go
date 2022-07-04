package smstype

// THIS FILE IS A MODIFIED VERSION OF THE ORIGINAL CODE FROM signal-back (WHICH IS DISTRIBUTED UNDER THE APACHE LICENSE, VERSION 2.0)

import (
	"fmt"
)

// SMSType is actually more like the status+origin of the message
type SMSType int64

const (
	// SMSDefaultType is a SMSType from https://developer.android.com/reference/android/provider/Telephony.TextBasedSmsColumns#MESSAGE_TYPE_ALL
	SMSDefaultType SMSType = iota
	// SMSInbox is a SMSType from https://developer.android.com/reference/android/provider/Telephony.TextBasedSmsColumns#MESSAGE_TYPE_INBOX
	SMSInbox
	// SMSSent is a SMSType from https://developer.android.com/reference/android/provider/Telephony.TextBasedSmsColumns#MESSAGE_TYPE_SENT
	SMSSent
	// SMSDraft is a SMSType from https://developer.android.com/reference/android/provider/Telephony.TextBasedSmsColumns#MESSAGE_TYPE_DRAFT
	SMSDraft
	// SMSOutbox is a SMSType from https://developer.android.com/reference/android/provider/Telephony.TextBasedSmsColumns#MESSAGE_TYPE_OUTBOX
	SMSOutbox
	// SMSFailed is a SMSType from https://developer.android.com/reference/android/provider/Telephony.TextBasedSmsColumns#MESSAGE_TYPE_FAILED
	SMSFailed
	// SMSQueued is a SMSType from https://developer.android.com/reference/android/provider/Telephony.TextBasedSmsColumns#MESSAGE_TYPE_QUEUED
	SMSQueued
)

// TranslateSMSType translates an SMS type from Signal's representation in its SQLite DB to an SMSType
func TranslateSMSType(t int64) (SMSType, error) {
	// Just get the lowest 5 bits, because everything else is masking:
	// https://github.com/signalapp/Signal-Android/blob/a4a4665aaa1a9c4ab1eb566cd249709a112aec69/app/src/main/java/org/thoughtcrime/securesms/database/MmsSmsColumns.java#L83-L90
	v := uint8(t) & 0x1F

	switch v {
	case 1, 20: // standard and signal standard
		return SMSInbox, nil
	case 2, 23: // standard and signal sent
		return SMSSent, nil
	case 3, 27: // standard and signal draft
		return SMSDraft, nil
	case 4, 21: // standard and signal outbox
		return SMSOutbox, nil
	case 5, 24: // standard and signal failed
		return SMSFailed, nil
	case 6, 22, 25, 26: // standard and signal queued, followed by signal fallbacks
		return SMSQueued, nil

	case 13:
		return SMSDefaultType, fmt.Errorf("found sms type BAD_DECRYPT_TYPE")

	default:
		return SMSDefaultType, fmt.Errorf("undefined SMS type: %#v\nplease report this issue, as well as (if possible) details about the SMS,\nsuch as whether it was sent, received, drafted, etc", t)
	}
}

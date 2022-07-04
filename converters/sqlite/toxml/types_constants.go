package toxml

import (
	"database/sql"
	"encoding/xml"
	"fmt"

	"github.com/jamescostian/signal-to-sms/internal/fromsignalback/smstype"
)

// sms represents a Short Message Service record.
// For the XSD, see https://synctech.com.au/wp-content/uploads/2018/01/sms.xsd_.txt
// For explanations of each field in that XSD, see https://developer.android.com/reference/android/provider/Telephony.TextBasedSmsColumns#constants_1
type sms struct {
	XMLName xml.Name `xml:"sms"`
	// Protocol identifier code. Seems to always be 0.
	Protocol int `xml:"protocol,attr"`
	// Never your phone number. Either the number you're sending to or receiving from
	PhoneNumber   string          `xml:"address,attr"`
	DateReceived  int64           `xml:"date,attr"`
	Type          smstype.SMSType `xml:"type,attr"`
	Subject       strOrNull       `xml:"subject,attr"`
	Body          *string         `xml:"body,attr"`
	ServiceCenter *string         `xml:"service_center,attr"`
	// A boolean stored as an int, because that's what SMS backup and restore expects
	Read     int8   `xml:"read,attr"`
	Status   int64  `xml:"status,attr"`
	DateSent *int64 `xml:"date_sent,attr"`
}

// signalSMSProtocol is the protocol Signal sometimes uses for its SMSes. It can safely be ignored and treated as the default, 0
const signalSMSProtocol = 31337

// mms represents a Multimedia Messaging Service record.
// For the XSD, see https://synctech.com.au/wp-content/uploads/2018/01/sms.xsd_.txt
// For explanations of each field in that XSD, see https://developer.android.com/reference/android/provider/Telephony.BaseMmsColumns#constants_1
type mms struct {
	XMLName xml.Name  `xml:"mms"`
	Parts   []mmsPart `xml:"parts>part"`
	Addrs   []mmsAddr `xml:"addrs>addr"`
	// Basically a boolean, except it's serialized and deserialized as a number, so it might as well be a number
	TextOnly       int8      `xml:"text_only,attr"`
	Subject        strOrNull `xml:"sub,attr"`
	RetrieveStatus strOrNull `xml:"retr_st,attr"`
	DateReceived   int64     `xml:"date,attr"`
	ContentClass   strOrNull `xml:"ct_cls,attr"`
	SubjectCharset ianaChset `xml:"sub_cs,attr"`
	// A boolean stored as an int, because that's what SMS backup and restore expects
	Read            int8          `xml:"read,attr"`
	ContentLocation strOrNull     `xml:"ct_l,attr"`
	TransactionID   strOrNull     `xml:"tr_id,attr"`
	Status          strOrNull     `xml:"st,attr"`
	MsgBox          mmsMessageBox `xml:"msg_box,attr"`
	// Never your phone number. Either the number you're sending to or receiving from
	PhoneNumber         string         `xml:"address,attr"`
	MessageClass        string         `xml:"m_cls,attr"`
	DeliveryTime        strOrNull      `xml:"d_tm,attr"`
	ReadStatus          strOrNull      `xml:"read_status,attr"`
	ContentType         mmsContentType `xml:"ct_t,attr"`
	RetrieveTextCharset strOrNull      `xml:"retr_txt_cs,attr"`
	DeliveryReport      mmsYesOrNo     `xml:"d_rpt,attr"`
	ID                  strOrNull      `xml:"m_id,attr"`
	DateSent            int64          `xml:"date_sent,attr"`
	// A boolean stored as an int, because that's what SMS backup and restore expects
	Seen int8    `xml:"seen,attr"`
	Type mmsType `xml:"m_type,attr"`
	// Always set to 18 in the dumps I've seen from SMS backup and restore
	MMSVersion    int         `xml:"v,attr"`
	Expiry        strOrNull   `xml:"exp,attr"`
	Priority      mmsPriority `xml:"pri,attr"`
	ReadReport    mmsYesOrNo  `xml:"rr,attr"`
	ResponseText  strOrNull   `xml:"resp_txt,attr"`
	ReportAllowed strOrNull   `xml:"rpt_a,attr"`
	// A boolean stored as an int, because that's what SMS backup and restore expects
	Locked         int8              `xml:"locked,attr"`
	RetrieveText   strOrNull         `xml:"retr_txt,attr"`
	ResponseStatus mmsResponseStatus `xml:"resp_st,attr"`
	Size           int               `xml:"m_size,attr"`
}

// NormalMMSVersion is the version SMS backup and restore uses. Not sure what's special about it, but it works.
const normalMMSVersion = 18

// mmsPart holds a data blob for an MMS.
// For the XSD, see https://synctech.com.au/wp-content/uploads/2018/01/sms.xsd_.txt
// For explanations of each field in that XSD, see https://developer.android.com/reference/android/provider/Telephony.Mms.Part#constants_1
type mmsPart struct {
	XMLName            xml.Name  `xml:"part"`
	Order              int       `xml:"seq,attr"`
	ContentType        string    `xml:"ct,attr"`
	Name               string    `xml:"name,attr"`
	Charset            ianaChset `xml:"chset,attr"`
	ContentDisposition strOrNull `xml:"cd,attr"`
	FileName           strOrNull `xml:"fn,attr"`
	ContentID          strOrNull `xml:"cid,attr"`
	ContentLocation    strOrNull `xml:"cl,attr"`
	ContentTypeStart   strOrNull `xml:"ctt_s,attr"`
	ContentTypeType    strOrNull `xml:"ctt_t,attr"`
	Text               strOrNull `xml:"text,attr"`
	Data               string    `xml:"data,attr,omitempty"`
}

// mmsAddr conveys a phone number that is sending/receiving an MMS.
// There is no XSD for this one, but when I made a back up with the app, it only included address, type, and charset attributes.
// For explanations of each of those fields, see https://developer.android.com/reference/android/provider/Telephony.Mms.Addr#constants_1
type mmsAddr struct {
	XMLName xml.Name `xml:"addr"`
	// Technically this should be called Address, but PhoneNumber is more obvious
	PhoneNumber string `xml:"address,attr"`
	// See mmsAddrType for possible values
	Type    mmsAddrType `xml:"type,attr"`
	Charset ianaChset   `xml:"charset,attr"`
}

// ianaChset is a number IANA assigned to a character set. See https://www.iana.org/assignments/character-sets/character-sets.xhtml
type ianaChset int

// utf8 is the number IANA assigned to the UTF-8 character set.
// Ideally I'd like to import it, but it's internal: https://go.googlesource.com/text/+/b1379a7b4714109485d36884eea575780d466cbd/encoding/internal/identifier/mib.go#729
const utf8 ianaChset = 106

var utf8Str = fmt.Sprintf("%v", utf8)

// mmsAddrType says whether an address is From, To, or CC
type mmsAddrType int

const (
	// mmsAddrCC is for CC (like mmsAddrTo, but it it's more like "let this person know as well" than "this is the main recipient").
	// It comes from https://android.googlesource.com/platform/frameworks/base/+/53ada2ab282c1b6b72365bc1c6b7aaa29e170eca/telephony/common/com/google/android/mms/pdu/PduHeaders.java#32
	mmsAddrCC mmsAddrType = 0x82
	// mmsAddrFrom is for the sender of a message.
	// It comes from https://android.googlesource.com/platform/frameworks/base/+/53ada2ab282c1b6b72365bc1c6b7aaa29e170eca/telephony/common/com/google/android/mms/pdu/PduHeaders.java#39
	mmsAddrFrom mmsAddrType = 0x89
	// mmsAddrTo is for the primary recipient of a message.
	// It comes from https://android.googlesource.com/platform/frameworks/base/+/53ada2ab282c1b6b72365bc1c6b7aaa29e170eca/telephony/common/com/google/android/mms/pdu/PduHeaders.java#55
	mmsAddrTo mmsAddrType = 0x97
)

// mmsMessageBox tells you which "box" the message is in (for example, there's a difference between a message being in the sent box vs the inbox)
type mmsMessageBox int

const (
	// mmsInbox is a possible value for mms.MsgBox from https://developer.android.com/reference/android/provider/Telephony.BaseMmsColumns#MESSAGE_BOX_INBOX
	mmsInbox mmsMessageBox = 1
	// mmsSentBox is a possible value for mms.MsgBox from https://developer.android.com/reference/android/provider/Telephony.BaseMmsColumns#MESSAGE_BOX_SENT
	mmsSentBox mmsMessageBox = 2
)

// mmsContentType is a content type for MMSes
type mmsContentType string

const (
	// wap is the Wireless Application Protocol content type.
	// It's a good value for mms.ContentType
	wap mmsContentType = "application/vnd.wap.multipart.related"
)

// mmsYesOrNo is for Delivery and Read reports
type mmsYesOrNo int

const (
	// mmsNo comes from https://android.googlesource.com/platform/frameworks/base/+/53ada2ab282c1b6b72365bc1c6b7aaa29e170eca/telephony/common/com/google/android/mms/pdu/PduHeaders.java#142
	mmsNo mmsYesOrNo = 0x81
)

// mmsResponseStatus lets you know if an operation was successful< or what error it had
type mmsResponseStatus int

const (
	// mmsResponseOk is a value for mms.ResponseStatus that comes from https://android.googlesource.com/platform/frameworks/base/+/53ada2ab282c1b6b72365bc1c6b7aaa29e170eca/telephony/common/com/google/android/mms/pdu/PduHeaders.java#214
	mmsResponseOk mmsResponseStatus = 0x80
)

// mmsPriority is a valid priority of MMSes
type mmsPriority int

const (
	// mmsPriorityNormal is a value for mms.Priority that comes from https://android.googlesource.com/platform/frameworks/base/+/53ada2ab282c1b6b72365bc1c6b7aaa29e170eca/telephony/common/com/google/android/mms/pdu/PduHeaders.java#208
	mmsPriorityNormal mmsPriority = 0x81
)

// mmsType is the type of an MMS
type mmsType int

const (
	// mmsSendReq is for MMSes that were sent.
	// It comes from https://android.googlesource.com/platform/frameworks/base/+/53ada2ab282c1b6b72365bc1c6b7aaa29e170eca/telephony/common/com/google/android/mms/pdu/PduHeaders.java#102
	mmsSendReq mmsType = 0x80
)

// strOrNull is database/sql.NullString, except it appears as "null" in XML if it's a null string
type strOrNull sql.NullString

// Scan calls Scan on a sql.NullString and puts the value into the string.
// This is meant to be used by database/sql.
func (s *strOrNull) Scan(val interface{}) error {
	ns := sql.NullString{}
	err := ns.Scan(val)
	if err != nil {
		return err
	}
	*s = strOrNull(ns)
	return nil
}

// Technically, you should never reuse a global variable. However, this is *specifically* for go's encoding/xml, so it'll be fine. And a bit faster.
var nullBytes = []byte("null")

// MarshalText allows encoding/xml to encode this as either "null" (if the string is nil) or the string value.
// DO NOT CALL THIS DIRECTLY! It reuses a buffer containing the word "null" for better performance when marshalling
func (s strOrNull) MarshalText() ([]byte, error) {
	if !s.Valid {
		return nullBytes, nil
	}
	return []byte(s.String), nil
}

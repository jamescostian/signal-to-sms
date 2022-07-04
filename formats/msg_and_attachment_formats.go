package formats

var BackupFrame = MsgFormat{IncludesAttachments: true, Name: "backup frames"}
var Encrypted = MsgFormat{
	IncludesAttachments: true,
	Name:                "encrypted",
	OpenForReads:        OpenFileForReads,
	OpenForWrites:       OpenFileForWrites,
	RequiredCtxKeys:     []interface{}{SignalBackupDBPasswordCtxKey},
}
var SQLiteMsgAndAttachment = MsgFormat{
	IncludesAttachments: true,
	Name:                "sqlite (with attachments)",
	OpenForReads:        SQLiteMsg.OpenForReads,
	OpenForWrites:       SQLiteMsg.OpenForWrites,
}
var XML = MsgFormat{
	IncludesAttachments: true,
	Name:                "xml",
	OpenForWrites:       OpenFileForWrites,
	RequiredCtxKeys:     []interface{}{MyPhoneNumberCtxKey},
}

var MsgAndAttachmentFormats = map[string]MsgFormat{
	"frames":    BackupFrame,
	"encrypted": Encrypted,
	"sqlite":    SQLiteMsgAndAttachment,
	"xml":       XML,
}

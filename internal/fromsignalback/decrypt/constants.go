package decrypt

// THIS FILE IS A MODIFIED VERSION OF THE ORIGINAL CODE FROM signal-back (WHICH IS LICENSED UNDER THE APACHE LICENSE, VERSION 2.0, AND WHICH IS DERIVATIVE OF Signal-Android, WHICH IS LICENSED UNDER THE GPLv3)

const (
	// BlobPieceSize is a specific size used by Signal for Android, and if you don't match it, your mac will be messed up.
	// Source: https://github.com/signalapp/Signal-Android/blob/7fc9876b1ea88ff1893ef627d1ecf0d2c2913508/app/src/main/java/org/thoughtcrime/securesms/backup/FullBackupImporter.java#L368
	BlobPieceSize = 8192

	// DigestIterations is how many times to perform a SHA512 digest on a user's password. Must be the same as the number Signal for Android uses.
	// Source: https://github.com/signalapp/Signal-Android/blob/8d4419705bbb3ab6b2bec5c85a9e3de806723dd6/app/src/main/java/org/thoughtcrime/securesms/backup/FullBackupBase.java#L15
	DigestIterations = 250000

	// MACSize is the length every MAC is truncated down to. Must be the same size as in Signal for Android.
	// Source: https://github.com/signalapp/Signal-Android/blob/4e01336b2f9aff2e727046132bac3d582e7f6e0c/app/src/main/java/org/thoughtcrime/securesms/backup/FullBackupExporter.java#L683
	MACSize = 10
)

// HKDFInfo is used to build an HKDF from which the cipher key and MAC key can be extracted - must be the same as Signal for Android.
// Source: https://github.com/signalapp/Signal-Android/blob/4e01336b2f9aff2e727046132bac3d582e7f6e0c/app/src/main/java/org/thoughtcrime/securesms/backup/FullBackupExporter.java#L546
var HKDFInfo = []byte("Backup Export")

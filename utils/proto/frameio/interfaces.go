// Package frameio contains tools to work with Signal for Android's BackupFrames
package frameio

import (
	"io"

	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
)

//go:generate mockery --name=Frame.+ --inpackage

// FrameReader allows reading *signal.BackupFrames, similar to io.Reader
type FrameReader interface {
	// Read puts a *signal.BackupFrame into the first argument you provide, and puts any blob read (e.g. an attachment) into the second argument you provide.
	// This allows you to reuse things, e.g. reuse a buffer for attachments.
	// Should return io.EOF on last frame.
	Read(backupFrameDest *signal.BackupFrame, blobDest *[]byte) error
}

// FrameReadCloser = FrameReader + io.Closer
type FrameReadCloser interface {
	FrameReader
	io.Closer
}

// FrameWriter allows writing *signal.BackupFrames, similar to io.Writer
type FrameWriter interface {
	// Writes out a BackupFrame, and optionally an attachment as well
	Write(backupFrame *signal.BackupFrame, blob []byte) error
}

// FrameWriteCloser = FrameWriter + io.Closer
type FrameWriteCloser interface {
	FrameWriter
	io.Closer
}

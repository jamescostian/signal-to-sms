package frameio

import (
	"context"
	"io"

	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
)

// CopyFrames copies all frames from a FrameReader to a FrameWriter. It does not close either.
func CopyFrames(ctx context.Context, writer FrameWriter, reader FrameReader) error {
	var (
		frame      signal.BackupFrame
		attachment []byte
	)
	var err error
	for {
		if err = ctx.Err(); err != nil {
			return err
		}
		if err = reader.Read(&frame, &attachment); err != nil {
			if err == io.EOF {
				return nil // No error! Finished writing!
			}
			return err
		}
		if err = writer.Write(&frame, attachment); err != nil {
			return err
		}
	}
}

package frameio

import (
	"io"

	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"google.golang.org/protobuf/proto"
)

func NewMultiWriter(writers ...FrameWriter) FrameWriter {
	return &multiWriter{writers}
}

type multiWriter struct {
	Writers []FrameWriter
}

func (w *multiWriter) Write(backupFrame *signal.BackupFrame, blob []byte) (err error) {
	for _, writer := range w.Writers {
		if err = writer.Write(backupFrame, blob); err != nil {
			return err
		}
	}
	return nil
}

// NewMultiReader creates a FrameReader that can read from several FrameReaders, sequentially.
// The first frame reader provided should provide the Header, and so on.
// If an End frame is encountered, it will be skipped.
// After calling Read enough times to go through all of the frames, an End frame will be provided via Read.
func NewMultiReader(readers ...FrameReader) FrameReader {
	return &multiReader{readers, false}
}

type multiReader struct {
	Readers             []FrameReader
	HasProvidedEndFrame bool
}

func (m *multiReader) Read(backupFrameDest *signal.BackupFrame, blobDest *[]byte) error {
	// If there are no readers left to read from, provide an end frame. After providing an end frame, return io.EOF
	if m.HasProvidedEndFrame {
		return io.EOF
	}
	if len(m.Readers) == 0 {
		m.HasProvidedEndFrame = true
		backupFrameDest.Reset()
		backupFrameDest.End = proto.Bool(true)
		return nil
	}

	err := m.Readers[0].Read(backupFrameDest, blobDest)
	if err == io.EOF {
		// The current reader ran out of content. Try the next one!
		// And the one after that if it fails, and keep going until all readers have been tried or one of them has data
		m.Readers = m.Readers[1:]
		return m.Read(backupFrameDest, blobDest)
	}
	if err != nil {
		return err
	}

	// Don't show an end frame until all readers have been exhausted. Skip to the next frame
	if backupFrameDest.End != nil {
		return m.Read(backupFrameDest, blobDest)
	}

	return nil
}

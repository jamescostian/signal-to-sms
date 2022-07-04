package decrypt

// THIS FILE IS A MODIFIED VERSION OF THE ORIGINAL CODE FROM signal-back (WHICH IS LICENSED UNDER THE APACHE LICENSE, VERSION 2.0, AND WHICH IS DERIVATIVE OF Signal-Android, WHICH IS LICENSED UNDER THE GPLv3)

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/pkg/errors"
)

// Read reads a frame of the backup. Not thread-safe, reads from the io.Reader originally provided
func (bf *BackupFile) Read(backupFrame *signal.BackupFrame, blob *[]byte) error {
	if bf.finished {
		return io.EOF
	}
	_, err := io.ReadFull(bf.file, bf.frameLen)
	if err != nil {
		return errors.Wrap(err, "reading frame length")
	}

	frameLength := binary.BigEndian.Uint32(bf.frameLen)
	frame := make([]byte, frameLength)

	if _, err := io.ReadFull(bf.file, frame); err != nil {
		return errors.Wrap(err, "reading frame")
	}
	// Split the frame MAC out
	frameLength -= MACSize
	theirMac := frame[frameLength:]
	frame = frame[:frameLength]

	bf.mac.Reset()
	bf.mac.Write(frame)
	ourMac := bf.mac.Sum(nil)

	if !bytes.Equal(theirMac, ourMac[:MACSize]) {
		return errors.New("Bad MAC in backup frame")
	}

	stream, err := NextAESCTRCipher(bf.cipherKey, bf.iv)
	if err != nil {
		return err
	}

	if uint32(cap(bf.decryptedFrame)) < frameLength {
		bf.decryptedFrame = make([]byte, frameLength*2)
	}
	bf.decryptedFrame = bf.decryptedFrame[:frameLength]
	stream.XORKeyStream(bf.decryptedFrame, frame)

	if backupFrame == nil {
		backupFrame = new(signal.BackupFrame)
	}
	if err = bf.encoder.Unmarshal(bf.decryptedFrame, backupFrame); err != nil {
		return errors.Wrap(err, "unmarshalling protobuf")
	}
	blobSize := frameio.BlobSizeOf(backupFrame)
	if blobSize > 0 {
		return errors.Wrap(bf.decryptBlob(blobSize, blob), "decrypting blob")
	}
	if backupFrame.End != nil && *backupFrame.End {
		bf.finished = true
	}
	return nil
}

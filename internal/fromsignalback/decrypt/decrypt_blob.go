package decrypt

// THIS FILE IS A MODIFIED VERSION OF THE ORIGINAL CODE FROM signal-back (WHICH IS LICENSED UNDER THE APACHE LICENSE, VERSION 2.0, AND WHICH IS DERIVATIVE OF Signal-Android, WHICH IS LICENSED UNDER THE GPLv3)

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
)

// decryptBlob decrypts the next blob (e.g. an attachment) in the file. Not thread-safe.
func (bf *BackupFile) decryptBlob(length uint32, output *[]byte) error {
	stream, err := NextAESCTRCipher(bf.cipherKey, bf.iv)
	if err != nil {
		return err
	}
	bf.mac.Write(bf.iv)

	// Grow output if needed. Resize it down to the correct length afterwards
	if output == nil || cap(*output) <= int(length) {
		*output = make([]byte, length)
	}
	*output = (*output)[:length]

	currentOutputOffset := 0
	bf.decryptBlobPiece = bf.decryptBlobPiece[:cap(bf.decryptBlobPiece)]
	for length > 0 {
		// If the scratch buffer is too big, only use a part of it.
		if int(length) < BlobPieceSize {
			bf.decryptBlobPiece = bf.decryptBlobPiece[:length]
		}
		n, err := bf.file.Read(bf.decryptBlobPiece)
		if err != nil {
			return errors.Wrap(err, "failed to read att")
		}
		bf.mac.Write(bf.decryptBlobPiece)

		stream.XORKeyStream((*output)[currentOutputOffset:], bf.decryptBlobPiece)
		currentOutputOffset = currentOutputOffset + len(bf.decryptBlobPiece)
		length -= uint32(n)
	}

	if _, err := io.ReadFull(bf.file, bf.decryptBlobTheirMac); err != nil {
		return err
	}
	ourMac := bf.mac.Sum(nil)

	if bytes.Equal(bf.decryptBlobTheirMac, ourMac) {
		return errors.New("Bad MAC in attachment")
	}

	return nil
}

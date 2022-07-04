package decrypt

// THIS FILE IS A MODIFIED VERSION OF THE ORIGINAL CODE FROM signal-back (WHICH IS LICENSED UNDER THE APACHE LICENSE, VERSION 2.0, AND WHICH IS DERIVATIVE OF Signal-Android, WHICH IS LICENSED UNDER THE GPLv3)

import (
	"crypto"
	"crypto/hmac"
	"encoding/binary"

	// Fix "crypto: requested hash function #5 is unavailable"
	_ "crypto/sha256"
	// Fix "crypto: requested hash function #7 is unavailable"
	_ "crypto/sha512"

	"hash"
	"io"

	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/pkg/errors"
)

// Chosen based on a few trials to see what initial number worked out best. There's nothing special about this specific number.
const decryptedFrameDefaultSize = 1 << 12

type BackupFile struct {
	file      io.Reader
	cipherKey []byte
	macKey    []byte
	mac       hash.Hash
	iv        []byte
	encoder   frameio.Encoder
	finished  bool

	decryptBlobTheirMac []byte
	decryptBlobPiece    []byte

	decryptedFrame []byte
	frameLen       []byte
}

// NewBackupFile allows reading decrypted backup frames from encrypted data in an io.Reader. Not thread-safe.
func NewBackupFile(file io.Reader, password string, encoder frameio.Encoder) (*BackupFile, error) {
	headerLengthBytes := make([]byte, 4)
	_, err := io.ReadFull(file, headerLengthBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read headerLengthBytes")
	}
	headerLength := binary.BigEndian.Uint32(headerLengthBytes)

	headerFrame := make([]byte, headerLength)
	_, err = io.ReadFull(file, headerFrame)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read headerFrame")
	}
	frame := &signal.BackupFrame{}
	if err := encoder.Unmarshal(headerFrame, frame); err != nil {
		return nil, errors.Wrap(err, "failed to decode header")
	}

	iv := frame.Header.Iv
	if len(iv) != 16 {
		return nil, errors.New("The IV in the header isn't 16 bytes long")
	}

	cipherKey, macKey, err := DeriveSecrets(password, frame.Header.Salt)
	if err != nil {
		return nil, err
	}

	// Each iteration of decryption requires incrementing the IV. Allow this to work transparently by first decrementing once
	negative1 := -1
	AddToUint32Bytes(iv, uint32(negative1))

	return &BackupFile{
		file:      file,
		cipherKey: cipherKey,
		macKey:    macKey,
		mac:       hmac.New(crypto.SHA256.New, macKey),
		iv:        iv,
		encoder:   encoder,

		decryptBlobTheirMac: make([]byte, MACSize),
		decryptBlobPiece:    make([]byte, BlobPieceSize),

		decryptedFrame: make([]byte, decryptedFrameDefaultSize),
		frameLen:       headerLengthBytes,
	}, nil
}

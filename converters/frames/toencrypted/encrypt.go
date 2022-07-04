package toencrypted

import (
	"bytes"
	"context"
	"crypto"
	"crypto/hmac"
	"encoding/binary"
	"io"

	// Fix "crypto: requested hash function #5 is unavailable"
	_ "crypto/sha256"
	// Fix "crypto: requested hash function #7 is unavailable"
	_ "crypto/sha512"

	"hash"
	"os"

	"github.com/jamescostian/signal-to-sms/internal/fromsignalback/decrypt"
	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/protobuf"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
)

func Encrypt(ctx context.Context, input frameio.FrameReader, output *os.File, password string) error {
	enc, err := newEncryptor(input, output, password, protobuf.ProtobufEncoder)
	if err != nil {
		return err
	}
	for err == nil {
		if err = ctx.Err(); err != nil {
			return err
		}
		err = enc.EncryptNextFrame()
	}
	if err == io.EOF {
		return nil
	}
	return err
}

type encryptor struct {
	input     frameio.FrameReader
	output    *os.File
	cipherKey []byte
	macKey    []byte
	mac       hash.Hash
	iv        []byte
	encoder   frameio.Encoder

	// The blob that is currently being encrypted will be encrypted piece by piece in here
	encryptBlobPiece []byte

	encryptedFrame   *bytes.Buffer
	frameLenBytes    []byte
	frame            signal.BackupFrame
	frameBytes       *[]byte
	blobForCurrFrame []byte
}

// NewBackupFile allows reading decrypted backup frames from encrypted data in an io.Reader. Not thread-safe.
func newEncryptor(input frameio.FrameReader, output *os.File, password string, encoder frameio.Encoder) (*encryptor, error) {
	salt := make([]byte, 32)
	iv := make([]byte, 16)
	if err := fillRandBytes(salt); err != nil {
		return nil, err
	}
	if err := fillRandBytes(iv); err != nil {
		return nil, err
	}
	cipherKey, macKey, err := decrypt.DeriveSecrets(password, salt)
	if err != nil {
		return nil, err
	}
	headerFrame, err := encoder.Marshal(&signal.BackupFrame{Header: &signal.Header{Salt: salt, Iv: iv}})
	if err != nil {
		return nil, err
	}
	headerFrameLen := make([]byte, 4)
	binary.BigEndian.PutUint32(headerFrameLen, uint32(len(headerFrame)))
	if _, err = output.Write(headerFrameLen); err != nil {
		return nil, err
	}
	if _, err = output.Write(headerFrame); err != nil {
		return nil, err
	}

	// Each iteration of decryption requires incrementing the IV. Allow this to work transparently by first decrementing once
	negative1 := -1
	decrypt.AddToUint32Bytes(iv, uint32(negative1))

	return &encryptor{
		input:     input,
		output:    output,
		cipherKey: cipherKey,
		macKey:    macKey,
		mac:       hmac.New(crypto.SHA256.New, macKey),
		iv:        iv,
		encoder:   encoder,

		encryptBlobPiece: make([]byte, decrypt.BlobPieceSize),

		encryptedFrame:   bytes.NewBuffer(nil),
		frameLenBytes:    headerFrameLen,
		frame:            signal.BackupFrame{},
		frameBytes:       &[]byte{},
		blobForCurrFrame: make([]byte, 0),
	}, nil
}

func (e *encryptor) EncryptNextFrame() error {
	if err := e.input.Read(&e.frame, &e.blobForCurrFrame); err != nil {
		return err
	}
	if err := e.encryptFrame(&e.frame); err != nil {
		return err
	}
	if len(e.blobForCurrFrame) == 0 {
		return nil
	}
	return e.encryptBlob(e.blobForCurrFrame)
}

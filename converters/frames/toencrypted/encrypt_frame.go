package toencrypted

import (
	"encoding/binary"

	"github.com/jamescostian/signal-to-sms/internal/fromsignalback/decrypt"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
)

func (e *encryptor) encryptFrame(frame *signal.BackupFrame) error {
	cipher, err := decrypt.NextAESCTRCipher(e.cipherKey, e.iv)
	if err != nil {
		return err
	}
	if *e.frameBytes, err = e.encoder.Marshal(frame); err != nil {
		return err
	}

	frameLen := len(*e.frameBytes)
	binary.BigEndian.PutUint32(e.frameLenBytes, uint32(frameLen+decrypt.MACSize))
	if _, err = e.output.Write(e.frameLenBytes); err != nil {
		return err
	}

	e.encryptedFrame.Grow(frameLen)
	encryptedFrameBytes := e.encryptedFrame.Bytes()[:frameLen]
	cipher.XORKeyStream(encryptedFrameBytes, *e.frameBytes)
	if _, err = e.output.Write(encryptedFrameBytes); err != nil {
		return err
	}

	e.mac.Reset()
	e.mac.Write(encryptedFrameBytes)
	_, err = e.output.Write(e.mac.Sum(nil)[:decrypt.MACSize])
	return err
}

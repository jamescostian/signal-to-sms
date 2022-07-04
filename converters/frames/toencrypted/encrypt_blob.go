package toencrypted

import "github.com/jamescostian/signal-to-sms/internal/fromsignalback/decrypt"

func (e *encryptor) encryptBlob(blob []byte) error {
	cipher, err := decrypt.NextAESCTRCipher(e.cipherKey, e.iv)
	if err != nil {
		return err
	}
	e.mac.Write(e.iv)
	for len(e.blobForCurrFrame) > 0 {
		pieceSize := intMin(decrypt.BlobPieceSize, len(e.blobForCurrFrame))
		e.encryptBlobPiece = e.encryptBlobPiece[:pieceSize]
		cipher.XORKeyStream(e.encryptBlobPiece, e.blobForCurrFrame[:pieceSize])
		e.blobForCurrFrame = e.blobForCurrFrame[pieceSize:]
		if _, err = e.output.Write(e.encryptBlobPiece); err != nil {
			return err
		}
		e.mac.Write(e.encryptBlobPiece)
	}
	_, err = e.output.Write(e.mac.Sum(nil)[:decrypt.MACSize])
	return err
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

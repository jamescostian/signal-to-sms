package decrypt

// THIS FILE IS A MODIFIED VERSION OF THE ORIGINAL CODE FROM signal-back (WHICH IS LICENSED UNDER THE APACHE LICENSE, VERSION 2.0, AND WHICH IS DERIVATIVE OF Signal-Android, WHICH IS LICENSED UNDER THE GPLv3)

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"io"
	"strings"

	"golang.org/x/crypto/hkdf"
)

func AddToUint32Bytes(b []byte, val uint32) {
	num := binary.BigEndian.Uint32(b)
	binary.BigEndian.PutUint32(b, num+val)
}

func deriveKey(password string, salt []byte) []byte {
	digest := crypto.SHA512.New()
	input := []byte(strings.Replace(strings.TrimSpace(password), " ", "", -1))
	hash := input

	if salt != nil {
		digest.Write(salt)
	}

	for i := 0; i < DigestIterations; i++ {
		digest.Write(hash)
		digest.Write(input)
		hash = digest.Sum(nil)
		digest.Reset()
	}

	return hash[:32]
}

func DeriveSecrets(password string, salt []byte) (cipherKey []byte, macKey []byte, err error) {
	sha := crypto.SHA256.New
	okm := make([]byte, 64)

	hkdf := hkdf.New(sha, deriveKey(password, salt), make([]byte, sha().Size()), HKDFInfo)
	_, err = io.ReadFull(hkdf, okm)
	cipherKey = okm[:32]
	macKey = okm[32:]
	return
}

func NextAESCTRCipher(cipherKey []byte, iv []byte) (cipher.Stream, error) {
	// Increment the IV for each iteration just like signal does before creating a new cipher stream based on this new IV
	AddToUint32Bytes(iv, 1)

	aesCipher, err := aes.NewCipher(cipherKey)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesCipher, iv)

	return stream, nil
}

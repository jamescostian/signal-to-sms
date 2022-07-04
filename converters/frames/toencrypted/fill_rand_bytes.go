package toencrypted

import "crypto/rand"

var fillRandBytes = func(bytes []byte) error {
	_, err := rand.Read(bytes)
	return err
}

func fillRandBytesDeterministically(bytes []byte) error {
	for i := range bytes {
		bytes[i] = 42 // I call this "v2" since the original version was missing the 2: https://xkcd.com/221/
	}
	return nil
}

//go:build deterministic_output
// +build deterministic_output

package toencrypted

func init() {
	fillRandBytes = fillRandBytesDeterministically
}

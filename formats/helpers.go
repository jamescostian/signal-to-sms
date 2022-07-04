package formats

import (
	"io"
	"os"
)

func OpenFileForReads(path string) (io.Closer, error) {
	return os.Open(path)
}

func OpenFileForWrites(path string, flag int, perm os.FileMode) (io.Closer, error) {
	return os.OpenFile(path, flag, perm)
}

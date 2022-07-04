package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

// Create makes a SQLite DB and opens it up
func Create(dbPath string, flag int, perm os.FileMode) (*sql.DB, error) {
	file, err := os.OpenFile(dbPath, flag, perm)
	if err != nil {
		return nil, err
	}
	file.Close()
	return Open(dbPath)
}

// CreateTemp can be used to create a temporary SQLite DB. It gives you a cleanup function whose execution you can defer.
// The cleanup function will panic if there's an error deleting the file, since databases can contain sensitive info.
func CreateTemp() (db *sql.DB, deleteOrPanic func(), err error) {
	db, _, deleteOrPanic, err = CreateTempCopyOf("")
	return
}

// Don't let others read/write these temporary files, they could contain sensitive info like text messages and attachments
var tempFileMode fs.FileMode = 0600

func CreateTempCopyOf(basedOn string) (db *sql.DB, path string, deleteOrPanic func(), err error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return
	}
	if basedOn != "" {
		var bytes []byte
		bytes, err = ioutil.ReadFile(basedOn)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(f.Name(), bytes, tempFileMode)
		if err != nil {
			return
		}
	}
	db, err = Create(f.Name(), os.O_WRONLY|os.O_CREATE, tempFileMode)
	if err != nil {
		return
	}
	deleteOrPanic = func() {
		// Try to close. Even if it fails, still try to delete the file
		db.Close()
		if err := os.Remove(f.Name()); err != nil && !errors.Is(err, fs.ErrNotExist) {
			panic(fmt.Errorf("unable to delete this temporary file that contains your signal database (you should delete it yourself): %v", f.Name()))
		}
	}
	return
}

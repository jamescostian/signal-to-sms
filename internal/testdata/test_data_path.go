package testdata

import (
	"path"
	"runtime"
)

var TestDataPassword = "111111111111111111111111111111"

var TestDataPath string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	TestDataPath = path.Dir(filename)
}

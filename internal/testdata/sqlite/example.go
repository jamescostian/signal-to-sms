package sqlite

import (
	"os"
	"path"

	"github.com/jamescostian/signal-to-sms/internal/testdata"
)

//go:generate go run -tags=deterministic_output ../../../. -I encrypted -i ../encrypted/example.backup -O sqlite -o example.sqlite -p 111111111111111111111111111111 -t

var ExamplePath = path.Join(testdata.TestDataPath, "sqlite", "example.sqlite")
var Example, _ = os.Open(ExamplePath)

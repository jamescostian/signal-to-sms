package encrypted

import (
	"os"
	"path"

	"github.com/jamescostian/signal-to-sms/internal/testdata"
)

//go:generate go run -tags=deterministic_output ../../../. -I prototext,sqlite -i ../prototext/example.prototext --attachment-input ../attachments/sqlite/example.sqlite -O encrypted -o example.backup -p 111111111111111111111111111111 -t

var ExamplePath = path.Join(testdata.TestDataPath, "encrypted", "example.backup")
var Example, _ = os.Open(ExamplePath)
var ExamplePassword = "111111111111111111111111111111"

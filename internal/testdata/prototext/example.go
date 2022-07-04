package prototext

import (
	"os"
	"path"

	"github.com/jamescostian/signal-to-sms/internal/testdata"
)

//go:generate go run -tags=deterministic_output ../../../. -I encrypted -i ../encrypted/example.backup -O prototext,sqlite -o example.prototext --attachment-output ../attachments/sqlite/example.sqlite -p 111111111111111111111111111111 -t

var ExamplePath = path.Join(testdata.TestDataPath, "prototext", "example.prototext")
var Example, _ = os.Open(ExamplePath)

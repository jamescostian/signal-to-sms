package xml

import (
	"os"
	"path"

	"github.com/jamescostian/signal-to-sms/internal/testdata"
)

//go:generate go run -tags=deterministic_output ../../../. -I sqlite,sqlite -i ../sqlite/example.sqlite --attachment-input ../attachments/sqlite/example.sqlite -O xml -o example.xml -m 6054756968 -t

var ExamplePath = path.Join(testdata.TestDataPath, "xml", "example.xml")
var Example, _ = os.Open(ExamplePath)
var ExampleIndent = "\t"
var ExamplePhone = "6054756968"

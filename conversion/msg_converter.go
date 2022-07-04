package conversion

import (
	"context"

	"github.com/jamescostian/signal-to-sms/formats"
)

type MsgConverter struct {
	InputFormat  formats.MsgFormat
	OutputFormat formats.MsgFormat
	// Convert between the formats, respecting if the context is canceled
	// Note: if your input format isn't intermediate, then input format should be a file path string. If your output format isn't intermediate, then you should save your output to the output path provided
	Convert func(ctx context.Context, c *Conversion) (intermediateOutput interface{}, err error)
}

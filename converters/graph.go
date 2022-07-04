package converters

import (
	"github.com/jamescostian/signal-to-sms/conversion"
	"github.com/jamescostian/signal-to-sms/formats"
)

func ComputeConversionGraph() (*conversion.Graph, error) {
	allFormats := make([]formats.MsgFormat, len(formats.MsgAndAttachmentFormats)+len(formats.MsgOnlyFormats))
	for _, format := range formats.MsgAndAttachmentFormats {
		allFormats = append(allFormats, format)
	}
	for _, format := range formats.MsgOnlyFormats {
		allFormats = append(allFormats, format)
	}
	return conversion.NewGraph(MsgConverter, allFormats)
}

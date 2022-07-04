package cmd

import (
	"github.com/jamescostian/signal-to-sms/conversion"
	"github.com/jamescostian/signal-to-sms/converters"
	"github.com/jamescostian/signal-to-sms/formats"
)

func findConversionPath(in formats.MsgFormat, out formats.MsgFormat) ([]conversion.MsgConverter, error) {
	conversionGraph, err := converters.ComputeConversionGraph()
	if err != nil {
		return nil, err
	}
	return conversionGraph.FindPath(in, out)
}

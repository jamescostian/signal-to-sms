package cmd

import (
	"context"
	"log"

	"github.com/jamescostian/signal-to-sms/conversion"
	"github.com/jamescostian/signal-to-sms/utils/attachments"
)

// Run Convert and print out the progress of the conversion as it happens
func consoleFriendlyConvert(ctx context.Context, conversionPath []conversion.MsgConverter, attachmentIn attachments.Store, attachmentOut attachments.Store, prgArgs prgArgsT) error {
	conversionProgress := make(chan *conversion.MsgConverter, len(conversionPath))
	go func() {
		for converter := range conversionProgress {
			log.Printf("Converting to %v...\n", converter.OutputFormat.Name)
		}
	}()
	return conversion.Convert(
		context.WithValue(ctx, conversion.ConversionProgressCtxKey, conversionProgress),
		conversionPath,
		prgArgs.MsgInputPath,
		prgArgs.MsgOutputPath,
		attachmentIn,
		attachmentOut,
		prgArgs.MsgFlags(),
		prgArgs.MsgPerms(),
	)
}

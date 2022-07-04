package cmd

import (
	"fmt"

	"github.com/jamescostian/signal-to-sms/formats"
	"github.com/jamescostian/signal-to-sms/utils/attachments"
)

func validateConversionFormats(msgInputFormat formats.MsgFormat, msgOutputFormat formats.MsgFormat, attachmentIn attachments.Store, attachmentOut attachments.Store) error {
	if attachmentIn == nil && !msgInputFormat.IncludesAttachments {
		return fmt.Errorf("when using %v, a separate file containing attachments must be provided", msgInputFormat.Name)
	}
	if attachmentOut == nil && !msgOutputFormat.IncludesAttachments {
		return fmt.Errorf("when using %v, a separate file to store attachments in must be provided", msgOutputFormat.Name)
	}

	if attachmentIn != nil && msgInputFormat.IncludesAttachments {
		return fmt.Errorf("%v includes attachments, so you cannot specify another set of attachments to use", msgInputFormat.Name)
	}
	if attachmentOut != nil && msgOutputFormat.IncludesAttachments {
		return fmt.Errorf("%v specified includes attachments, so you cannot specify a separate place to put attachments", msgOutputFormat.Name)
	}

	if msgInputFormat.OpenForReads == nil {
		return fmt.Errorf("cannot use %v as a messages input format", msgInputFormat.Name)
	}
	if msgOutputFormat.OpenForWrites == nil {
		return fmt.Errorf("cannot use %v as a message output format", msgOutputFormat.Name)
	}

	return nil
}

package cmd

import (
	"fmt"

	"github.com/jamescostian/signal-to-sms/formats"
	"github.com/jamescostian/signal-to-sms/utils/attachments"
)

func openAttachments(args *prgArgsT) (in attachments.Store, out attachments.Store, err error) {
	inFormat := args.AttachmentInputFormat()
	if inFormat != nil {
		attachmentFmt, found := formats.AttachmentFormats[*inFormat]
		if !found {
			return nil, nil, fmt.Errorf("no such attachment format: %v", *inFormat)
		}
		in, err = attachmentFmt.Open(args.AttachmentInputPath)
		if err != nil {
			return
		}
	}
	outFormat := args.AttachmentOututFormat()
	if outFormat != nil {
		attachmentFmt, found := formats.AttachmentFormats[*outFormat]
		if !found {
			if in != nil {
				in.Close()
			}
			return nil, nil, fmt.Errorf("no such attachment format: %v", *outFormat)
		}
		if err = attachmentFmt.Initialize(args.AttachmentOutputPath, args.AttachmentFlags(), args.AttachmentPerms()); err != nil {
			return
		}
		out, err = attachmentFmt.Open(args.AttachmentOutputPath)
	}
	return
}

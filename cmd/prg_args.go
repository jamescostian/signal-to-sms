package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jamescostian/signal-to-sms/conversion"
	"github.com/jamescostian/signal-to-sms/formats"
	"github.com/pkg/errors"
)

type prgArgsT struct {
	PhoneNumber string
	Password    string

	OverwriteExisting bool

	MsgInputPath         string
	MsgOutputPath        string
	AttachmentInputPath  string
	AttachmentOutputPath string

	// These two can be things that have separate attachments (e.g. "prototext,tar") or unified messages and attachments (e.g. "xml").
	// Use other functions that read from these, instead of reading from them directly.
	RawInputFormat  string
	RawOutputFormat string
}

var formatDelimeter = ","

func (args *prgArgsT) MsgInputFormat() (format formats.MsgFormat, found bool) {
	msgFmt := firstPart(args.RawInputFormat, formatDelimeter)
	if msgFmt == args.RawInputFormat {
		format, found = formats.MsgAndAttachmentFormats[msgFmt]
	} else {
		format, found = formats.MsgOnlyFormats[msgFmt]
	}
	return
}
func (args *prgArgsT) MsgOutputFormat() (format formats.MsgFormat, found bool) {
	msgFmt := firstPart(args.RawOutputFormat, formatDelimeter)
	if msgFmt == args.RawOutputFormat {
		format, found = formats.MsgAndAttachmentFormats[msgFmt]
	} else {
		format, found = formats.MsgOnlyFormats[msgFmt]
	}
	return
}
func (args *prgArgsT) AttachmentInputFormat() *string {
	return secondPart(args.RawInputFormat, formatDelimeter)
}
func (args *prgArgsT) AttachmentOututFormat() *string {
	return secondPart(args.RawOutputFormat, formatDelimeter)
}

func firstPart(str string, delimeter string) string {
	delimIdx := strings.Index(str, delimeter)
	if delimIdx == -1 {
		return str
	}
	return str[:delimIdx]
}

func secondPart(str string, delimeter string) *string {
	delimIdx := strings.Index(str, delimeter)
	if delimIdx == -1 {
		return nil
	}
	part := str[delimIdx+1:]
	return &part
}

func numbersOnly(input string) string {
	output := strings.Builder{}
	for _, char := range input {
		if char >= '0' && char <= '9' {
			output.WriteRune(char)
		}
	}
	return output.String()
}

func (args *prgArgsT) AskForMissingArgs(ctx context.Context, conversionPath []conversion.MsgConverter) (context.Context, error) {
	ctx = context.WithValue(ctx, formats.SignalBackupDBPasswordCtxKey, args.Password)
	ctx = context.WithValue(ctx, formats.MyPhoneNumberCtxKey, args.PhoneNumber)
	ctxKeysNeeded := findAllCtxKeysNeeded(ctx, conversionPath)
	var (
		val interface{}
		err error
	)
	for key := range ctxKeysNeeded {
		val, err = ValidateAndFillCtxFromUser(ctx, key)
		if err != nil {
			return ctx, errors.Wrapf(err, "error getting a context key (%v)", key)
		}
		ctx = context.WithValue(ctx, key, val)
	}
	return ctx, nil
}

// ValidateAndFillCtxFromUser is a variable that can be overwritten, and that will be called with every context key that's required for every MsgFormat in the conversion path.
// Even if a context key is present in the context, this function will still be called, allowing you to perform validation here.
// This function can ask a user to enter in some data needed to proceed with a conversion, or throw errors
var ValidateAndFillCtxFromUser = func(ctx context.Context, ctxKey interface{}) (value interface{}, err error) {
	switch ctxKey {
	case formats.SignalBackupDBPasswordCtxKey:
		// TODO: do better validation here of the password here
		if val := ctx.Value(ctxKey); val != nil && val.(string) != "" {
			return val, nil
		}
		// Not hidden like a normal password, shown in the clear, because:
		// 1. This password isn't something a user chooses, so giving it away gives away nothing but the user's password
		// 2. This password isn't useful without the encrypted backup
		// 3. Under what circumstances does your attacker have your back up and have visibility of your screen, but can't see what you're physically typing via a keylogger or their eyes?
		var password string
		err := survey.AskOne(&survey.Input{
			Message: "The Signal backup password (a bunch of numbers)",
			Help: `This is not the PIN code that Signal makes you type in every once in a while.
				This is the password that they tell you once when you turn on back ups.
				It's only numbers, and is pretty long. You may have screenshotted it.
				If you lost it, you'll need to disable backups on your Android phone, then turn them back on and grab the new backup file.`,
		}, &password, survey.WithValidator(survey.Required))
		return numbersOnly(password), err
	case formats.MyPhoneNumberCtxKey:
		// TODO: do better validation here of the phone number here
		if val := ctx.Value(ctxKey); val != nil && val.(string) != "" {
			return val, nil
		}
		var phoneNumber string
		err := survey.AskOne(&survey.Input{
			Message: "Your phone number, without the country code",
		}, &phoneNumber, survey.WithValidator(survey.Required))
		return numbersOnly(phoneNumber), err
	default:
		return nil, fmt.Errorf("unknown context key (%v) is required. Consider setting cmd.ValidateAndFillCtxFromUser to be your own function that calls the original function, and adds support for your custom context keys", ctxKey)
	}
}

func findAllCtxKeysNeeded(ctx context.Context, conversionPath []conversion.MsgConverter) map[interface{}]bool {
	ctxKeyMap := map[interface{}]bool{}
	for _, converter := range conversionPath {
		for _, key := range converter.InputFormat.RequiredCtxKeys {
			ctxKeyMap[key] = true
		}
		for _, key := range converter.OutputFormat.RequiredCtxKeys {
			ctxKeyMap[key] = true
		}
	}
	return ctxKeyMap
}

func (args *prgArgsT) AttachmentPerms() os.FileMode {
	return 0600
}
func (args *prgArgsT) AttachmentFlags() int {
	flags := os.O_CREATE | os.O_RDWR
	if args.OverwriteExisting {
		return flags | os.O_TRUNC
	}
	return flags
}
func (args *prgArgsT) MsgPerms() os.FileMode {
	return 0600
}
func (args *prgArgsT) MsgFlags() int {
	flags := os.O_CREATE | os.O_WRONLY
	if args.OverwriteExisting {
		return flags | os.O_TRUNC
	}
	return flags
}

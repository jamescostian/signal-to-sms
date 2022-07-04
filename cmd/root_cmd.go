package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var prgArgs prgArgsT

var rootCmd = &cobra.Command{
	Use:           "signal-to-sms [path to signal export file] [path to save SMSes to]",
	Short:         "Convert Signal for Android encrypted backups to SMSes",
	Long:          "Convert Signal for Android encrypted backups to XML that can be read by SMS Backup & Restore. Combined with other tools, signal-to-sms allows moving Signal messages from Android to iOS.\n\nAlso allows converting between other formats, and can handle using 2 separate files, one for messages and one for attachments. Example formats include encrypted (the default input format, which is what Signal for Android provides), xml (the default output format, which can be imported into SMS Backup & Restore), and prototext,tar (provides a way to inspect the raw, but decrypted data Signal for Android provided).",
	Version:       "unknown",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	// If there's an error, that shouldn't cause the usage menu to pop-up; use --help if you want the usage menu.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
	},
	// Here's the actual code to run everything
	RunE: func(cmd *cobra.Command, args []string) error {
		attachmentIn, attachmentOut, err := openAttachments(&prgArgs)
		if err != nil {
			return err
		}
		msgInputFormat, found := prgArgs.MsgInputFormat()
		if !found {
			return fmt.Errorf("unknown message input format")
		}
		msgOutputFormat, found := prgArgs.MsgOutputFormat()
		if !found {
			return fmt.Errorf("unknown message input format")
		}
		err = validateConversionFormats(msgInputFormat, msgOutputFormat, attachmentIn, attachmentOut)
		if err != nil {
			return err
		}
		conversionPath, err := findConversionPath(msgInputFormat, msgOutputFormat)
		if err != nil {
			return err
		}
		ctx, err := prgArgs.AskForMissingArgs(cmd.Context(), conversionPath)
		if err != nil {
			return err
		}
		return consoleFriendlyConvert(
			ctx,
			conversionPath,
			attachmentIn,
			attachmentOut,
			prgArgs,
		)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&prgArgs.RawInputFormat, "input-format", "I", "encrypted", "The format the input(s) come(s) in, like encrypted or prototext,tar")
	rootCmd.Flags().StringVarP(&prgArgs.RawOutputFormat, "output-format", "O", "xml", "The format the output(s) is/are desired in, like xml or prototext,tar")

	rootCmd.Flags().StringVarP(&prgArgs.MsgInputPath, "msg-input", "i", "", "The path to the input containing messages")
	if err := rootCmd.MarkFlagRequired("msg-input"); err != nil {
		panic(err)
	}
	rootCmd.Flags().StringVarP(&prgArgs.MsgOutputPath, "msg-output", "o", "", "The path to the output for where to store messages")
	if err := rootCmd.MarkFlagRequired("msg-output"); err != nil {
		panic(err)
	}

	rootCmd.Flags().StringVar(&prgArgs.AttachmentInputPath, "attachment-input", "", "The path to the input containing messages")
	rootCmd.Flags().StringVar(&prgArgs.AttachmentOutputPath, "attachment-output", "", "The path to the output for where to store messages")

	rootCmd.Flags().BoolVarP(&prgArgs.OverwriteExisting, "overwrite", "t", false, "Overwrite the output file(s) if they already exist")

	rootCmd.Flags().StringVarP(&prgArgs.PhoneNumber, "my-number", "m", "", "Your phone number, without the country code")
	rootCmd.Flags().StringVarP(&prgArgs.Password, "password", "p", "", "The password signal told you to write down when you enable backups (a bunch of numbers)")
}

// Execute runs the whole program using Cobra to parse CLI args
func Execute(ctx context.Context, version string) error {
	rootCmd.Version = version
	return rootCmd.ExecuteContext(ctx)
}

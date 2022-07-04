// Package conversion is for the process of taking an input format and running it through various converters to arrive at an output format
package conversion

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/jamescostian/signal-to-sms/formats"
	"github.com/jamescostian/signal-to-sms/utils/attachments"
	"github.com/pkg/errors"
)

type Conversion struct {
	// ExternalAttachments should only be used by a MsgConverter which has a MsgFormat that has IncludesAttachments set to false.
	// If both the input and output formats include attachments within them, then this should be nil.
	// This is also made read-only or write-only if only one of the formats have attachments.
	ExternalAttachments attachments.Store

	MsgIn  interface{}
	MsgOut interface{}

	MsgFileFlags       int
	MsgFilePermissions os.FileMode

	// Where the final attachment output belongs. Will be nil if the final output MsgFormat includes attachments in it
	attachmentOut attachments.Store
	// Path to the file where the final message output belongs
	msgOutPath string
	// Path to the file where messages are read from. This can be the original input file, or a temporary file
	msgInPath string

	conversionPath   []MsgConverter
	currConverterIdx int

	// If true, then treat this like it's the last conversion when it comes to things like cleanUp functions
	unrecoverableError bool

	// cleanUpFuncs will run the next time CleanUp is called
	cleanUpFuncs []func(conversionSuccessful bool) error
}

// CurrentConverter returns the MsgConverter that will be used next time Convert is called, or nil if no converter will be used because the Conversion is finished
func (c *Conversion) CurrentConverter() *MsgConverter {
	if c.Finished() {
		return nil
	}
	return &c.conversionPath[c.currConverterIdx]
}

// Finished says whether or not the Conversion has successfully finished (because it has gone through all the MsgConverters needed for this Conversion)
func (c *Conversion) Finished() bool {
	return len(c.conversionPath) <= c.currConverterIdx
}

// IsFinalConversion lets you know if this is the last conversion between formats in the entire conversion path
func (c *Conversion) IsFinalConversion() bool {
	return len(c.conversionPath) <= c.currConverterIdx+1
}

// convertToNextFormat only runs 1 MsgConverter, but sets it up to run perfectly - things like temporary attachments, cleanups, etc are handled.
func (c *Conversion) convertToNextFormat(ctx context.Context) (err error) {
	converter := c.CurrentConverter()
	if converter == nil {
		return fmt.Errorf("attempting to run a conversion with no current converter")
	}

	// Set up attachments, the message input format, and the message output format.
	// Note that ExternalAttachments will be set to an attachment.Store that's been purposefully limited (e.g. calling Close on it will fail).
	// The real, non-limited attachment.Store will be put into the realExtAttachments variable.
	// This real attachment.Store will be restored later on (in a deferred block) as ExternalAttachments at the end of this function, so that:
	// 1. During the next run, setupExtAttachmentsFor can use the real ExternalAttachments if needed, and impose whichever limits on it are correct
	// 2. Clean up functions can access the real attachments
	realExtAttachments, err := c.setupExtAttachmentsFor(converter)
	if err != nil {
		return err
	}

	// The result of calling Convert() on the current converter will be stored here
	var converterOutput interface{}
	// If there's a panic along the way, recover from it.
	// Because we won't know for sure if there's an error until this deferred func, several other things happen here as well, including setting up
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during conversion: %v", r)
		}
		// Restore the real external attachments now so that clean up functions can work with them.
		// The reason this is here and not inside a defer block near where realExtAttachments is declared is that go's defer blocks run in opposite-order, so it could mess things up.
		c.ExternalAttachments = realExtAttachments
		// This line of code is here to ensure cleanUp happens even if there's a panic
		errCleaningUp := c.cleanUp(err == nil)
		if err == nil {
			err = errCleaningUp
		}
		// Only advance to the next converter (using the output provided) if this conversion was a success.
		// This is so that CurrentConverter will show the converter that failed, and MsgIn will continue to show the output that caused the error - useful behavior for debugging.
		if err == nil {
			c.MsgIn = converterOutput
			c.currConverterIdx++
		}
	}()

	if converter.InputFormat.OpenForReads != nil {
		var msgIn io.Closer
		if msgIn, err = c.openForReads(converter.InputFormat, c.msgInPath); err != nil {
			return
		}
		c.MsgIn = msgIn // Only set c.MsgIn if there was no error
		c.addCleanUpOnceWrittenToDisk(func(a bool) error {
			return msgIn.Close()
		})
	}
	if converter.OutputFormat.OpenForWrites != nil {
		outPath := c.msgOutPath
		if !c.IsFinalConversion() {
			// This isn't the final conversion, so store the outputs in a temporary file
			var f *os.File
			if f, err = ioutil.TempFile("", ""); err != nil {
				return
			}
			c.AddFinalCleanUpFn(func(_ bool) error {
				return os.Remove(f.Name())
			})
			outPath = f.Name()
			c.AddCleanUpFn(func(conversionSuccessful bool) error {
				if conversionSuccessful {
					c.msgInPath = outPath
				}
				return nil
			})
		}
		var msgOut io.Closer
		msgOut, err = c.openForWrites(converter.OutputFormat, outPath)
		if err != nil {
			return
		}
		c.MsgOut = msgOut
		c.addCleanUpOnceWrittenToDisk(func(a bool) error {
			return msgOut.Close()
		})
	}

	// Finally, now that everything has been set up, actually run the conversion, and then restore
	converterOutput, err = converter.Convert(ctx, c)
	return
}

func (c *Conversion) openForReads(format formats.MsgFormat, path string) (io.Closer, error) {
	opened, err := format.OpenForReads(path)
	return opened, errors.Wrapf(err, "unable to open input format %v for reads", format.Name)
}
func (c *Conversion) openForWrites(format formats.MsgFormat, path string) (io.Closer, error) {
	opened, err := format.OpenForWrites(path, c.MsgFileFlags, c.MsgFilePermissions)
	return opened, errors.Wrapf(err, "unable to open input format %v for writes", format.Name)
}

func (c *Conversion) shouldUseFinalAttachmentOut() bool {
	if c.attachmentOut == nil {
		return false
	}
	for i := c.currConverterIdx + 1; i < len(c.conversionPath); i++ {
		if c.conversionPath[i].OutputFormat.IncludesAttachments {
			return false
		}
	}
	return true
}

func (c *Conversion) setupExtAttachmentsFor(converter *MsgConverter) (realExtAttachments attachments.Store, err error) {
	// Make sure ExternalAttachments is untouchable if they aren't needed, and make sure it is usable if it is needed.
	// Say you write a converter for encrypted -> prototext. You'll expect that c.ExternalAttachments has no attachments in it, right?
	// But imagine there's also a converter for ??? -> encrypted where ??? has separate attachments, and a user converts ??? -> encrypted -> prototext.
	// For ???, attachments are loaded up, into ExternalAttachments, but then encrypted -> prototext now has different behavior - c.ExternalAttachments has some attachments in it!
	// To keep the behavior consistent for individual converters, after the ??? -> encrypted conversion, ExternalAttachments will be replaced with a blank attachments.Store
	if converter.OutputFormat.IncludesAttachments && converter.OutputFormat.OpenForWrites != nil {
		c.AddCleanUpFn(func(conversionSuccessful bool) error {
			// If it's all over and ExternalAttachments gets closed and set to nil by the time we get here, then respect that
			if c.ExternalAttachments == nil {
				return nil
			}
			err := c.ExternalAttachments.Close()
			if err != nil {
				return err
			}
			if c.shouldUseFinalAttachmentOut() {
				c.ExternalAttachments = c.attachmentOut
				return nil
			}
			// When this is closed, it will also delete any file that it creates
			c.ExternalAttachments, err = attachments.NewTempStore()
			return err
		})
	}

	// Restrict control over ExternalAttachments based on what the formats are, and restore the full capabilities of the formats at the very end
	realExtAttachments = c.ExternalAttachments
	c.ExternalAttachments = attachments.Unclosable(c.ExternalAttachments)
	if converter.InputFormat.IncludesAttachments {
		c.ExternalAttachments = attachments.NoReads(c.ExternalAttachments)
	}
	if converter.OutputFormat.IncludesAttachments {
		c.ExternalAttachments = attachments.NoWrites(c.ExternalAttachments)
	}
	return
}

package conversion

import (
	"fmt"

	"go.uber.org/multierr"
)

func runCleanUpFunc(cleanUpFunc func(bool) error, conversionSuccessful bool) (err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			err = fmt.Errorf("panic during conversion clean up: %v", r)
		}
	}()
	err = cleanUpFunc(conversionSuccessful)
	return
}

// Run all those scheduled functions
func (c *Conversion) cleanUp(conversionSuccessful bool) error {
	// Run the normal clean up functions - but only the ones already added!
	// This way you can write a clean up function that adds a clean up function for the *next* conversion
	cleanUpFuncs := make([]func(bool) error, len(c.cleanUpFuncs))
	copy(cleanUpFuncs, c.cleanUpFuncs)
	c.cleanUpFuncs = c.cleanUpFuncs[:0] // Clear out the clean up functions queued up since they're all about to be executed
	var err error
	for _, fn := range cleanUpFuncs {
		multierr.AppendInto(&err, runCleanUpFunc(fn, conversionSuccessful))
	}
	return err
}

// AddCleanUpFn will schedule a function to run after the conversion is done (whether it's successful or not)
func (c *Conversion) AddCleanUpFn(fn func(conversionSuccessful bool) error) {
	c.cleanUpFuncs = append(c.cleanUpFuncs, fn)
}

// AddFinalCleanUpFn adds a function that will run after the last conversion finishes, or if there was an error during the conversion.
// If there's an error during a step of the conversion, your function will still be run.
func (c *Conversion) AddFinalCleanUpFn(fn func(conversionSuccessful bool) error) {
	c.AddCleanUpFn(func(conversionSuccessful bool) error {
		if c.IsFinalConversion() || c.unrecoverableError || !conversionSuccessful {
			return fn(conversionSuccessful)
		}
		// This isn't the last conversion, so schedule the function to run during the next conversion
		c.AddFinalCleanUpFn(fn)
		return nil
	})
}

func (c *Conversion) addCleanUpOnceWrittenToDisk(fn func(conversionSuccessful bool) error) {
	c.AddCleanUpFn(func(conversionSuccessful bool) error {
		if c.CurrentConverter().OutputFormat.OpenForWrites != nil || c.unrecoverableError {
			return fn(conversionSuccessful)
		}
		// This conversion was to a format that doesn't get persisted to disk, so try again during the next conversion
		c.addCleanUpOnceWrittenToDisk(fn)
		return nil
	})
}

//go:build deterministic_output
// +build deterministic_output

package prototext

// Required in order to use go:linkname
import _ "unsafe"

// The following makes a disable function that will disable the randomness in prototext
//go:linkname disable google.golang.org/protobuf/internal/detrand.Disable
func disable()

// This part makes prototext output always be deterministic.
// This code won't run unless you build with the deterministic_output tag
func init() {
	disable()
}

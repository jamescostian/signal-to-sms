package tosqlite

import "github.com/jamescostian/signal-to-sms/utils/proto/frameio"

func ChangeShouldIgnoreImport(writer frameio.FrameWriter, fn func(s string) bool) {
	importer := writer.(*importer)
	importer.ShouldIgnoreImport = fn
}

func ShouldIgnoreImport(statement string) bool {
	return shouldIgnoreImport(statement)
}
func FastShouldIgnoreImport(statement string) bool {
	return fastShouldIgnoreImport(statement)
}

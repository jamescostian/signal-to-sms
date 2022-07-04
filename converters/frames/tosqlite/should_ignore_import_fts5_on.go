//go:build sqlite_fts5
// +build sqlite_fts5

package tosqlite

import "strings"

// This exists because:
// 1. It's faster to have specific code generated once, rather than continuously check things like fasterInsert during a hot loop
// 2. In addition to adjusting itself by fasterInsert, it also has to adjust itself based on the sqlite_fts5 build tag

func shouldIgnoreImport(statement string) bool {
	// Ignore things that will cause trouble during import because they are managed by sqlite itself
	return strings.HasPrefix(statement, "CREATE TABLE sqlite_")
}

func fastShouldIgnoreImport(statement string) bool {
	// Perf enhancement
	if !strings.HasPrefix(statement, "CREATE ") {
		return false
	}
	// Ignore things that will cause trouble during import because they are managed by sqlite itself
	if strings.HasPrefix(statement, "CREATE TABLE sqlite_") {
		return true
	}
	// Don't bother indexing data for full text search
	if strings.HasPrefix(statement, "CREATE TRIGGER") && strings.Contains(statement, "fts") {
		return true
	}
	// Don't bother creating indexes either
	if strings.HasPrefix(statement, "CREATE INDEX") {
		return true
	}
	return false
}

package sqlite

import (
	"context"
	"database/sql"
	"io"
)

// DBOrTx is an interface satisfied by *sql.DB and *sql.Tx.
// It lets you run whatever SQL queries you need, either in a transaction or not in a transaction
type DBOrTx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func CommitOrClose(dbOrTx DBOrTx) error {
	if closer, ok := dbOrTx.(io.Closer); ok {
		return closer.Close()
	}
	return dbOrTx.(*sql.Tx).Commit()
}

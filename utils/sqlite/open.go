package sqlite

import (
	"context"
	"database/sql"

	// Import go-sqlite3 library
	_ "github.com/mattn/go-sqlite3"
)

func Open(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite3", dbPath)
}

func OpenToTx(dbPath string) (*ClosableTx, error) {
	db, err := Open(dbPath)
	if err != nil {
		return nil, err
	}
	return NewClosableTx(db)
}

func NewClosableTx(db *sql.DB) (*ClosableTx, error) {
	tx, err := db.Begin()
	return &ClosableTx{DB: db, Tx: tx}, err
}

type ClosableTx struct {
	DB *sql.DB
	Tx *sql.Tx
}

func (dbAndTx *ClosableTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return dbAndTx.Tx.Exec(query, args...)
}
func (dbAndTx *ClosableTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return dbAndTx.Tx.ExecContext(ctx, query, args...)
}
func (dbAndTx *ClosableTx) Prepare(query string) (*sql.Stmt, error) {
	return dbAndTx.Tx.Prepare(query)
}
func (dbAndTx *ClosableTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return dbAndTx.Tx.PrepareContext(ctx, query)
}
func (dbAndTx *ClosableTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return dbAndTx.Tx.Query(query, args...)
}
func (dbAndTx *ClosableTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return dbAndTx.Tx.QueryContext(ctx, query, args...)
}
func (dbAndTx *ClosableTx) QueryRow(query string, args ...interface{}) *sql.Row {
	return dbAndTx.Tx.QueryRow(query, args...)
}
func (dbAndTx *ClosableTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return dbAndTx.Tx.QueryRowContext(ctx, query, args...)
}
func (dbAndTx *ClosableTx) Close() error {
	return dbAndTx.Commit()
}
func (dbAndTx *ClosableTx) Commit() error {
	errCommitting := dbAndTx.Tx.Commit()
	errClosing := dbAndTx.DB.Close()
	// Being unable to close is a more pressing matter than committing, for privacy reasons
	if errClosing != nil {
		return errClosing
	}
	return errCommitting
}
func (dbAndTx *ClosableTx) Rollback() error {
	errRollingBack := dbAndTx.Tx.Rollback()
	errClosing := dbAndTx.DB.Close()
	// Being unable to close is a more pressing matter than rolling back, for privacy reasons
	if errClosing != nil {
		return errClosing
	}
	return errRollingBack
}

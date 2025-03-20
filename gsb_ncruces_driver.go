//go:build ncruces_driver

package go_sqlite_bench

import (
	"database/sql"

	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// Notes
//
// * Transactions
//
// If `TxOptions.IsolationLevel` is `sql.LevelSerializable`, ncruces issues
// `BEGIN IMMEDIATE`.
//
// If `TxOptions.IsolationLevel` is `sql.LevelLinearizable`, ncruces issues
// `BEGIN EXCLUSIVE`.
//
// If `TxOptions.IsolationLevel` is `sql.LevelDefault`, and
// `TxOptions.ReadOnly` is true, ncruces issues `BEGIN DEFERRED`.
//
// If `TxOptions.IsolationLevel` is `sql.LevelDefault`, and
// `TxOptions.ReadOnly` is false, ncruces issues `BEGIN <txlock>`, where
// <txlock> is the `_txlock` value passed in the Open DSN..

func init() {
	sqlite3.Initialize()
}

func OpenDB(filename string) (*sql.DB, error) {
	return sql.Open("sqlite3", "file:"+filename+"?_pragma=busy_timeout(10000)&_pragma=foreign_keys(true)&_pragma=journal_mode(wal)&_pragma=synchronous(normal)")
}

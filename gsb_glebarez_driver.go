//go:build glebarez_driver

package go_sqlite_bench

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

// Notes
//
// * Exec
//
// glebarez only steps an `Exec` statement once. For some queries that can
// produce incorrect results.
//
//     https://github.com/glebarez/go-sqlite/blob/e6de9fc0c320a357ddc0201f8086bf44b7603436/sqlite.go#L536
//
// * Transactions
//
// If `TxOptions.ReadOnly` is true, glebarez issues `BEGIN DEFERRED`. If
// `TxOptions.ReadOnly` is false, glebarez issues `BEGIN <txlock>`, where
// <txlock> is the `_txlock` value passed in the Open DSN.
//
//     https://github.com/glebarez/go-sqlite/blob/e6de9fc0c320a357ddc0201f8086bf44b7603436/sqlite.go#L698

func OpenDB(filename string) (*sql.DB, error) {
	return sql.Open("sqlite", "file:"+filename+"?_pragma=busy_timeout(10000)&_pragma=foreign_keys(true)&_pragma=journal_mode(wal)&_pragma=synchronous(normal)&_time_format=sqlite&_txlock=immediate")
}

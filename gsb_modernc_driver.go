//go:build modernc_driver

package go_sqlite_bench

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Notes
//
// * Exec
//
// modernc only steps an `Exec` statement once. For some queries that can
// produce incorrect results.
//
//     https://gitlab.com/cznic/sqlite/-/blob/master/sqlite.go#L558
//
// * Transactions
//
// If `TxOptions.ReadOnly` is true, modernc issues `BEGIN DEFERRED`. If
// `TxOptions.ReadOnly` is false, modernc issues `BEGIN <txlock>`, where
// <txlock> is the `_txlock` value passed in the Open DSN.
//
//     https://gitlab.com/cznic/sqlite/-/blob/master/sqlite.go?ref_type=heads#L730

func OpenDB(filename string) (*sql.DB, error) {
	return sql.Open("sqlite", "file:"+filename+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(true)&_pragma=journal_mode(wal)&_pragma=synchronous(normal)&_time_format=sqlite&_txlock=immediate")
}

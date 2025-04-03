//go:build mattn_driver

package go_sqlite_bench

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Notes
//
// * Exec
//
// mattn only steps an `Exec` statement once. For some queries that can
// produce incorrect results.
//
//     https://github.com/mattn/go-sqlite3/blob/7658c06970ecf5588d8cd930ed1f2de7223f1010/sqlite3.go#L2093C10-L2093C36
//     https://github.com/mattn/go-sqlite3/blob/7658c06970ecf5588d8cd930ed1f2de7223f1010/sqlite3_opt_unlock_notify.c#L49
//     https://github.com/mattn/go-sqlite3/blob/7658c06970ecf5588d8cd930ed1f2de7223f1010/sqlite3.go#L126
//
// * Transactions
//
// mattn doesn't support read and write transactions on the same connection
// because it ignores `TxOptions` completely.
//
//     https://github.com/mattn/go-sqlite3/blob/7658c06970ecf5588d8cd930ed1f2de7223f1010/sqlite3_go18.go#L42
//
// Instead, mattn uses the value of `_txlock` whenever it starts a transaction.
//
//     https://github.com/mattn/go-sqlite3/blob/7658c06970ecf5588d8cd930ed1f2de7223f1010/sqlite3.go#L969
//
// The default value of `_txlock` is, effectively, "deferred".
//
// To get "correct" read/write behaviour you'd need two connection pools,
// configured with two different `_txlock` values.
//
// If you don't do that, you'll get "deferred" transactions for both reading
// and writing. That's not good. It will cause deferred tranactions to get
// upgraded to immediate once SQLite realizes that it needs to write to the
// database. In WAL mode that can cause "database is locked" errors even if
// you've set a timeout.
//
// See https://berthub.eu/articles/posts/a-brief-post-on-sqlite3-database-locked-despite-timeout/
// for a more detailed explanation.

func OpenDB(filename string) (*sql.DB, error) {
	return sql.Open("sqlite3", "file:"+filename+"?_busy_timeout=10000&_foreign_keys=true&_journal_mode=WAL&_synchronous=normal&_mutex=no")
}

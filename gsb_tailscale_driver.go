//go:build tailscale_driver

package go_sqlite_bench

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/tailscale/sqlite"
)

// Notes
//
// * Exec
//
// tailscale only steps an `Exec` statement once. For some queries that can
// produce incorrect results.
//
//     https://github.com/tailscale/sqlite/blob/9328d0478def66b5188b4b80b5db90c670c77b0f/sqlite.go#L514
//
// * Transactions
//
// tailscale always starts a deferred transaction when `TxOptions.ReadOnly`
// is true, or when the connection is read-only. Otherwise it starts an
// immediate transaction. (If `TxOptions.IsolationLevel` is set, it must be
// set to `sql.LevelSerializable`. But it doesn't *need* to be set.)
//
//     https://github.com/tailscale/sqlite/blob/9328d0478def66b5188b4b80b5db90c670c77b0f/sqlite.go#L295
//     https://github.com/tailscale/sqlite/blob/9328d0478def66b5188b4b80b5db90c670c77b0f/sqlite.go#L331

func OpenDB(filename string) (*sql.DB, error) {
	connInitFunc := func(ctx context.Context, conn driver.ConnPrepareContext) error {
		return sqlite.ExecScript(conn.(sqlite.SQLConn), `
			PRAGMA busy_timeout(5000);
			PRAGMA foreign_keys(true);
			PRAGMA journal_mode(WAL);
			PRAGMA synchronous(normal);
		`)
	}
	return sql.OpenDB(sqlite.Connector("file:"+filename, connInitFunc, nil)), nil
}

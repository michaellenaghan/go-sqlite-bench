//go:build zombiezen_direct

package go_sqlite_bench

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type DB struct {
	readPool  *sqlitex.Pool
	writePool *sqlitex.Pool
}

func NewDB(ctx context.Context, filename string, maxReadConnections, maxWriteConnections int) (*DB, error) {
	if !(maxReadConnections >= 0) {
		return nil, errors.New("maxReadConnections must be >= 0")
	}
	if !(maxWriteConnections >= 1) {
		return nil, errors.New("maxWriteConnections must be >= 1")
	}

	var prepareConn = func(conn *sqlite.Conn) error {
		err := sqlitex.ExecuteTransient(conn, "PRAGMA busy_timeout(10000)", &sqlitex.ExecOptions{})
		if err != nil {
			return err
		}
		err = sqlitex.ExecuteTransient(conn, "PRAGMA foreign_keys(true)", &sqlitex.ExecOptions{})
		if err != nil {
			return err
		}
		err = sqlitex.ExecuteTransient(conn, "PRAGMA journal_mode(wal)", &sqlitex.ExecOptions{})
		if err != nil {
			return err
		}
		err = sqlitex.ExecuteTransient(conn, "PRAGMA synchronous(normal)", &sqlitex.ExecOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	if maxReadConnections == 0 {
		pool, err := sqlitex.NewPool(
			filename,
			sqlitex.PoolOptions{
				PoolSize:    maxWriteConnections,
				PrepareConn: prepareConn,
			})
		if err != nil {
			return nil, err
		}

		return &DB{readPool: pool, writePool: pool}, nil
	} else {
		readPool, err := sqlitex.NewPool(
			filename,
			sqlitex.PoolOptions{
				PoolSize:    maxReadConnections,
				PrepareConn: prepareConn,
			})
		if err != nil {
			return nil, err
		}
		writePool, err := sqlitex.NewPool(
			filename,
			sqlitex.PoolOptions{
				PoolSize:    maxWriteConnections,
				PrepareConn: prepareConn,
			})
		if err != nil {
			return nil, err
		}

		return &DB{readPool: readPool, writePool: writePool}, nil
	}
}

func (db *DB) Close() error {
	var readErr error
	var writeErr error

	readErr = db.readPool.Close()
	if db.writePool != db.readPool {
		writeErr = db.writePool.Close()
	}

	if readErr != nil {
		return readErr
	}
	if writeErr != nil {
		return writeErr
	}

	return nil
}

func (db *DB) PrepareDBWithTx(ctx context.Context) error {
	conn, err := db.writePool.Take(ctx)
	if err != nil {
		return err
	}
	defer db.writePool.Put(conn)

	tx, err := sqlitex.ImmediateTransaction(conn)
	if err != nil {
		return err
	}
	defer tx(&err)

	for _, s := range SQLForSchema {
		err := sqlitex.ExecuteTransient(conn, strings.TrimSpace(s), &sqlitex.ExecOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) PopulateDB(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) error {
	conn, err := db.writePool.Take(ctx)
	if err != nil {
		return err
	}
	defer db.writePool.Put(conn)

	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON
	postDate := NewPostDate(posts)

	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON
	commentDate := postDate.CommentDate

	for i := range posts {
		postSeq := i + 1
		title := fmt.Sprintf("Post %d", postSeq)
		content := postContent
		stats := postStats
		created := postDate.NextFormatted()

		err = sqlitex.Execute(conn, SQLForInsertPostWithCreated, &sqlitex.ExecOptions{
			Args: []any{title, content, stats, created},
		})
		if err != nil {
			return err
		}

		postID := conn.LastInsertRowID()

		for j := range comments {
			commentSeq := j + 1
			name := fmt.Sprintf("Comment %d.%d", postSeq, commentSeq)
			content := commentContent
			stats := commentStats
			created := commentDate.NextFormatted()

			err = sqlitex.Execute(conn, SQLForInsertCommentWithCreated, &sqlitex.ExecOptions{
				Args: []any{postID, name, content, stats, created},
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *DB) PopulateDBWithTx(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) error {
	conn, err := db.writePool.Take(ctx)
	if err != nil {
		return err
	}
	defer db.writePool.Put(conn)

	tx, err := sqlitex.ImmediateTransaction(conn)
	if err != nil {
		return err
	}
	defer tx(&err)

	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON
	postDate := NewPostDate(posts)

	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON
	commentDate := postDate.CommentDate

	for i := range posts {
		postSeq := i + 1
		title := fmt.Sprintf("Post %d", postSeq)
		content := postContent
		stats := postStats
		created := postDate.NextFormatted()

		err = sqlitex.Execute(conn, SQLForInsertPostWithCreated, &sqlitex.ExecOptions{
			Args: []any{title, content, stats, created},
		})
		if err != nil {
			return err
		}

		postID := conn.LastInsertRowID()
		if err != nil {
			return err
		}

		for j := range comments {
			commentSeq := j + 1
			name := fmt.Sprintf("Comment %d.%d", postSeq, commentSeq)
			content := commentContent
			stats := commentStats
			created := commentDate.NextFormatted()

			err = sqlitex.Execute(conn, SQLForInsertCommentWithCreated, &sqlitex.ExecOptions{
				Args: []any{postID, name, content, stats, created},
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *DB) PopulateDBWithTxs(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) error {
	conn, err := db.writePool.Take(ctx)
	if err != nil {
		return err
	}
	defer db.writePool.Put(conn)

	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON
	postDate := NewPostDate(posts)

	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON
	commentDate := postDate.CommentDate

	writePostAndComments := func(i int) error {
		tx, err := sqlitex.ImmediateTransaction(conn)
		if err != nil {
			return err
		}
		defer tx(&err)

		postSeq := i + 1
		title := fmt.Sprintf("Post %d", postSeq)
		content := postContent
		stats := postStats
		created := postDate.NextFormatted()

		err = sqlitex.Execute(conn, SQLForInsertPostWithCreated, &sqlitex.ExecOptions{
			Args: []any{title, content, stats, created},
		})
		if err != nil {
			return err
		}

		postID := conn.LastInsertRowID()
		if err != nil {
			return err
		}

		for j := range comments {
			commentSeq := j + 1
			name := fmt.Sprintf("Comment %d.%d", postSeq, commentSeq)
			content := commentContent
			stats := commentStats
			created := commentDate.NextFormatted()

			err = sqlitex.Execute(conn, SQLForInsertCommentWithCreated, &sqlitex.ExecOptions{
				Args: []any{postID, name, content, stats, created},
			})
			if err != nil {
				return err
			}
		}

		return nil
	}

	for i := range posts {
		err = writePostAndComments(i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) CountPosts(ctx context.Context) (int64, error) {
	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareTransient(SQLForCountPosts)
	if err != nil {
		return 0, err
	}
	defer stmt.Finalize()

	n, err := sqlitex.ResultInt64(stmt)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) CountComments(ctx context.Context) (int64, error) {
	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareTransient(SQLForCountComments)
	if err != nil {
		return 0, err
	}
	defer stmt.Finalize()

	n, err := sqlitex.ResultInt64(stmt)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) ReadPost(ctx context.Context, id int64) (*Post, error) {
	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return nil, err
	}
	defer db.readPool.Put(conn)

	p := &Post{ID: id}

	err = sqlitex.ExecuteTransient(conn, SQLForSelectPostByID, &sqlitex.ExecOptions{
		Args: []any{id},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			p.Title = stmt.ColumnText(0)
			p.Content = stmt.ColumnText(1)
			p.Created = stmt.ColumnText(2)
			p.Stats = stmt.ColumnText(3)

			return nil
		},
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *DB) ReadPostWithTx(ctx context.Context, id int64) (*Post, error) {
	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return nil, err
	}
	defer db.readPool.Put(conn)

	tx := sqlitex.Transaction(conn)
	defer tx(&err)

	p := &Post{ID: id}

	err = sqlitex.ExecuteTransient(conn, SQLForSelectPostByID, &sqlitex.ExecOptions{
		Args: []any{id},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			p.Title = stmt.ColumnText(0)
			p.Content = stmt.ColumnText(1)
			p.Created = stmt.ColumnText(2)
			p.Stats = stmt.ColumnText(3)

			return nil
		},
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *DB) ReadPostAndComments(ctx context.Context, id int64) (*Post, []*Comment, error) {
	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer db.readPool.Put(conn)

	p := &Post{ID: id}

	err = sqlitex.ExecuteTransient(conn, SQLForSelectPostByID, &sqlitex.ExecOptions{
		Args: []any{id},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			p.Title = stmt.ColumnText(0)
			p.Content = stmt.ColumnText(1)
			p.Created = stmt.ColumnText(2)
			p.Stats = stmt.ColumnText(3)

			return nil
		},
	})
	if err != nil {
		return nil, nil, err
	}

	cs := make([]*Comment, 0)

	err = sqlitex.ExecuteTransient(conn, SQLForSelectCommentsByPostID, &sqlitex.ExecOptions{
		Args: []any{id},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			c := &Comment{}

			c.ID = stmt.ColumnInt64(0)
			c.Name = stmt.ColumnText(1)
			c.Content = stmt.ColumnText(2)
			c.Created = stmt.ColumnText(3)
			c.Stats = stmt.ColumnText(4)

			cs = append(cs, c)

			return nil
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return p, cs, nil
}

func (db *DB) ReadPostAndCommentsWithTx(ctx context.Context, id int64) (*Post, []*Comment, error) {
	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer db.readPool.Put(conn)

	tx := sqlitex.Transaction(conn)
	defer tx(&err)

	p := &Post{ID: id}

	err = sqlitex.ExecuteTransient(conn, SQLForSelectPostByID, &sqlitex.ExecOptions{
		Args: []any{id},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			p.Title = stmt.ColumnText(0)
			p.Content = stmt.ColumnText(1)
			p.Created = stmt.ColumnText(2)
			p.Stats = stmt.ColumnText(3)

			return nil
		},
	})
	if err != nil {
		return nil, nil, err
	}

	cs := make([]*Comment, 0)

	err = sqlitex.ExecuteTransient(conn, SQLForSelectCommentsByPostID, &sqlitex.ExecOptions{
		Args: []any{id},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			c := &Comment{}

			c.ID = stmt.ColumnInt64(0)
			c.Name = stmt.ColumnText(1)
			c.Content = stmt.ColumnText(2)
			c.Created = stmt.ColumnText(3)
			c.Stats = stmt.ColumnText(4)

			cs = append(cs, c)

			return nil
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return p, cs, nil
}

func (db *DB) WritePost(ctx context.Context, title, content, stats string) (int64, error) {
	conn, err := db.writePool.Take(ctx)
	if err != nil {
		return 0, err
	}
	defer db.writePool.Put(conn)

	err = sqlitex.ExecuteTransient(conn, SQLForInsertPost, &sqlitex.ExecOptions{
		Args: []any{title, content, stats},
	})
	if err != nil {
		return 0, err
	}

	postID := conn.LastInsertRowID()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (db *DB) WritePostWithTx(ctx context.Context, title, content, stats string) (int64, error) {
	conn, err := db.writePool.Take(ctx)
	if err != nil {
		return 0, err
	}
	defer db.writePool.Put(conn)

	tx, err := sqlitex.ImmediateTransaction(conn)
	if err != nil {
		return 0, err
	}
	defer tx(&err)

	err = sqlitex.ExecuteTransient(conn, SQLForInsertPost, &sqlitex.ExecOptions{
		Args: []any{title, content, stats},
	})
	if err != nil {
		return 0, err
	}

	postID := conn.LastInsertRowID()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (db *DB) WritePostAndComments(ctx context.Context, postTitle, postContent, postStats string, comments int, commentName, commentContent, commentStats string) (int64, error) {
	conn, err := db.writePool.Take(ctx)
	if err != nil {
		return 0, err
	}
	defer db.writePool.Put(conn)

	err = sqlitex.ExecuteTransient(conn, SQLForInsertPost, &sqlitex.ExecOptions{
		Args: []any{postTitle, postContent, postStats},
	})
	if err != nil {
		return 0, err
	}

	postID := conn.LastInsertRowID()
	if err != nil {
		return 0, err
	}

	for range comments {
		err = sqlitex.ExecuteTransient(conn, SQLForInsertComment, &sqlitex.ExecOptions{
			Args: []any{postID, commentName, commentContent, commentStats},
		})
		if err != nil {
			return 0, err
		}
	}

	return postID, nil
}

func (db *DB) WritePostAndCommentsWithTx(ctx context.Context, postTitle, postContent, postStats string, comments int, commentName, commentContent, commentStats string) (int64, error) {
	conn, err := db.writePool.Take(ctx)
	if err != nil {
		return 0, err
	}
	defer db.writePool.Put(conn)

	tx, err := sqlitex.ImmediateTransaction(conn)
	if err != nil {
		return 0, err
	}
	defer tx(&err)

	err = sqlitex.ExecuteTransient(conn, SQLForInsertPost, &sqlitex.ExecOptions{
		Args: []any{postTitle, postContent, postStats},
	})
	if err != nil {
		return 0, err
	}

	postID := conn.LastInsertRowID()
	if err != nil {
		return 0, err
	}

	for range comments {
		err = sqlitex.ExecuteTransient(conn, SQLForInsertComment, &sqlitex.ExecOptions{
			Args: []any{postID, commentName, commentContent, commentStats},
		})
		if err != nil {
			return 0, err
		}
	}

	return postID, nil
}

// ===

func (db *DB) QueryCorrelated(ctx context.Context) (int, error) {
	n := 0

	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	err = sqlitex.Execute(conn, SQLForQueryCorrelated, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			n += 1

			return nil
		},
	})
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) QueryGroupBy(ctx context.Context) (int, error) {
	n := 0

	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	err = sqlitex.Execute(conn, SQLForQueryGroupBy, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			n += 1

			return nil
		},
	})
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) QueryJSON(ctx context.Context) (int, error) {
	n := 0

	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	err = sqlitex.Execute(conn, SQLForQueryJSON, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			n += 1

			return nil
		},
	})
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) QueryNonRecursiveCTE(ctx context.Context) (int, error) {
	n := 0

	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	err = sqlitex.Execute(conn, SQLForQueryNonRecursiveCTE, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			n += 1

			return nil
		},
	})
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) QueryOrderBy(ctx context.Context) (int, error) {
	n := 0

	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	err = sqlitex.Execute(conn, SQLForQueryOrderBy, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			n += 1

			return nil
		},
	})
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) QueryRecursiveCTE(ctx context.Context) (int, error) {
	n := 0

	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	err = sqlitex.Execute(conn, SQLForQueryRecursiveCTE, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			n += 1

			return nil
		},
	})
	if err != nil {
		return 0, err
	}

	return n, nil
}

// ===

func (db *DB) Options(ctx context.Context) ([]string, error) {
	options := make([]string, 0)

	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return nil, err
	}
	defer db.readPool.Put(conn)

	err = sqlitex.ExecuteTransient(conn, "PRAGMA compile_options", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			option := stmt.ColumnText(0)

			options = append(options, option)

			return nil
		},
	})
	if err != nil {
		return nil, err
	}

	return options, nil
}

func (db *DB) Pragmas(ctx context.Context, names []string) ([]string, error) {
	pragmas := make([]string, 0)

	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return nil, err
	}
	defer db.readPool.Put(conn)

	for _, name := range names {
		err = sqlitex.ExecuteTransient(conn, "PRAGMA"+" "+name, &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				value := stmt.ColumnText(0)

				pragmas = append(pragmas, fmt.Sprintf("%s=%s", name, value))

				return nil
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return pragmas, nil
}

func (db *DB) Select1(ctx context.Context) error {
	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return err
	}
	defer db.readPool.Put(conn)

	err = sqlitex.ExecuteTransient(conn, "SELECT 1", &sqlitex.ExecOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Time(ctx context.Context, in time.Time) (time.Time, error) {
	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return time.Time{}, err
	}
	defer db.readPool.Put(conn)

	var out time.Time

	err = sqlitex.ExecuteTransient(conn, "SELECT ?", &sqlitex.ExecOptions{
		Args: []any{in.Format(time.RFC3339)},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			out, err = time.Parse(time.RFC3339, stmt.ColumnText(0))
			if err != nil {
				return err
			}
			return nil
		},
	})
	if err != nil {
		return time.Time{}, err
	}

	return out, nil
}

func (db *DB) Version(ctx context.Context) (string, error) {
	conn, err := db.readPool.Take(ctx)
	if err != nil {
		return "", err
	}
	defer db.readPool.Put(conn)

	var version string

	err = sqlitex.ExecuteTransient(conn, "SELECT sqlite_version()", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			version = stmt.ColumnText(0)
			return nil
		},
	})
	if err != nil {
		return "", err
	}

	return version, nil
}

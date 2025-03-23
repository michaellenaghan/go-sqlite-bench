//go:build ncruces_direct

package go_sqlite_bench

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ncruces/go-sqlite3"

	_ "github.com/ncruces/go-sqlite3/embed"
)

type DB struct {
	readPool  *Pool[*Conn]
	writePool *Pool[*Conn]
}

type Conn struct {
	*sqlite3.Conn

	prepared map[string]*sqlite3.Stmt
}

func init() {
	sqlite3.Initialize()
}

func NewDB(ctx context.Context, filename string, maxReadConnections, maxWriteConnections int) (*DB, error) {
	if !(maxReadConnections >= 0) {
		return nil, errors.New("maxReadConnections must be >= 0")
	}
	if !(maxWriteConnections >= 1) {
		return nil, errors.New("maxWriteConnections must be >= 1")
	}

	if maxReadConnections == 0 {
		pool, err := newPool(filename, 0, maxWriteConnections, 0)
		if err != nil {
			return nil, err
		}

		return &DB{readPool: pool, writePool: pool}, nil
	} else {
		readPool, err := newPool(filename, 0, maxReadConnections, 0)
		if err != nil {
			return nil, err
		}
		writePool, err := newPool(filename, 0, maxWriteConnections, 0)
		if err != nil {
			readPool.Stop()
			return nil, err
		}

		return &DB{readPool: readPool, writePool: writePool}, nil
	}
}

func newPool(filename string, minConnections, maxConnections int, maxConnectionIdleTime time.Duration) (*Pool[*Conn], error) {
	pool, err := NewPool(
		minConnections,
		maxConnections,
		maxConnectionIdleTime,
		func() (*Conn, error) {
			// "Order matters: encryption keys, busy timeout and locking mode
			// should be the first PRAGMAs set, in that order."
			// https://github.com/ncruces/go-sqlite3/blob/main/driver/driver.go
			conn, err := sqlite3.Open("file:" + filename + "?_pragma=busy_timeout(10000)&_pragma=journal_mode(wal)&_pragma=synchronous(normal)")
			if err != nil {
				return nil, err
			}

			return &Conn{Conn: conn}, nil
		}, func(conn *Conn) {
			defer conn.Close()

			if conn.Stmts() != nil {
				for stmt := range conn.Stmts() {
					err := stmt.Close()
					if err != nil {
						log.Printf("failed to close: %v", err)
					}
				}
			}
		},
	)
	if err != nil {
		return nil, err
	}

	err = pool.Start(false)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (db *DB) Close() error {
	db.writePool.Stop()

	if db.readPool != db.writePool {
		db.readPool.Stop()
	}

	return nil
}

func (db *DB) PrepareDBWithTx(ctx context.Context) error {
	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return err
	}
	defer db.writePool.Put(conn)

	tx, err := conn.BeginImmediate()
	if err != nil {
		return err
	}
	defer tx.End(&err)

	for _, s := range SQLForSchema {
		err := conn.Exec(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) PopulateDB(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) error {
	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return err
	}
	defer db.writePool.Put(conn)

	postStmt, _, err := conn.PrepareAndPersist(SQLForInsertPostWithCreated)
	if err != nil {
		return err
	}

	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON
	postDate := NewPostDate(posts)

	commentStmt, _, err := conn.PrepareAndPersist(SQLForInsertCommentWithCreated)
	if err != nil {
		return err
	}

	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON
	commentDate := postDate.CommentDate

	for i := range posts {
		postSeq := i + 1
		title := fmt.Sprintf("Post %d", postSeq)
		content := postContent
		stats := postStats
		created := postDate.NextFormatted()

		err = postStmt.BindText(1, title)
		if err != nil {
			return err
		}
		err = postStmt.BindText(2, content)
		if err != nil {
			return err
		}
		err = postStmt.BindText(3, stats)
		if err != nil {
			return err
		}
		err = postStmt.BindText(4, created)
		if err != nil {
			return err
		}

		err = postStmt.Exec()
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

			err = commentStmt.BindInt64(1, postID)
			if err != nil {
				return err
			}
			err = commentStmt.BindText(2, name)
			if err != nil {
				return err
			}
			err = commentStmt.BindText(3, content)
			if err != nil {
				return err
			}
			err = commentStmt.BindText(4, stats)
			if err != nil {
				return err
			}
			err = commentStmt.BindText(5, created)
			if err != nil {
				return err
			}

			err = commentStmt.Exec()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *DB) PopulateDBWithTx(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) error {
	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return err
	}
	defer db.writePool.Put(conn)

	tx, err := conn.BeginImmediate()
	if err != nil {
		return err
	}
	defer tx.End(&err)

	postStmt, _, err := conn.PrepareAndPersist(SQLForInsertPostWithCreated)
	if err != nil {
		return err
	}

	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON
	postDate := NewPostDate(posts)

	commentStmt, _, err := conn.PrepareAndPersist(SQLForInsertCommentWithCreated)
	if err != nil {
		return err
	}

	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON
	commentDate := postDate.CommentDate

	for i := range posts {
		postSeq := i + 1
		title := fmt.Sprintf("Post %d", postSeq)
		content := postContent
		stats := postStats
		created := postDate.NextFormatted()

		err = postStmt.BindText(1, title)
		if err != nil {
			return err
		}
		err = postStmt.BindText(2, content)
		if err != nil {
			return err
		}
		err = postStmt.BindText(3, stats)
		if err != nil {
			return err
		}
		err = postStmt.BindText(4, created)
		if err != nil {
			return err
		}

		err = postStmt.Exec()
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

			err = commentStmt.BindInt64(1, postID)
			if err != nil {
				return err
			}
			err = commentStmt.BindText(2, name)
			if err != nil {
				return err
			}
			err = commentStmt.BindText(3, content)
			if err != nil {
				return err
			}
			err = commentStmt.BindText(4, stats)
			if err != nil {
				return err
			}
			err = commentStmt.BindText(5, created)
			if err != nil {
				return err
			}

			err = commentStmt.Exec()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *DB) PopulateDBWithTxs(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) error {
	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return err
	}
	defer db.writePool.Put(conn)

	postStmt, _, err := conn.PrepareAndPersist(SQLForInsertPostWithCreated)
	if err != nil {
		return err
	}

	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON
	postDate := NewPostDate(posts)

	commentStmt, _, err := conn.PrepareAndPersist(SQLForInsertCommentWithCreated)
	if err != nil {
		return err
	}

	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON
	commentDate := postDate.CommentDate

	writePostAndComments := func(i int) error {
		tx, err := conn.BeginImmediate()
		if err != nil {
			return err
		}
		defer tx.End(&err)

		postSeq := i + 1
		title := fmt.Sprintf("Post %d", postSeq)
		content := postContent
		stats := postStats
		created := postDate.NextFormatted()

		err = postStmt.BindText(1, title)
		if err != nil {
			return err
		}
		err = postStmt.BindText(2, content)
		if err != nil {
			return err
		}
		err = postStmt.BindText(3, stats)
		if err != nil {
			return err
		}
		err = postStmt.BindText(4, created)
		if err != nil {
			return err
		}

		err = postStmt.Exec()
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

			err = commentStmt.BindInt64(1, postID)
			if err != nil {
				return err
			}
			err = commentStmt.BindText(2, name)
			if err != nil {
				return err
			}
			err = commentStmt.BindText(3, content)
			if err != nil {
				return err
			}
			err = commentStmt.BindText(4, stats)
			if err != nil {
				return err
			}
			err = commentStmt.BindText(5, created)
			if err != nil {
				return err
			}

			err = commentStmt.Exec()
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
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareAndPersist(SQLForCountPosts)
	if err != nil {
		return 0, err
	}

	var n int64

	for stmt.Step() {
		n = stmt.ColumnInt64(0)
	}

	err = stmt.Reset()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) CountComments(ctx context.Context) (int64, error) {
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareAndPersist(SQLForCountComments)
	if err != nil {
		return 0, err
	}

	var n int64

	for stmt.Step() {
		n = stmt.ColumnInt64(0)
	}

	err = stmt.Reset()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) ReadPost(ctx context.Context, id int64) (*Post, error) {
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareAndPersist(SQLForSelectPostByID)
	if err != nil {
		return nil, err
	}

	err = stmt.BindInt64(1, id)
	if err != nil {
		return nil, err
	}

	p := &Post{ID: id}

	for stmt.Step() {
		p.Title = stmt.ColumnText(0)
		p.Content = stmt.ColumnText(1)
		p.Created = stmt.ColumnText(2)
		p.Stats = stmt.ColumnText(3)
	}

	err = stmt.Reset()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *DB) ReadPostWithTx(ctx context.Context, id int64) (*Post, error) {
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer db.readPool.Put(conn)

	tx := conn.Begin()
	defer tx.End(&err)

	stmt, _, err := conn.PrepareAndPersist(SQLForSelectPostByID)
	if err != nil {
		return nil, err
	}

	err = stmt.BindInt64(1, id)
	if err != nil {
		return nil, err
	}

	p := &Post{ID: id}

	for stmt.Step() {
		p.Title = stmt.ColumnText(0)
		p.Content = stmt.ColumnText(1)
		p.Created = stmt.ColumnText(2)
		p.Stats = stmt.ColumnText(3)
	}

	err = stmt.Reset()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *DB) ReadPostAndComments(ctx context.Context, id int64) (*Post, []*Comment, error) {
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareAndPersist(SQLForSelectPostByID)
	if err != nil {
		return nil, nil, err
	}

	err = stmt.BindInt64(1, id)
	if err != nil {
		return nil, nil, err
	}

	p := &Post{ID: id}

	for stmt.Step() {
		p.Title = stmt.ColumnText(0)
		p.Content = stmt.ColumnText(1)
		p.Created = stmt.ColumnText(2)
		p.Stats = stmt.ColumnText(3)
	}

	err = stmt.Reset()
	if err != nil {
		return nil, nil, err
	}

	stmt, _, err = conn.PrepareAndPersist(SQLForSelectCommentsByPostID)
	if err != nil {
		return nil, nil, err
	}

	err = stmt.BindInt64(1, id)
	if err != nil {
		return nil, nil, err
	}

	cs := make([]*Comment, 0)

	for stmt.Step() {
		c := &Comment{}

		c.ID = stmt.ColumnInt64(0)
		c.Name = stmt.ColumnText(1)
		c.Content = stmt.ColumnText(2)
		c.Created = stmt.ColumnText(3)
		c.Stats = stmt.ColumnText(4)

		cs = append(cs, c)
	}

	err = stmt.Reset()
	if err != nil {
		return nil, nil, err
	}

	return p, cs, nil
}

func (db *DB) ReadPostAndCommentsWithTx(ctx context.Context, id int64) (*Post, []*Comment, error) {
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer db.readPool.Put(conn)

	tx := conn.Begin()
	defer tx.End(&err)

	stmt, _, err := conn.PrepareAndPersist(SQLForSelectPostByID)
	if err != nil {
		return nil, nil, err
	}

	err = stmt.BindInt64(1, id)
	if err != nil {
		return nil, nil, err
	}

	p := &Post{ID: id}

	for stmt.Step() {
		p.Title = stmt.ColumnText(0)
		p.Content = stmt.ColumnText(1)
		p.Created = stmt.ColumnText(2)
		p.Stats = stmt.ColumnText(3)
	}

	err = stmt.Reset()
	if err != nil {
		return nil, nil, err
	}

	stmt, _, err = conn.PrepareAndPersist(SQLForSelectCommentsByPostID)
	if err != nil {
		return nil, nil, err
	}

	err = stmt.BindInt64(1, id)
	if err != nil {
		return nil, nil, err
	}

	cs := make([]*Comment, 0)

	for stmt.Step() {
		c := &Comment{}

		c.ID = stmt.ColumnInt64(0)
		c.Name = stmt.ColumnText(1)
		c.Content = stmt.ColumnText(2)
		c.Created = stmt.ColumnText(3)
		c.Stats = stmt.ColumnText(4)

		cs = append(cs, c)
	}

	err = stmt.Reset()
	if err != nil {
		return nil, nil, err
	}

	return p, cs, nil
}

func (db *DB) WritePost(ctx context.Context, title, content, stats string) (int64, error) {
	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.writePool.Put(conn)

	stmt, _, err := conn.PrepareAndPersist(SQLForInsertPost)
	if err != nil {
		return 0, err
	}

	err = stmt.BindText(1, title)
	if err != nil {
		return 0, err
	}
	err = stmt.BindText(2, content)
	if err != nil {
		return 0, err
	}
	err = stmt.BindText(3, stats)
	if err != nil {
		return 0, err
	}

	err = stmt.Exec()
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
	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.writePool.Put(conn)

	tx, err := conn.BeginImmediate()
	if err != nil {
		return 0, err
	}
	defer tx.End(&err)

	stmt, _, err := conn.PrepareAndPersist(SQLForInsertPost)
	if err != nil {
		return 0, err
	}

	err = stmt.BindText(1, title)
	if err != nil {
		return 0, err
	}
	err = stmt.BindText(2, content)
	if err != nil {
		return 0, err
	}
	err = stmt.BindText(3, stats)
	if err != nil {
		return 0, err
	}

	err = stmt.Exec()
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
	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.writePool.Put(conn)

	postStmt, _, err := conn.PrepareAndPersist(SQLForInsertPost)
	if err != nil {
		return 0, err
	}

	err = postStmt.BindText(1, postTitle)
	if err != nil {
		return 0, err
	}
	err = postStmt.BindText(2, postContent)
	if err != nil {
		return 0, err
	}
	err = postStmt.BindText(3, postStats)
	if err != nil {
		return 0, err
	}

	err = postStmt.Exec()
	if err != nil {
		return 0, err
	}

	postID := conn.LastInsertRowID()
	if err != nil {
		return 0, err
	}

	commentStmt, _, err := conn.PrepareAndPersist(SQLForInsertComment)
	if err != nil {
		return 0, err
	}

	for range comments {
		err = commentStmt.BindInt64(1, postID)
		if err != nil {
			return 0, err
		}
		err = commentStmt.BindText(2, commentName)
		if err != nil {
			return 0, err
		}
		err = commentStmt.BindText(3, commentContent)
		if err != nil {
			return 0, err
		}
		err = commentStmt.BindText(4, commentStats)
		if err != nil {
			return 0, err
		}

		err = commentStmt.Exec()
		if err != nil {
			return 0, err
		}
	}

	return postID, nil
}

func (db *DB) WritePostAndCommentsWithTx(ctx context.Context, postTitle, postContent, postStats string, comments int, commentName, commentContent, commentStats string) (int64, error) {
	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.writePool.Put(conn)

	tx, err := conn.BeginImmediate()
	if err != nil {
		return 0, err
	}
	defer tx.End(&err)

	postStmt, _, err := conn.PrepareAndPersist(SQLForInsertPost)
	if err != nil {
		return 0, err
	}

	err = postStmt.BindText(1, postTitle)
	if err != nil {
		return 0, err
	}
	err = postStmt.BindText(2, postContent)
	if err != nil {
		return 0, err
	}
	err = postStmt.BindText(3, postStats)
	if err != nil {
		return 0, err
	}

	err = postStmt.Exec()
	if err != nil {
		return 0, err
	}

	postID := conn.LastInsertRowID()
	if err != nil {
		return 0, err
	}

	commentStmt, _, err := conn.PrepareAndPersist(SQLForInsertComment)
	if err != nil {
		return 0, err
	}

	for range comments {
		err = commentStmt.BindInt64(1, postID)
		if err != nil {
			return 0, err
		}
		err = commentStmt.BindText(2, commentName)
		if err != nil {
			return 0, err
		}
		err = commentStmt.BindText(3, commentContent)
		if err != nil {
			return 0, err
		}
		err = commentStmt.BindText(4, commentStats)
		if err != nil {
			return 0, err
		}

		err = commentStmt.Exec()
		if err != nil {
			return 0, err
		}
	}

	return postID, nil
}

// ===

func (db *DB) QueryCorrelated(ctx context.Context) (int, error) {
	n := 0

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareAndPersist(SQLForQueryCorrelated)
	if err != nil {
		return 0, err
	}

	for stmt.Step() {
		n += 1
	}

	err = stmt.Reset()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) QueryGroupBy(ctx context.Context) (int, error) {
	n := 0

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareAndPersist(SQLForQueryGroupBy)
	if err != nil {
		return 0, err
	}

	for stmt.Step() {
		n += 1
	}

	err = stmt.Reset()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) QueryJSON(ctx context.Context) (int, error) {
	n := 0

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareAndPersist(SQLForQueryJSON)
	if err != nil {
		return 0, err
	}

	for stmt.Step() {
		n += 1
	}

	err = stmt.Reset()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) QueryNonRecursiveCTE(ctx context.Context) (int, error) {
	n := 0

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareAndPersist(SQLForQueryNonRecursiveCTE)
	if err != nil {
		return 0, err
	}

	for stmt.Step() {
		n += 1
	}

	err = stmt.Reset()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) QueryOrderBy(ctx context.Context) (int, error) {
	n := 0

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareAndPersist(SQLForQueryOrderBy)
	if err != nil {
		return 0, err
	}

	for stmt.Step() {
		n += 1
	}

	err = stmt.Reset()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) QueryRecursiveCTE(ctx context.Context) (int, error) {
	n := 0

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareAndPersist(SQLForQueryRecursiveCTE)
	if err != nil {
		return 0, err
	}

	for stmt.Step() {
		n += 1
	}

	err = stmt.Reset()
	if err != nil {
		return 0, err
	}

	return n, nil
}

// ===

func (db *DB) Analyze(ctx context.Context) error {
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.Prepare("ANALYZE")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Conn(ctx context.Context) error {
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return err
	}
	defer db.readPool.Put(conn)

	return nil
}

func (db *DB) Options(ctx context.Context) ([]string, error) {
	options := make([]string, 0)

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.Prepare("PRAGMA compile_options")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for stmt.Step() {
		option := stmt.ColumnText(0)

		options = append(options, option)
	}

	err = stmt.Reset()
	if err != nil {
		return nil, err
	}

	return options, nil
}

func (db *DB) Pragmas(ctx context.Context, names []string) ([]string, error) {
	pragmas := make([]string, 0)

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer db.readPool.Put(conn)

	pragma := func(name string) error {
		stmt, _, err := conn.Prepare("PRAGMA" + " " + name)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for stmt.Step() {
			value := stmt.ColumnText(0)

			pragmas = append(pragmas, fmt.Sprintf("%s=%s", name, value))
		}

		err = stmt.Reset()
		if err != nil {
			return err
		}

		return nil
	}

	for _, name := range names {
		err = pragma(name)
		if err != nil {
			return nil, err
		}
	}

	return pragmas, nil
}

func (db *DB) Select1(ctx context.Context) error {
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.Prepare("SELECT 1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Select1PrePrepared(ctx context.Context) error {
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.PrepareAndPersist("SELECT 1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Time(ctx context.Context, in time.Time) (time.Time, error) {
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return time.Time{}, err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.Prepare("SELECT ?")
	if err != nil {
		return time.Time{}, err
	}
	defer stmt.Close()

	err = stmt.BindTime(1, in, sqlite3.TimeFormatDefault)
	if err != nil {
		return time.Time{}, err
	}

	var out time.Time

	for stmt.Step() {
		out = stmt.ColumnTime(0, sqlite3.TimeFormatDefault)
	}

	err = stmt.Reset()
	if err != nil {
		return time.Time{}, err
	}

	return out, nil
}

func (db *DB) Version(ctx context.Context) (string, error) {
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return "", err
	}
	defer db.readPool.Put(conn)

	stmt, _, err := conn.Prepare("SELECT sqlite_version()")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var version string

	for stmt.Step() {
		version = stmt.ColumnText(0)
	}

	err = stmt.Reset()
	if err != nil {
		return "", err
	}

	return version, nil
}

func (c *Conn) PrepareAndPersist(sql string) (*sqlite3.Stmt, string, error) {
	if stmt, ok := c.prepared[sql]; ok {
		err := stmt.Reset()
		if err != nil {
			return nil, "", err
		}

		err = stmt.ClearBindings()
		if err != nil {
			return nil, "", err
		}

		return stmt, "", nil
	}

	stmt, tail, err := c.PrepareFlags(sql, sqlite3.PREPARE_PERSISTENT)
	if err != nil {
		return nil, "", err
	}
	if c.prepared == nil {
		c.prepared = make(map[string]*sqlite3.Stmt)
	}
	c.prepared[sql] = stmt

	return stmt, tail, nil
}

//go:build eatonphil_direct

package go_sqlite_bench

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/eatonphil/gosqlite"
)

type DB struct {
	readPool  *Pool[*Conn]
	writePool *Pool[*Conn]
}

type Conn struct {
	*gosqlite.Conn

	prepared map[string]*gosqlite.Stmt
}

// ===

func NewDB(ctx context.Context, filename string, maxReadConnections, maxWriteConnections int) (*DB, error) {
	if !(maxReadConnections >= 0) {
		return nil, errors.New("maxReadConnections must be >= 0")
	}
	if !(maxWriteConnections >= 1) {
		return nil, errors.New("maxWriteConnections must be >= 1")
	}

	if maxReadConnections == 0 {
		pool, err := newPool(filename, maxWriteConnections, maxWriteConnections, 0)
		if err != nil {
			return nil, err
		}

		return &DB{readPool: pool, writePool: pool}, nil
	} else {
		readPool, err := newPool(filename, maxReadConnections, maxReadConnections, 0)
		if err != nil {
			return nil, err
		}
		writePool, err := newPool(filename, maxWriteConnections, maxWriteConnections, 0)
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
			conn, err := gosqlite.Open("file:"+filename, gosqlite.OPEN_CREATE|gosqlite.OPEN_READWRITE|gosqlite.OPEN_URI|gosqlite.OPEN_NOMUTEX)
			if err != nil {
				return nil, err
			}

			err = conn.Exec(`
				PRAGMA busy_timeout(10000);
				PRAGMA foreign_keys(true);
				PRAGMA journal_mode(wal);
				PRAGMA synchronous(normal);
			`)
			if err != nil {
				conn.Close()
				return nil, err
			}

			return &Conn{Conn: conn}, nil
		}, func(conn *Conn) error {
			if !conn.AutoCommit() {
				return conn.Rollback()
			}

			return nil
		}, func(conn *Conn) {
			conn.Close()
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

	err = conn.WithTxImmediate(func() error {
		for _, s := range SQLForSchema {
			err := conn.Exec(s)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) PopulateDB(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) error {
	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return err
	}
	defer db.writePool.Put(conn)

	postStmt, err := conn.PrepareAndPersist(SQLForInsertPostWithCreated)
	if err != nil {
		return err
	}
	defer postStmt.Reset() // Purely defensive.

	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON
	postDate := NewPostDate(posts)

	commentStmt, err := conn.PrepareAndPersist(SQLForInsertCommentWithCreated)
	if err != nil {
		return err
	}
	defer commentStmt.Reset() // Purely defensive.

	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON
	commentDate := postDate.CommentDate

	for i := range posts {
		postSeq := i + 1
		title := fmt.Sprintf("Post %d", postSeq)
		content := postContent
		stats := postStats
		created := postDate.NextFormatted()

		err = postStmt.Exec(title, content, stats, created)
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

			err = commentStmt.Exec(postID, name, content, stats, created)
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

	err = conn.WithTxImmediate(func() error {
		postStmt, err := conn.PrepareAndPersist(SQLForInsertPostWithCreated)
		if err != nil {
			return err
		}
		defer postStmt.Reset() // Purely defensive.

		postContent := Paragraphs(LoremIpsum, postParagraphs)
		postStats := LoremIpsumJSON
		postDate := NewPostDate(posts)

		commentStmt, err := conn.PrepareAndPersist(SQLForInsertCommentWithCreated)
		if err != nil {
			return err
		}
		defer commentStmt.Reset() // Purely defensive.

		commentContent := Paragraphs(LoremIpsum, commentParagraphs)
		commentStats := LoremIpsumJSON
		commentDate := postDate.CommentDate

		for i := range posts {
			postSeq := i + 1
			title := fmt.Sprintf("Post %d", postSeq)
			content := postContent
			stats := postStats
			created := postDate.NextFormatted()

			err = postStmt.Exec(title, content, stats, created)
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

				err = commentStmt.Exec(postID, name, content, stats, created)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) PopulateDBWithTxs(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) error {
	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return err
	}
	defer db.writePool.Put(conn)

	postStmt, err := conn.PrepareAndPersist(SQLForInsertPostWithCreated)
	if err != nil {
		return err
	}
	defer postStmt.Reset() // Purely defensive.

	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON
	postDate := NewPostDate(posts)

	commentStmt, err := conn.PrepareAndPersist(SQLForInsertCommentWithCreated)
	if err != nil {
		return err
	}
	defer commentStmt.Reset() // Purely defensive.

	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON
	commentDate := postDate.CommentDate

	for i := range posts {
		err = conn.WithTxImmediate(func() error {
			postSeq := i + 1
			title := fmt.Sprintf("Post %d", postSeq)
			content := postContent
			stats := postStats
			created := postDate.NextFormatted()

			err = postStmt.Exec(title, content, stats, created)
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

				err = commentStmt.Exec(postID, name, content, stats, created)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) CountPosts(ctx context.Context) (int64, error) {
	n := int64(0)

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, err := conn.PrepareAndPersist(SQLForCountPosts)
	if err != nil {
		return 0, err
	}
	defer stmt.Reset() // Purely defensive.

	row, err := stmt.Step()
	for row {
		err = stmt.Scan(&n)
		if err != nil {
			return 0, err
		}

		row, err = stmt.Step()
	}
	if err != nil {
		return 0, err
	}

	err = stmt.Reset()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) CountComments(ctx context.Context) (int64, error) {
	n := int64(0)

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, err := conn.PrepareAndPersist(SQLForCountComments)
	if err != nil {
		return 0, err
	}
	defer stmt.Reset() // Purely defensive.

	row, err := stmt.Step()
	for row {
		err = stmt.Scan(&n)
		if err != nil {
			return 0, err
		}

		row, err = stmt.Step()
	}
	if err != nil {
		return 0, err
	}

	err = stmt.Reset()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) ReadPost(ctx context.Context, id int64) (*Post, error) {
	p := &Post{ID: id}

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer db.readPool.Put(conn)

	stmt, err := conn.PrepareAndPersist(SQLForSelectPostByID)
	if err != nil {
		return nil, err
	}
	defer stmt.Reset() // Purely defensive.

	err = stmt.Bind(id)
	if err != nil {
		return nil, err
	}

	row, err := stmt.Step()
	for row {
		err = stmt.Scan(&p.Title, &p.Content, &p.Created, &p.Stats)
		if err != nil {
			return nil, err
		}

		row, err = stmt.Step()
	}
	if err != nil {
		return nil, err
	}

	err = stmt.Reset()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *DB) ReadPostWithTx(ctx context.Context, id int64) (*Post, error) {
	p := &Post{ID: id}

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer db.readPool.Put(conn)

	err = conn.WithTx(func() error {
		stmt, err := conn.PrepareAndPersist(SQLForSelectPostByID)
		if err != nil {
			return err
		}
		defer stmt.Reset() // Purely defensive.

		err = stmt.Bind(id)
		if err != nil {
			return err
		}

		row, err := stmt.Step()
		for row {
			err = stmt.Scan(&p.Title, &p.Content, &p.Created, &p.Stats)
			if err != nil {
				return err
			}

			row, err = stmt.Step()
		}
		if err != nil {
			return err
		}

		err = stmt.Reset()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *DB) ReadPostAndComments(ctx context.Context, id int64) (*Post, []*Comment, error) {
	p := &Post{ID: id}
	cs := make([]*Comment, 0)

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer db.readPool.Put(conn)

	stmt, err := conn.PrepareAndPersist(SQLForSelectPostByID)
	if err != nil {
		return nil, nil, err
	}
	defer stmt.Reset() // Purely defensive.

	err = stmt.Bind(id)
	if err != nil {
		return nil, nil, err
	}

	row, err := stmt.Step()
	for row {
		err = stmt.Scan(&p.Title, &p.Content, &p.Created, &p.Stats)
		if err != nil {
			return nil, nil, err
		}

		row, err = stmt.Step()
	}
	if err != nil {
		return nil, nil, err
	}

	err = stmt.Reset()
	if err != nil {
		return nil, nil, err
	}

	stmt, err = conn.PrepareAndPersist(SQLForSelectCommentsByPostID)
	if err != nil {
		return nil, nil, err
	}
	defer stmt.Reset() // Purely defensive.

	err = stmt.Bind(id)
	if err != nil {
		return nil, nil, err
	}

	row, err = stmt.Step()
	for row {
		c := &Comment{}
		err = stmt.Scan(&c.ID, &c.Name, &c.Content, &c.Created, &c.Stats)
		if err != nil {
			return nil, nil, err
		}
		cs = append(cs, c)

		row, err = stmt.Step()
	}
	if err != nil {
		return nil, nil, err
	}

	err = stmt.Reset()
	if err != nil {
		return nil, nil, err
	}

	return p, cs, nil
}

func (db *DB) ReadPostAndCommentsWithTx(ctx context.Context, id int64) (*Post, []*Comment, error) {
	p := &Post{ID: id}
	cs := make([]*Comment, 0)

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer db.readPool.Put(conn)

	err = conn.WithTx(func() error {
		stmt, err := conn.PrepareAndPersist(SQLForSelectPostByID)
		if err != nil {
			return err
		}
		defer stmt.Reset() // Purely defensive.

		err = stmt.Bind(id)
		if err != nil {
			return err
		}

		row, err := stmt.Step()
		for row {
			err = stmt.Scan(&p.Title, &p.Content, &p.Created, &p.Stats)
			if err != nil {
				return err
			}

			row, err = stmt.Step()
		}
		if err != nil {
			return err
		}

		err = stmt.Reset()
		if err != nil {
			return err
		}

		stmt, err = conn.PrepareAndPersist(SQLForSelectCommentsByPostID)
		if err != nil {
			return err
		}
		defer stmt.Reset() // Purely defensive.

		err = stmt.Bind(id)
		if err != nil {
			return err
		}

		row, err = stmt.Step()
		for row {
			c := &Comment{}
			err = stmt.Scan(&c.ID, &c.Name, &c.Content, &c.Created, &c.Stats)
			if err != nil {
				return err
			}
			cs = append(cs, c)

			row, err = stmt.Step()
		}
		if err != nil {
			return err
		}

		err = stmt.Reset()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return p, cs, nil
}

func (db *DB) WritePost(ctx context.Context, title, content, stats string) (int64, error) {
	postID := int64(0)

	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.writePool.Put(conn)

	stmt, err := conn.PrepareAndPersist(SQLForInsertPost)
	if err != nil {
		return 0, err
	}
	defer stmt.Reset() // Purely defensive.

	err = stmt.Exec(title, content, stats)
	if err != nil {
		return 0, err
	}

	postID = conn.LastInsertRowID()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (db *DB) WritePostWithTx(ctx context.Context, title, content, stats string) (int64, error) {
	postID := int64(0)

	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.writePool.Put(conn)

	err = conn.WithTxImmediate(func() error {
		stmt, err := conn.PrepareAndPersist(SQLForInsertPost)
		if err != nil {
			return err
		}
		defer stmt.Reset() // Purely defensive.

		err = stmt.Exec(title, content, stats)
		if err != nil {
			return err
		}

		postID = conn.LastInsertRowID()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (db *DB) WritePostAndComments(ctx context.Context, postTitle, postContent, postStats string, comments int, commentName, commentContent, commentStats string) (int64, error) {
	postID := int64(0)

	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.writePool.Put(conn)

	postStmt, err := conn.PrepareAndPersist(SQLForInsertPost)
	if err != nil {
		return 0, err
	}
	defer postStmt.Reset() // Purely defensive.

	err = postStmt.Exec(postTitle, postContent, postStats)
	if err != nil {
		return 0, err
	}

	postID = conn.LastInsertRowID()
	if err != nil {
		return 0, err
	}

	commentStmt, err := conn.PrepareAndPersist(SQLForInsertComment)
	if err != nil {
		return 0, err
	}
	defer commentStmt.Reset() // Purely defensive.

	for range comments {
		err = commentStmt.Exec(postID, commentName, commentContent, commentStats)
		if err != nil {
			return 0, err
		}
	}

	return postID, nil
}

func (db *DB) WritePostAndCommentsWithTx(ctx context.Context, postTitle, postContent, postStats string, comments int, commentName, commentContent, commentStats string) (int64, error) {
	postID := int64(0)

	conn, err := db.writePool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.writePool.Put(conn)

	err = conn.WithTxImmediate(func() error {
		postStmt, err := conn.PrepareAndPersist(SQLForInsertPost)
		if err != nil {
			return err
		}
		defer postStmt.Reset() // Purely defensive.

		err = postStmt.Exec(postTitle, postContent, postStats)
		if err != nil {
			return err
		}

		postID = conn.LastInsertRowID()
		if err != nil {
			return err
		}

		commentStmt, err := conn.PrepareAndPersist(SQLForInsertComment)
		if err != nil {
			return err
		}
		defer commentStmt.Reset() // Purely defensive.

		for range comments {
			err = commentStmt.Exec(postID, commentName, commentContent, commentStats)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return postID, nil
}

// ===

func (db *DB) query(ctx context.Context, sql string) (int, error) {
	n := 0

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return 0, err
	}
	defer db.readPool.Put(conn)

	stmt, err := conn.PrepareAndPersist(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Reset() // Purely defensive.

	row, err := stmt.Step()
	for row {
		n += 1

		row, err = stmt.Step()
	}
	if err != nil {
		return 0, err
	}

	err = stmt.Reset()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) QueryCorrelated(ctx context.Context) (int, error) {
	return db.query(ctx, SQLForQueryCorrelated)
}

func (db *DB) QueryGroupBy(ctx context.Context) (int, error) {
	return db.query(ctx, SQLForQueryGroupBy)
}

func (db *DB) QueryJSON(ctx context.Context) (int, error) {
	return db.query(ctx, SQLForQueryJSON)
}

func (db *DB) QueryOrderBy(ctx context.Context) (int, error) {
	return db.query(ctx, SQLForQueryOrderBy)
}

func (db *DB) QueryRecursiveCTE(ctx context.Context) (int, error) {
	return db.query(ctx, SQLForQueryRecursiveCTE)
}

func (db *DB) QueryWindow(ctx context.Context) (int, error) {
	return db.query(ctx, SQLForQueryWindow)
}

// ===

func (db *DB) Analyze(ctx context.Context) error {
	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return err
	}
	defer db.readPool.Put(conn)

	err = conn.Exec("ANALYZE")
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

	stmt, err := conn.Prepare("PRAGMA compile_options")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row, err := stmt.Step()
	for row {
		option := ""

		err = stmt.Scan(&option)
		if err != nil {
			return nil, err
		}

		options = append(options, option)

		row, err = stmt.Step()
	}
	if err != nil {
		return nil, err
	}

	err = stmt.Reset()
	if err != nil {
		return nil, err
	}

	return options, nil
}

func (db *DB) Pragma(ctx context.Context, name string) (string, error) {
	value := ""

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return "", err
	}
	defer db.readPool.Put(conn)

	stmt, err := conn.Prepare("PRAGMA" + " " + name)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	row, err := stmt.Step()
	for row {
		err = stmt.Scan(&value)
		if err != nil {
			return "", err
		}

		row, err = stmt.Step()
	}
	if err != nil {
		return "", err
	}

	err = stmt.Reset()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (db *DB) Pragmas(ctx context.Context, names []string) ([]string, error) {
	pragmas := make([]string, 0)

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer db.readPool.Put(conn)

	pragma := func(name string) error {
		stmt, err := conn.Prepare("PRAGMA" + " " + name)
		if err != nil {
			return err
		}
		defer stmt.Close()

		row, err := stmt.Step()
		for row {
			value := ""

			err = stmt.Scan(&value)
			if err != nil {
				return err
			}

			pragmas = append(pragmas, fmt.Sprintf("%s=%s", name, value))

			row, err = stmt.Step()
		}
		if err != nil {
			return err
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

	stmt, err := conn.Prepare("SELECT 1")
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

	stmt, err := conn.PrepareAndPersist("SELECT 1")
	if err != nil {
		return err
	}
	defer stmt.Reset() // Purely defensive.

	err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Time(ctx context.Context, in time.Time) (time.Time, error) {
	return time.Time{}, errors.New("no built-in support for time")
}

func (db *DB) Version(ctx context.Context) (string, error) {
	version := ""

	conn, err := db.readPool.Get(ctx)
	if err != nil {
		return "", err
	}
	defer db.readPool.Put(conn)

	stmt, err := conn.Prepare("SELECT sqlite_version()")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	row, err := stmt.Step()
	for row {
		err = stmt.Scan(&version)
		if err != nil {
			return "", err
		}

		row, err = stmt.Step()
	}
	if err != nil {
		return "", err
	}

	err = stmt.Reset()
	if err != nil {
		return "", err
	}

	return version, nil
}

// ===

func (c *Conn) Close() error {
	for name, stmt := range c.prepared {
		err := stmt.Close()
		if err != nil {
			log.Printf("failed to close: %v", err)
		}

		delete(c.prepared, name)
	}

	return c.Conn.Close()
}

func (c *Conn) PrepareAndPersist(sql string) (*gosqlite.Stmt, error) {
	if stmt, ok := c.prepared[sql]; ok {
		err := stmt.ClearBindings()
		if err != nil {
			return nil, err
		}

		err = stmt.Reset()
		if err != nil {
			return nil, err
		}

		return stmt, nil
	}

	stmt, err := c.Prepare(sql)
	if err != nil {
		return nil, err
	}
	if c.prepared == nil {
		c.prepared = make(map[string]*gosqlite.Stmt)
	}
	c.prepared[sql] = stmt

	return stmt, nil
}

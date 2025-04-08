//go:build glebarez_driver || mattn_driver || modernc_driver || ncruces_driver || tailscale_driver

package go_sqlite_bench

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime"
	"time"
)

type DB struct {
	readDB  *sql.DB
	writeDB *sql.DB

	// Non-Schema ("NewDB") Statements

	readDBSelect1 *sql.Stmt

	// Schema ("PrepareDB") Statements

	writeDBInsertPost               *sql.Stmt
	writeDBInsertPostWithCreated    *sql.Stmt
	writeDBInsertComment            *sql.Stmt
	writeDBInsertCommentWithCreated *sql.Stmt

	readDBSelectPostByID         *sql.Stmt
	readDBSelectCommentsByPostID *sql.Stmt

	readDBQueryCorrelated   *sql.Stmt
	readDBQueryGroupBy      *sql.Stmt
	readDBQueryJSON         *sql.Stmt
	readDBQueryOrderBy      *sql.Stmt
	readDBQueryRecursiveCTE *sql.Stmt
	readDBQueryWindow       *sql.Stmt
}

// ===

func NewDB(ctx context.Context, filename string, maxReadConnections, maxWriteConnections int) (*DB, error) {
	var db *DB

	if !(maxReadConnections >= 0) {
		return nil, errors.New("maxReadConnections must be >= 0")
	}
	if !(maxWriteConnections >= 1) {
		return nil, errors.New("maxWriteConnections must be >= 1")
	}

	if maxReadConnections == 0 {
		readWriteDB, err := OpenDB(filename)
		if err != nil {
			return nil, err
		}

		readWriteDB.SetMaxIdleConns(runtime.GOMAXPROCS(0))
		readWriteDB.SetMaxOpenConns(maxWriteConnections)

		db = &DB{readDB: readWriteDB, writeDB: readWriteDB}
	} else {
		readDB, err := OpenDB(filename)
		if err != nil {
			return nil, err
		}

		readDB.SetMaxIdleConns(runtime.GOMAXPROCS(0))
		readDB.SetMaxOpenConns(maxReadConnections)

		writeDB, err := OpenDB(filename)
		if err != nil {
			err = errors.Join(readDB.Close(), err)
			return nil, err
		}

		writeDB.SetMaxIdleConns(runtime.GOMAXPROCS(0))
		writeDB.SetMaxOpenConns(maxWriteConnections)

		db = &DB{readDB: readDB, writeDB: writeDB}
	}

	err := db.prepareNewDBStatements()
	if err != nil {
		err = errors.Join(db.Close(), err)
		return nil, err
	}

	return db, nil
}

func (db *DB) prepareNewDBStatements() error {
	stmt, err := db.readDB.Prepare("SELECT 1")
	if err != nil {
		return err
	}
	db.readDBSelect1 = stmt

	return nil
}

func (db *DB) closeNewDBStatements() error {
	var err error

	if db.readDBSelect1 != nil {
		err = errors.Join(db.readDBSelect1.Close(), err)
		db.readDBSelect1 = nil
	}

	return err
}

func (db *DB) Close() error {
	var err error

	err = errors.Join(db.closePrepareDBStatements(), err)
	err = errors.Join(db.closeNewDBStatements(), err)

	err = errors.Join(db.writeDB.Close(), err)
	if db.readDB != db.writeDB {
		err = errors.Join(db.readDB.Close(), err)
	}

	return err
}

func (db *DB) PrepareDBWithTx(ctx context.Context) error {
	tx, err := db.writeDB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, s := range SQLForSchema {
		_, err = tx.ExecContext(ctx, s)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	err = db.preparePrepareDBStatements()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) preparePrepareDBStatements() error {
	stmt, err := db.writeDB.Prepare(SQLForInsertPost)
	if err != nil {
		return err
	}
	db.writeDBInsertPost = stmt

	stmt, err = db.writeDB.Prepare(SQLForInsertPostWithCreated)
	if err != nil {
		return err
	}
	db.writeDBInsertPostWithCreated = stmt

	stmt, err = db.writeDB.Prepare(SQLForInsertComment)
	if err != nil {
		return err
	}
	db.writeDBInsertComment = stmt

	stmt, err = db.writeDB.Prepare(SQLForInsertCommentWithCreated)
	if err != nil {
		return err
	}
	db.writeDBInsertCommentWithCreated = stmt

	stmt, err = db.readDB.Prepare(SQLForSelectPostByID)
	if err != nil {
		return err
	}
	db.readDBSelectPostByID = stmt

	stmt, err = db.readDB.Prepare(SQLForSelectCommentsByPostID)
	if err != nil {
		return err
	}
	db.readDBSelectCommentsByPostID = stmt

	stmt, err = db.readDB.Prepare(SQLForQueryCorrelated)
	if err != nil {
		return err
	}
	db.readDBQueryCorrelated = stmt

	stmt, err = db.readDB.Prepare(SQLForQueryGroupBy)
	if err != nil {
		return err
	}
	db.readDBQueryGroupBy = stmt

	stmt, err = db.readDB.Prepare(SQLForQueryJSON)
	if err != nil {
		return err
	}
	db.readDBQueryJSON = stmt

	stmt, err = db.readDB.Prepare(SQLForQueryOrderBy)
	if err != nil {
		return err
	}
	db.readDBQueryOrderBy = stmt

	stmt, err = db.readDB.Prepare(SQLForQueryRecursiveCTE)
	if err != nil {
		return err
	}
	db.readDBQueryRecursiveCTE = stmt

	stmt, err = db.readDB.Prepare(SQLForQueryWindow)
	if err != nil {
		return err
	}
	db.readDBQueryWindow = stmt

	return nil
}

func (db *DB) closePrepareDBStatements() error {
	var err error

	if db.writeDBInsertPost != nil {
		err = errors.Join(db.writeDBInsertPost.Close(), err)
		db.writeDBInsertPost = nil
	}

	if db.writeDBInsertPostWithCreated != nil {
		err = errors.Join(db.writeDBInsertPostWithCreated.Close(), err)
		db.writeDBInsertPostWithCreated = nil
	}

	if db.writeDBInsertComment != nil {
		err = errors.Join(db.writeDBInsertComment.Close(), err)
		db.writeDBInsertComment = nil
	}

	if db.writeDBInsertCommentWithCreated != nil {
		err = errors.Join(db.writeDBInsertCommentWithCreated.Close(), err)
		db.writeDBInsertCommentWithCreated = nil
	}

	if db.readDBSelectPostByID != nil {
		err = errors.Join(db.readDBSelectPostByID.Close(), err)
		db.readDBSelectPostByID = nil
	}

	if db.readDBSelectCommentsByPostID != nil {
		err = errors.Join(db.readDBSelectCommentsByPostID.Close(), err)
		db.readDBSelectCommentsByPostID = nil
	}

	if db.readDBQueryCorrelated != nil {
		err = errors.Join(db.readDBQueryCorrelated.Close(), err)
		db.readDBQueryCorrelated = nil
	}

	if db.readDBQueryGroupBy != nil {
		err = errors.Join(db.readDBQueryGroupBy.Close(), err)
		db.readDBQueryGroupBy = nil
	}

	if db.readDBQueryJSON != nil {
		err = errors.Join(db.readDBQueryJSON.Close(), err)
		db.readDBQueryJSON = nil
	}

	if db.readDBQueryOrderBy != nil {
		err = errors.Join(db.readDBQueryOrderBy.Close(), err)
		db.readDBQueryOrderBy = nil
	}

	if db.readDBQueryRecursiveCTE != nil {
		err = errors.Join(db.readDBQueryRecursiveCTE.Close(), err)
		db.readDBQueryRecursiveCTE = nil
	}

	if db.readDBQueryWindow != nil {
		err = errors.Join(db.readDBQueryWindow.Close(), err)
		db.readDBQueryWindow = nil
	}

	return err
}

func (db *DB) PopulateDB(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) error {
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

		result, err := db.writeDBInsertPostWithCreated.ExecContext(ctx, title, content, stats, created)
		if err != nil {
			return err
		}

		postID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		for j := range comments {
			commentSeq := j + 1
			name := fmt.Sprintf("Comment %d.%d", postSeq, commentSeq)
			content := commentContent
			stats := commentStats
			created := commentDate.NextFormatted()

			_, err = db.writeDBInsertCommentWithCreated.ExecContext(ctx, postID, name, content, stats, created)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *DB) PopulateDBWithTx(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) error {
	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON
	postDate := NewPostDate(posts)

	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON
	commentDate := postDate.CommentDate

	tx, err := db.writeDB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for i := range posts {
		postSeq := i + 1
		title := fmt.Sprintf("Post %d", postSeq)
		content := postContent
		stats := postStats
		created := postDate.NextFormatted()

		result, err := tx.Stmt(db.writeDBInsertPostWithCreated).ExecContext(ctx, title, content, stats, created)
		if err != nil {
			return err
		}

		postID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		for j := range comments {
			commentSeq := j + 1
			name := fmt.Sprintf("Comment %d.%d", postSeq, commentSeq)
			content := commentContent
			stats := commentStats
			created := commentDate.NextFormatted()

			_, err = tx.Stmt(db.writeDBInsertCommentWithCreated).ExecContext(ctx, postID, name, content, stats, created)
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) PopulateDBWithTxs(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) error {
	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON
	postDate := NewPostDate(posts)

	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON
	commentDate := postDate.CommentDate

	writePostAndComments := func(i int) (err error) {
		tx, err := db.writeDB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()

		postSeq := i + 1
		title := fmt.Sprintf("Post %d", postSeq)
		content := postContent
		stats := postStats
		created := postDate.NextFormatted()

		result, err := tx.Stmt(db.writeDBInsertPostWithCreated).ExecContext(ctx, title, content, stats, created)
		if err != nil {
			return err
		}

		postID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		for j := range comments {
			commentSeq := j + 1
			name := fmt.Sprintf("Comment %d.%d", postSeq, commentSeq)
			content := commentContent
			stats := commentStats
			created := commentDate.NextFormatted()

			_, err = tx.Stmt(db.writeDBInsertCommentWithCreated).ExecContext(ctx, postID, name, content, stats, created)
			if err != nil {
				return err
			}
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		return nil
	}

	for i := range posts {
		err := writePostAndComments(i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) CountPosts(ctx context.Context) (int64, error) {
	n := int64(0)

	row := db.readDB.QueryRowContext(ctx, SQLForCountPosts)
	err := row.Scan(&n)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) CountComments(ctx context.Context) (int64, error) {
	n := int64(0)

	row := db.readDB.QueryRowContext(ctx, SQLForCountComments)
	err := row.Scan(&n)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) ReadPost(ctx context.Context, id int64) (*Post, error) {
	p := &Post{ID: id}

	row := db.readDBSelectPostByID.QueryRowContext(ctx, id)
	err := row.Scan(&p.Title, &p.Content, &p.Created, &p.Stats)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *DB) ReadPostWithTx(ctx context.Context, id int64) (*Post, error) {
	p := &Post{ID: id}

	tx, err := db.readDB.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	row := tx.Stmt(db.readDBSelectPostByID).QueryRowContext(ctx, id)
	err = row.Scan(&p.Title, &p.Content, &p.Created, &p.Stats)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p = nil
		}
		return p, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *DB) ReadPostAndComments(ctx context.Context, id int64) (*Post, []*Comment, error) {
	p := &Post{ID: id}
	cs := make([]*Comment, 0)

	row := db.readDBSelectPostByID.QueryRowContext(ctx, id)
	err := row.Scan(&p.Title, &p.Content, &p.Created, &p.Stats)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p = nil
		}
		return nil, nil, err
	}

	rows, err := db.readDBSelectCommentsByPostID.QueryContext(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := &Comment{}

		err = rows.Scan(&c.ID, &c.Name, &c.Content, &c.Created, &c.Stats)
		if err != nil {
			return nil, nil, err
		}

		cs = append(cs, c)
	}

	err = rows.Err()
	if err != nil {
		return nil, nil, err
	}

	return p, cs, nil
}

func (db *DB) ReadPostAndCommentsWithTx(ctx context.Context, id int64) (*Post, []*Comment, error) {
	p := &Post{ID: id}
	cs := make([]*Comment, 0)

	tx, err := db.readDB.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	row := tx.Stmt(db.readDBSelectPostByID).QueryRowContext(ctx, id)
	err = row.Scan(&p.Title, &p.Content, &p.Created, &p.Stats)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p = nil
		}
		return nil, nil, err
	}

	rows, err := tx.Stmt(db.readDBSelectCommentsByPostID).QueryContext(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := &Comment{}

		err = rows.Scan(&c.ID, &c.Name, &c.Content, &c.Created, &c.Stats)
		if err != nil {
			return nil, nil, err
		}

		cs = append(cs, c)
	}

	err = rows.Err()
	if err != nil {
		return nil, nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, nil, err
	}

	return p, cs, nil
}

func (db *DB) WritePost(ctx context.Context, title, content, stats string) (int64, error) {
	postID := int64(0)

	result, err := db.writeDBInsertPost.ExecContext(ctx, title, content, stats)
	if err != nil {
		return 0, err
	}

	postID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (db *DB) WritePostWithTx(ctx context.Context, title, content, stats string) (int64, error) {
	postID := int64(0)

	tx, err := db.writeDB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err := tx.Stmt(db.writeDBInsertPost).ExecContext(ctx, title, content, stats)
	if err != nil {
		return 0, err
	}

	postID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (db *DB) WritePostAndComments(ctx context.Context, postTitle, postContent, postStats string, comments int, commentName, commentContent, commentStats string) (int64, error) {
	postID := int64(0)

	result, err := db.writeDBInsertPost.ExecContext(ctx, postTitle, postContent, postStats)
	if err != nil {
		return 0, err
	}

	postID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for range comments {
		_, err = db.writeDBInsertComment.ExecContext(ctx, postID, commentName, commentContent, commentStats)
		if err != nil {
			return 0, err
		}
	}

	return postID, nil
}

func (db *DB) WritePostAndCommentsWithTx(ctx context.Context, postTitle, postContent, postStats string, comments int, commentName, commentContent, commentStats string) (int64, error) {
	postID := int64(0)

	tx, err := db.writeDB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err := tx.Stmt(db.writeDBInsertPost).ExecContext(ctx, postTitle, postContent, postStats)
	if err != nil {
		return 0, err
	}

	postID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for range comments {
		_, err = tx.Stmt(db.writeDBInsertComment).ExecContext(ctx, postID, commentName, commentContent, commentStats)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

// ===

func (db *DB) query(ctx context.Context, stmt *sql.Stmt) (int, error) {
	n := 0

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		n += 1
	}

	err = rows.Err()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (db *DB) QueryCorrelated(ctx context.Context) (int, error) {
	return db.query(ctx, db.readDBQueryCorrelated)
}

func (db *DB) QueryGroupBy(ctx context.Context) (int, error) {
	return db.query(ctx, db.readDBQueryGroupBy)
}

func (db *DB) QueryJSON(ctx context.Context) (int, error) {
	return db.query(ctx, db.readDBQueryJSON)
}

func (db *DB) QueryOrderBy(ctx context.Context) (int, error) {
	return db.query(ctx, db.readDBQueryOrderBy)
}

func (db *DB) QueryRecursiveCTE(ctx context.Context) (int, error) {
	return db.query(ctx, db.readDBQueryRecursiveCTE)
}

func (db *DB) QueryWindow(ctx context.Context) (int, error) {
	return db.query(ctx, db.readDBQueryWindow)
}

// ===

func (db *DB) Analyze(ctx context.Context) error {
	_, err := db.readDB.ExecContext(ctx, "ANALYZE")
	return err
}

func (db *DB) Conn(ctx context.Context) error {
	conn, err := db.readDB.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	return err
}

func (db *DB) Explain(ctx context.Context, sql string) ([]Explain, error) {
	explains := make([]Explain, 0)

	rows, err := db.readDB.QueryContext(ctx, "EXPLAIN "+sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var addr int64
		var opcode string
		var p1 *string
		var p2 *string
		var p3 *string
		var p4 *string
		var p5 *string
		var comment *string

		err = rows.Scan(
			&addr,
			&opcode,
			&p1,
			&p2,
			&p3,
			&p4,
			&p5,
			&comment,
		)
		if err != nil {
			return nil, err
		}

		nilStr := ""
		if p1 == nil {
			p1 = &nilStr
		}
		if p2 == nil {
			p2 = &nilStr
		}
		if p3 == nil {
			p3 = &nilStr
		}
		if p4 == nil {
			p4 = &nilStr
		}
		if p5 == nil {
			p5 = &nilStr
		}
		if comment == nil {
			comment = &nilStr
		}

		explains = append(explains,
			Explain{
				Addr:    addr,
				Opcode:  opcode,
				P1:      *p1,
				P2:      *p2,
				P3:      *p3,
				P4:      *p4,
				P5:      *p5,
				Comment: *comment,
			},
		)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return explains, nil
}

func (db *DB) Options(ctx context.Context) ([]string, error) {
	options := make([]string, 0)

	rows, err := db.readDB.QueryContext(ctx, "PRAGMA compile_options")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		option := ""

		err = rows.Scan(&option)
		if err != nil {
			return nil, err
		}

		options = append(options, option)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return options, nil
}

func (db *DB) Pragma(ctx context.Context, name string) (string, error) {
	value := ""

	err := db.readDB.QueryRowContext(ctx, "PRAGMA"+" "+name).Scan(&value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (db *DB) Pragmas(ctx context.Context, names []string) ([]string, error) {
	pragmas := make([]string, 0)

	for _, name := range names {
		value := ""

		err := db.readDB.QueryRowContext(ctx, "PRAGMA"+" "+name).Scan(&value)
		if err != nil {
			return nil, err
		}

		pragmas = append(pragmas, fmt.Sprintf("%s=%s", name, value))
	}

	return pragmas, nil
}

func (db *DB) Select1(ctx context.Context) error {
	_, err := db.readDB.ExecContext(ctx, "SELECT 1")
	return err
}

func (db *DB) Select1PrePrepared(ctx context.Context) error {
	_, err := db.readDBSelect1.ExecContext(ctx)
	return err
}

func (db *DB) Time(ctx context.Context, in time.Time) (time.Time, error) {
	out := time.Time{}
	err := db.readDB.QueryRowContext(ctx, "SELECT ?", in).Scan(&out)
	return out, err
}

func (db *DB) Version(ctx context.Context) (string, error) {
	version := ""
	err := db.readDB.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&version)
	return version, err
}

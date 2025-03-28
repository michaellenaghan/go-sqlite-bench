//go:build mattn_driver || modernc_driver || ncruces_driver || tailscale_driver

package go_sqlite_bench

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	readDBQueryCorrelated      *sql.Stmt
	readDBQueryGroupBy         *sql.Stmt
	readDBQueryJSON            *sql.Stmt
	readDBQueryNonRecursiveCTE *sql.Stmt
	readDBQueryOrderBy         *sql.Stmt
	readDBQueryRecursiveCTE    *sql.Stmt
}

func NewDB(ctx context.Context, filename string, maxReadConnections, maxWriteConnections int) (*DB, error) {
	if !(maxReadConnections >= 0) {
		return nil, errors.New("maxReadConnections must be >= 0")
	}
	if !(maxWriteConnections >= 1) {
		return nil, errors.New("maxWriteConnections must be >= 1")
	}

	var db *DB

	if maxReadConnections == 0 {
		readWriteDB, err := OpenDB(filename)
		if err != nil {
			return nil, err
		}

		readWriteDB.SetMaxOpenConns(maxWriteConnections)

		db = &DB{readDB: readWriteDB, writeDB: readWriteDB}
	} else {
		readDB, err := OpenDB(filename)
		if err != nil {
			return nil, err
		}

		readDB.SetMaxOpenConns(maxReadConnections)

		writeDB, err := OpenDB(filename)
		if err != nil {
			err = errors.Join(readDB.Close(), err)
			return nil, err
		}

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
	var err error

	var stmt *sql.Stmt

	stmt, err = db.readDB.Prepare("SELECT 1")
	if err != nil {
		return err
	}
	db.readDBSelect1 = stmt

	return err
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

func (db *DB) PrepareDBWithTx(ctx context.Context) (err error) {
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

	return
}

func (db *DB) preparePrepareDBStatements() error {
	var err error

	var stmt *sql.Stmt

	stmt, err = db.writeDB.Prepare(SQLForInsertPost)
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

	stmt, err = db.readDB.Prepare(SQLForQueryNonRecursiveCTE)
	if err != nil {
		return err
	}
	db.readDBQueryNonRecursiveCTE = stmt

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

	if db.readDBQueryNonRecursiveCTE != nil {
		err = errors.Join(db.readDBQueryNonRecursiveCTE.Close(), err)
		db.readDBQueryNonRecursiveCTE = nil
	}

	if db.readDBQueryOrderBy != nil {
		err = errors.Join(db.readDBQueryOrderBy.Close(), err)
		db.readDBQueryOrderBy = nil
	}

	if db.readDBQueryRecursiveCTE != nil {
		err = errors.Join(db.readDBQueryRecursiveCTE.Close(), err)
		db.readDBQueryRecursiveCTE = nil
	}

	return err
}

func (db *DB) PopulateDB(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) (err error) {
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

	return
}

func (db *DB) PopulateDBWithTx(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) (err error) {
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

	return
}

func (db *DB) PopulateDBWithTxs(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) (err error) {
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
		err = writePostAndComments(i)
		if err != nil {
			return err
		}
	}

	return
}

func (db *DB) CountPosts(ctx context.Context) (n int64, err error) {
	row := db.readDB.QueryRowContext(ctx, SQLForCountPosts)
	if err = row.Scan(&n); err != nil {
		n = 0
		return
	}

	return
}

func (db *DB) CountComments(ctx context.Context) (n int64, err error) {
	row := db.readDB.QueryRowContext(ctx, SQLForCountComments)
	if err = row.Scan(&n); err != nil {
		n = 0
		return
	}

	return
}

func (db *DB) ReadPost(ctx context.Context, id int64) (p *Post, err error) {
	p = &Post{ID: id}

	row := db.readDBSelectPostByID.QueryRowContext(ctx, id)
	if err = row.Scan(&p.Title, &p.Content, &p.Created, &p.Stats); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p = nil
		}
		return
	}

	return
}

func (db *DB) ReadPostWithTx(ctx context.Context, id int64) (p *Post, err error) {
	tx, err := db.readDB.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	p = &Post{ID: id}

	row := tx.Stmt(db.readDBSelectPostByID).QueryRowContext(ctx, id)
	if err = row.Scan(&p.Title, &p.Content, &p.Created, &p.Stats); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p = nil
		}
		return p, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return
}

func (db *DB) ReadPostAndComments(ctx context.Context, id int64) (p *Post, cs []*Comment, err error) {
	p = &Post{ID: id}

	row := db.readDBSelectPostByID.QueryRowContext(ctx, id)
	if err = row.Scan(&p.Title, &p.Content, &p.Created, &p.Stats); err != nil {
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
		if err = rows.Scan(&c.ID, &c.Name, &c.Content, &c.Created, &c.Stats); err != nil {
			return nil, nil, err
		}
		cs = append(cs, c)
	}

	err = rows.Err()
	if err != nil {
		return nil, nil, err
	}

	return
}

func (db *DB) ReadPostAndCommentsWithTx(ctx context.Context, id int64) (p *Post, cs []*Comment, err error) {
	tx, err := db.readDB.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	p = &Post{ID: id}

	row := tx.Stmt(db.readDBSelectPostByID).QueryRowContext(ctx, id)
	if err = row.Scan(&p.Title, &p.Content, &p.Created, &p.Stats); err != nil {
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
		if err = rows.Scan(&c.ID, &c.Name, &c.Content, &c.Created, &c.Stats); err != nil {
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

	return
}

func (db *DB) WritePost(ctx context.Context, title, content, stats string) (postID int64, err error) {
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

func (db *DB) WritePostWithTx(ctx context.Context, title, content, stats string) (postID int64, err error) {
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

func (db *DB) WritePostAndComments(ctx context.Context, postTitle, postContent, postStats string, comments int, commentName, commentContent, commentStats string) (postID int64, err error) {
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

func (db *DB) WritePostAndCommentsWithTx(ctx context.Context, postTitle, postContent, postStats string, comments int, commentName, commentContent, commentStats string) (postID int64, err error) {
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

func (db *DB) QueryCorrelated(ctx context.Context) (int, error) {
	n := 0

	rows, err := db.readDBQueryCorrelated.QueryContext(ctx)
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

func (db *DB) QueryGroupBy(ctx context.Context) (int, error) {
	n := 0

	rows, err := db.readDBQueryGroupBy.QueryContext(ctx)
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

func (db *DB) QueryJSON(ctx context.Context) (int, error) {
	n := 0

	rows, err := db.readDBQueryJSON.QueryContext(ctx)
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

func (db *DB) QueryNonRecursiveCTE(ctx context.Context) (int, error) {
	n := 0

	rows, err := db.readDBQueryNonRecursiveCTE.QueryContext(ctx)
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

func (db *DB) QueryOrderBy(ctx context.Context) (int, error) {
	n := 0

	rows, err := db.readDBQueryOrderBy.QueryContext(ctx)
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

func (db *DB) QueryRecursiveCTE(ctx context.Context) (int, error) {
	n := 0

	rows, err := db.readDBQueryRecursiveCTE.QueryContext(ctx)
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

func (db *DB) Options(ctx context.Context) ([]string, error) {
	options := make([]string, 0)

	rows, err := db.readDB.QueryContext(ctx, "PRAGMA compile_options")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var option string
		if err = rows.Scan(&option); err != nil {
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
		var value string
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
	var out time.Time
	err := db.readDB.QueryRowContext(ctx, "SELECT ?", in).Scan(&out)
	return out, err
}

func (db *DB) Version(ctx context.Context) (string, error) {
	var s string
	err := db.readDB.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&s)
	return s, err
}

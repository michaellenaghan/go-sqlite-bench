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
}

func NewDB(ctx context.Context, filename string, maxReadConnections, maxWriteConnections int) (*DB, error) {
	if !(maxReadConnections >= 0) {
		return nil, errors.New("maxReadConnections must be >= 0")
	}
	if !(maxWriteConnections >= 1) {
		return nil, errors.New("maxWriteConnections must be >= 1")
	}

	if maxReadConnections == 0 {
		db, err := OpenDB(filename)
		if err != nil {
			return nil, err
		}

		db.SetMaxOpenConns(maxWriteConnections)

		return &DB{readDB: db, writeDB: db}, nil
	} else {
		readDB, err := OpenDB(filename)
		if err != nil {
			return nil, err
		}

		readDB.SetMaxOpenConns(maxReadConnections)

		writeDB, err := OpenDB(filename)
		if err != nil {
			return nil, err
		}

		writeDB.SetMaxOpenConns(maxWriteConnections)

		return &DB{readDB: readDB, writeDB: writeDB}, nil
	}
}

func (db *DB) Close() error {
	var readErr error
	var writeErr error

	readErr = db.readDB.Close()
	if db.writeDB != db.readDB {
		writeErr = db.writeDB.Close()
	}

	if readErr != nil {
		return readErr
	}
	if writeErr != nil {
		return writeErr
	}

	return nil
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

	return
}

func (db *DB) PopulateDB(ctx context.Context, posts, postParagraphs, comments, commentParagraphs int) (err error) {
	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON
	postDate := NewPostDate(posts)

	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON
	commentDate := postDate.CommentDate

	postStmt, err := db.writeDB.PrepareContext(ctx, SQLForInsertPostWithCreated)
	if err != nil {
		return err
	}
	defer postStmt.Close()

	commentStmt, err := db.writeDB.PrepareContext(ctx, SQLForInsertCommentWithCreated)
	if err != nil {
		return err
	}
	defer commentStmt.Close()

	for i := range posts {
		postSeq := i + 1
		title := fmt.Sprintf("Post %d", postSeq)
		content := postContent
		stats := postStats
		created := postDate.NextFormatted()

		result, err := postStmt.ExecContext(ctx, title, content, stats, created)
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

			_, err = commentStmt.ExecContext(ctx, postID, name, content, stats, created)
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

	postStmt, err := tx.PrepareContext(ctx, SQLForInsertPostWithCreated)
	if err != nil {
		return err
	}
	defer postStmt.Close()

	commentStmt, err := tx.PrepareContext(ctx, SQLForInsertCommentWithCreated)
	if err != nil {
		return err
	}
	defer commentStmt.Close()

	for i := range posts {
		postSeq := i + 1
		title := fmt.Sprintf("Post %d", postSeq)
		content := postContent
		stats := postStats
		created := postDate.NextFormatted()

		result, err := postStmt.ExecContext(ctx, title, content, stats, created)
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

			_, err = commentStmt.ExecContext(ctx, postID, name, content, stats, created)
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

		postStmt, err := tx.PrepareContext(ctx, SQLForInsertPostWithCreated)
		if err != nil {
			return err
		}
		defer postStmt.Close()

		postSeq := i + 1
		title := fmt.Sprintf("Post %d", postSeq)
		content := postContent
		stats := postStats
		created := postDate.NextFormatted()

		result, err := postStmt.ExecContext(ctx, title, content, stats, created)
		if err != nil {
			return err
		}

		postID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		commentStmt, err := tx.PrepareContext(ctx, SQLForInsertCommentWithCreated)
		if err != nil {
			return err
		}
		defer commentStmt.Close()

		for j := range comments {
			commentSeq := j + 1
			name := fmt.Sprintf("Comment %d.%d", postSeq, commentSeq)
			content := commentContent
			stats := commentStats
			created := commentDate.NextFormatted()

			_, err = commentStmt.ExecContext(ctx, postID, name, content, stats, created)
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

	row := db.readDB.QueryRowContext(ctx, SQLForSelectPostByID, id)
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

	row := tx.QueryRowContext(ctx, SQLForSelectPostByID, id)
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
	postStmt, err := db.readDB.PrepareContext(ctx, SQLForSelectPostByID)
	if err != nil {
		return nil, nil, err
	}
	defer postStmt.Close()

	p = &Post{ID: id}

	row := postStmt.QueryRowContext(ctx, id)
	if err = row.Scan(&p.Title, &p.Content, &p.Created, &p.Stats); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p = nil
		}
		return nil, nil, err
	}

	commentStmt, err := db.readDB.PrepareContext(ctx, SQLForSelectCommentsByPostID)
	if err != nil {
		return nil, nil, err
	}
	defer commentStmt.Close()

	rows, err := commentStmt.QueryContext(ctx, id)
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

	postStmt, err := tx.PrepareContext(ctx, SQLForSelectPostByID)
	if err != nil {
		return nil, nil, err
	}
	defer postStmt.Close()

	p = &Post{ID: id}

	row := postStmt.QueryRowContext(ctx, id)
	if err = row.Scan(&p.Title, &p.Content, &p.Created, &p.Stats); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p = nil
		}
		return nil, nil, err
	}

	commentStmt, err := tx.PrepareContext(ctx, SQLForSelectCommentsByPostID)
	if err != nil {
		return nil, nil, err
	}
	defer commentStmt.Close()

	rows, err := commentStmt.QueryContext(ctx, id)
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
	result, err := db.writeDB.ExecContext(ctx, SQLForInsertPost, title, content, stats)
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

	result, err := tx.ExecContext(ctx, SQLForInsertPost, title, content, stats)
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
	postStmt, err := db.writeDB.PrepareContext(ctx, SQLForInsertPost)
	if err != nil {
		return 0, err
	}
	defer postStmt.Close()

	result, err := postStmt.ExecContext(ctx, postTitle, postContent, postStats)
	if err != nil {
		return 0, err
	}

	postID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	commentStmt, err := db.writeDB.PrepareContext(ctx, SQLForInsertComment)
	if err != nil {
		return 0, err
	}
	defer commentStmt.Close()

	for range comments {
		_, err = commentStmt.ExecContext(ctx, postID, commentName, commentContent, commentStats)
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

	postStmt, err := tx.PrepareContext(ctx, SQLForInsertPost)
	if err != nil {
		return 0, err
	}
	defer postStmt.Close()

	result, err := postStmt.ExecContext(ctx, postTitle, postContent, postStats)
	if err != nil {
		return 0, err
	}

	postID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	commentStmt, err := tx.PrepareContext(ctx, SQLForInsertComment)
	if err != nil {
		return 0, err
	}
	defer commentStmt.Close()

	for range comments {
		_, err = commentStmt.ExecContext(ctx, postID, commentName, commentContent, commentStats)
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

	rows, err := db.readDB.QueryContext(ctx, SQLForQueryCorrelated)
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

	rows, err := db.readDB.QueryContext(ctx, SQLForQueryGroupBy)
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

	rows, err := db.readDB.QueryContext(ctx, SQLForQueryJSON)
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

	rows, err := db.readDB.QueryContext(ctx, SQLForQueryNonRecursiveCTE)
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

	rows, err := db.readDB.QueryContext(ctx, SQLForQueryOrderBy)
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

	rows, err := db.readDB.QueryContext(ctx, SQLForQueryRecursiveCTE)
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

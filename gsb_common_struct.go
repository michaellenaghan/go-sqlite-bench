package go_sqlite_bench

import (
	"time"
)

type Post struct {
	ID      int64
	Title   string
	Content string
	Created string
	Stats   string
}

type Comment struct {
	ID      int64
	Name    string
	Content string
	Created string
	Stats   string
}

type PostDate struct {
	next time.Time

	*CommentDate
}

func NewPostDate(posts int) *PostDate {
	next := time.Now().UTC().Truncate(time.Second).AddDate(0, 0, -posts)
	return &PostDate{next: next, CommentDate: &CommentDate{next: next}}
}

func (pd *PostDate) NextFormatted() string {
	prev := pd.next
	pd.next = prev.Add(time.Duration(24) * time.Hour)
	pd.CommentDate.next = pd.next

	return prev.Format(time.RFC3339)
}

type CommentDate struct {
	next time.Time
}

func (cd *CommentDate) NextFormatted() string {
	prev := cd.next
	cd.next = prev.Add(time.Duration(1) * time.Hour)

	return prev.Format(time.RFC3339)
}

type Explain struct {
	Addr    int64
	Opcode  string
	P1      string
	P2      string
	P3      string
	P4      string
	P5      string
	Comment string
}

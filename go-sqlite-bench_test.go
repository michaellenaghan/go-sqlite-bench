package go_sqlite_bench

import (
	"errors"
	"flag"
	"fmt"
	"math/rand/v2"
	"path"
	"regexp"
	"testing"
	"time"
)

var (
	defaultMaxReadConnections  = flag.Int("gsb-max-read-connections", 0, "max read connections (>= 0, 0 = use write connections)")
	defaultMaxWriteConnections = flag.Int("gsb-max-write-connections", 32, "max write connections (>= 1)")

	defaultPosts          = flag.Int("gsb-posts", 1000, "number of posts")
	defaultPostParagraphs = flag.Int("gsb-post-paragraphs", 50, "number of paragraphs per post")

	defaultComments          = flag.Int("gsb-comments", 25, "number of comments per post")
	defaultCommentParagraphs = flag.Int("gsb-comment-paragraphs", 5, "number of paragraphs per comment")
)

// Log all compile-time options.
func TestOptions(t *testing.T) {
	db := newDB(t, 0, 1)
	defer db.Close()

	options, err := db.Options(t.Context())
	noErr(t, err)

	for _, option := range options {
		t.Log("OPTION", option)
	}
}

// Test selected pragma values.
func TestPragma(t *testing.T) {
	db := newDB(t, 0, 1)
	defer db.Close()

	type testCase struct {
		name  string
		value string
	}

	tests := []testCase{
		{name: "busy_timeout", value: "10000"},
		{name: "foreign_keys", value: "1"},
		{name: "journal_mode", value: "wal"},
		{name: "synchronous", value: "1"},
	}

	for _, tc := range tests {
		want := tc.value
		got, err := db.Pragma(t.Context(), tc.name)
		noErr(t, err)
		if got != want {
			t.Errorf("pragma name %v: want value %v, got %v", tc.name, want, got)
		}
	}
}

// Log selected pragma values.
func TestPragmas(t *testing.T) {
	db := newDB(t, 0, 1)
	defer db.Close()

	pragmas, err := db.Pragmas(t.Context(), []string{
		"auto_vacuum",
		"automatic_index",
		"busy_timeout",
		"cache_size",
		"cache_spill",
		"cell_size_check",
		"checkpoint_fullfsync",
		"defer_foreign_keys",
		"encoding",
		"foreign_keys",
		"fullfsync",
		"hard_heap_limit",
		"journal_mode",
		"journal_size_limit",
		"locking_mode",
		"mmap_size",
		"page_size",
		"query_only",
		"read_uncommitted",
		"recursive_triggers",
		"reverse_unordered_selects",
		"secure_delete",
		"soft_heap_limit",
		"synchronous",
		"temp_store",
		"threads",
		"wal_autocheckpoint",
	})
	noErr(t, err)

	for _, pragma := range pragmas {
		t.Log("PRAGMA", pragma)
	}
}

// Test whether or not time values round trip without additional effort.
func TestTime(t *testing.T) {
	db := newDB(t, 0, 1)
	defer db.Close()

	in := time.Now().Truncate(time.Second)
	out, err := db.Time(t.Context(), in)
	if err != nil {
		t.Skipf("skip: can't roundtrip time without additional effort: %v", err)
	}
	if in.Compare(out) != 0 {
		t.Errorf("want %v, got %v", in, out)
	}

	t.Log("TIME", in, out)
}

// Log the SQLite version.
func TestVersion(t *testing.T) {
	db := newDB(t, 0, 1)
	defer db.Close()

	version, err := db.Version(t.Context())
	noErr(t, err)

	t.Log("VERSION", version)
}

// ===

// Test the Populate methods.
func TestPopulate(t *testing.T) {
	posts := 10
	postParagraphs := 10

	postID := int64(1)

	postTitle := regexp.MustCompile(`^Post \d+$`)
	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postCreated := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)
	postStats := LoremIpsumJSON

	comments := 10
	commentParagraphs := 10

	commentName := regexp.MustCompile(`^Comment \d+\.\d+$`)
	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentCreated := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)
	commentStats := LoremIpsumJSON

	// Test the PopulateDB method.
	t.Run("PopulateDB", func(t *testing.T) {
		db := newPreparedDB(t, 0, 1)
		defer db.Close()

		err := db.PopulateDB(t.Context(), posts, postParagraphs, comments, commentParagraphs)
		noErr(t, err)

		// Check counts

		n, err := db.CountPosts(t.Context())
		noErr(t, err)

		if int64(posts) != n {
			t.Errorf("want posts %v, got %v", posts, n)
		}

		n, err = db.CountComments(t.Context())
		noErr(t, err)

		if int64(posts*comments) != n {
			t.Errorf("want comments %v, got %v", posts*comments, n)
		}

		// Check first post and comments

		p, cs, err := db.ReadPostAndCommentsWithTx(t.Context(), postID)
		noErr(t, err)

		if postID != p.ID {
			t.Errorf("want postID %v, got %v", postID, p.ID)
		}
		if !postTitle.MatchString(p.Title) {
			t.Errorf("want postTitle %v, got %v", postTitle, p.Title)
		}
		if postContent != p.Content {
			t.Errorf("want postContent %v, got %v", postContent, p.Content)
		}
		if !postCreated.MatchString(p.Created) {
			t.Errorf("want postCreated %v, got %v", postCreated, p.Created)
		}
		if postStats != p.Stats {
			t.Errorf("want postStats %v, got %v", postStats, p.Stats)
		}

		if comments != len(cs) {
			t.Errorf("want comments %v, got %v", comments, len(cs))
		}

		for _, c := range cs {
			if !commentName.MatchString(c.Name) {
				t.Errorf("want commentName %v, got %v", commentName, c.Name)
			}
			if commentContent != c.Content {
				t.Errorf("want commentContent %v, got %v", commentContent, c.Content)
			}
			if !commentCreated.MatchString(c.Created) {
				t.Errorf("want commentCreated %v, got %v", commentCreated, c.Created)
			}
			if commentStats != c.Stats {
				t.Errorf("want commentStats %v, got %v", commentStats, c.Stats)
			}
		}
	})

	// Test the PopulateDBWithTx method.
	t.Run("PopulateDBWithTx", func(t *testing.T) {
		db := newPreparedDB(t, 0, 1)
		defer db.Close()

		err := db.PopulateDBWithTx(t.Context(), posts, postParagraphs, comments, commentParagraphs)
		noErr(t, err)

		// Check counts

		n, err := db.CountPosts(t.Context())
		noErr(t, err)

		if int64(posts) != n {
			t.Errorf("want posts %v, got %v", posts, n)
		}

		n, err = db.CountComments(t.Context())
		noErr(t, err)

		if int64(posts*comments) != n {
			t.Errorf("want comments %v, got %v", posts*comments, n)
		}

		// Check first post and comments

		p, cs, err := db.ReadPostAndCommentsWithTx(t.Context(), postID)
		noErr(t, err)

		if postID != p.ID {
			t.Errorf("want postID %v, got %v", postID, p.ID)
		}
		if !postTitle.MatchString(p.Title) {
			t.Errorf("want postTitle %v, got %v", postTitle, p.Title)
		}
		if postContent != p.Content {
			t.Errorf("want postContent %v, got %v", postContent, p.Content)
		}
		if !postCreated.MatchString(p.Created) {
			t.Errorf("want postCreated %v, got %v", postCreated, p.Created)
		}
		if postStats != p.Stats {
			t.Errorf("want postStats %v, got %v", postStats, p.Stats)
		}

		if comments != len(cs) {
			t.Errorf("want comments %v, got %v", comments, len(cs))
		}

		for _, c := range cs {
			if !commentName.MatchString(c.Name) {
				t.Errorf("want commentName %v, got %v", commentName, c.Name)
			}
			if commentContent != c.Content {
				t.Errorf("want commentContent %v, got %v", commentContent, c.Content)
			}
			if !commentCreated.MatchString(c.Created) {
				t.Errorf("want commentCreated %v, got %v", commentCreated, c.Created)
			}
			if commentStats != c.Stats {
				t.Errorf("want commentStats %v, got %v", commentStats, c.Stats)
			}
		}
	})

	// Test the PopulateDBWithTxs method.
	t.Run("PopulateDBWithTxs", func(t *testing.T) {
		db := newPreparedDB(t, 0, 1)
		defer db.Close()

		err := db.PopulateDBWithTxs(t.Context(), posts, postParagraphs, comments, commentParagraphs)
		noErr(t, err)

		// Check counts

		n, err := db.CountPosts(t.Context())
		noErr(t, err)

		if int64(posts) != n {
			t.Errorf("want posts %v, got %v", posts, n)
		}

		n, err = db.CountComments(t.Context())
		noErr(t, err)

		if int64(posts*comments) != n {
			t.Errorf("want comments %v, got %v", posts*comments, n)
		}

		// Check first post and comments

		p, cs, err := db.ReadPostAndCommentsWithTx(t.Context(), postID)
		noErr(t, err)

		if postID != p.ID {
			t.Errorf("want postID %v, got %v", postID, p.ID)
		}
		if !postTitle.MatchString(p.Title) {
			t.Errorf("want postTitle %v, got %v", postTitle, p.Title)
		}
		if postContent != p.Content {
			t.Errorf("want postContent %v, got %v", postContent, p.Content)
		}
		if !postCreated.MatchString(p.Created) {
			t.Errorf("want postCreated %v, got %v", postCreated, p.Created)
		}
		if postStats != p.Stats {
			t.Errorf("want postStats %v, got %v", postStats, p.Stats)
		}

		if comments != len(cs) {
			t.Errorf("want comments %v, got %v", comments, len(cs))
		}

		for _, c := range cs {
			if !commentName.MatchString(c.Name) {
				t.Errorf("want commentName %v, got %v", commentName, c.Name)
			}
			if commentContent != c.Content {
				t.Errorf("want commentContent %v, got %v", commentContent, c.Content)
			}
			if !commentCreated.MatchString(c.Created) {
				t.Errorf("want commentCreated %v, got %v", commentCreated, c.Created)
			}
			if commentStats != c.Stats {
				t.Errorf("want commentStats %v, got %v", commentStats, c.Stats)
			}
		}
	})
}

// Test the ReadWrite methods.
func TestReadWrite(t *testing.T) {
	db := newPreparedDB(t, 0, 1)
	defer db.Close()

	postParagraphs := 10

	postTitle := LoremIpsum
	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON

	comments := 10
	commentParagraphs := 10

	commentName := LoremIpsum
	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON

	// Test the ReadPost and WritePost methods.
	t.Run("ReadWritePost", func(t *testing.T) {
		postID, err := db.WritePost(t.Context(), postTitle, postContent, postStats)
		noErr(t, err)

		p, err := db.ReadPost(t.Context(), postID)
		noErr(t, err)

		if postID != p.ID {
			t.Errorf("want postID %v, got %v", postID, p.ID)
		}
		if postTitle != p.Title {
			t.Errorf("want postTitle %v, got %v", postTitle, p.Title)
		}
		if p.Content != p.Content {
			t.Errorf("want postContent %v, got %v", postContent, p.Content)
		}
	})

	// Test the ReadPostWithTx and WritePostWithTx methods.
	t.Run("ReadWritePostWithTx", func(t *testing.T) {
		postID, err := db.WritePostWithTx(t.Context(), postTitle, postContent, postStats)
		noErr(t, err)

		p, err := db.ReadPostWithTx(t.Context(), postID)
		noErr(t, err)

		if postID != p.ID {
			t.Errorf("want postID %v, got %v", postID, p.ID)
		}
		if postTitle != p.Title {
			t.Errorf("want postTitle %v, got %v", postTitle, p.Title)
		}
		if p.Content != p.Content {
			t.Errorf("want postContent %v, got %v", postContent, p.Content)
		}
	})

	// Test the ReadPostAndComments and WritePostAndComments methods.
	t.Run("ReadWritePostAndComments", func(t *testing.T) {
		postID, err := db.WritePostAndComments(t.Context(), postTitle, postContent, postStats, comments, commentName, commentContent, commentStats)
		noErr(t, err)

		p, cs, err := db.ReadPostAndComments(t.Context(), postID)
		noErr(t, err)

		if postID != p.ID {
			t.Errorf("want postID %v, got %v", postID, p.ID)
		}
		if postTitle != p.Title {
			t.Errorf("want postTitle %v, got %v", postTitle, p.Title)
		}
		if postContent != p.Content {
			t.Errorf("want postContent %v, got %v", postContent, p.Content)
		}

		if comments != len(cs) {
			t.Errorf("want comments %v, got %v", comments, len(cs))
		}

		for _, c := range cs {
			if commentName != c.Name {
				t.Errorf("want commentName %v, got %v", commentName, c.Name)
			}
			if commentContent != c.Content {
				t.Errorf("want commentContent %v, got %v", commentContent, c.Content)
			}
		}
	})

	// Test the ReadPostAndCommentsWithTx and WritePostAndCommentsWithTx methods.
	t.Run("ReadWritePostAndCommentsWithTx", func(t *testing.T) {
		postID, err := db.WritePostAndCommentsWithTx(t.Context(), postTitle, postContent, postStats, comments, commentName, commentContent, commentStats)
		noErr(t, err)

		p, cs, err := db.ReadPostAndCommentsWithTx(t.Context(), postID)
		noErr(t, err)

		if postID != p.ID {
			t.Errorf("want postID %v, got %v", postID, p.ID)
		}
		if postTitle != p.Title {
			t.Errorf("want postTitle %v, got %v", postTitle, p.Title)
		}
		if postContent != p.Content {
			t.Errorf("want postContent %v, got %v", postContent, p.Content)
		}

		if comments != len(cs) {
			t.Errorf("want comments %v, got %v", comments, len(cs))
		}

		for _, c := range cs {
			if commentName != c.Name {
				t.Errorf("want commentName %v, got %v", commentName, c.Name)
			}
			if commentContent != c.Content {
				t.Errorf("want commentContent %v, got %v", commentContent, c.Content)
			}
		}
	})
}

// Test the Query methods.
func TestQuery(t *testing.T) {
	posts := 10
	postParagraphs := 10

	comments := 10
	commentParagraphs := 10

	db := newPopulatedDB(t, 0, 1, posts, postParagraphs, comments, commentParagraphs)
	defer db.Close()

	// The following tests assume some things about the PostDate{} implementation —
	// e.g., that there's one post per day and that posts start "-posts" days ago.

	// Test the QueryCorrelated method.
	t.Run("Correlated", func(t *testing.T) {
		n, err := db.QueryCorrelated(t.Context())
		noErr(t, err)

		// 1 row per post
		if n != posts {
			t.Errorf("want n %d, got %d", posts, n)
		}
	})

	// Test the QueryGroupBy method.
	t.Run("GroupBy", func(t *testing.T) {
		n, err := db.QueryGroupBy(t.Context())
		noErr(t, err)

		months := func(posts int) int {
			if posts > 0 {
				pStart, err := db.ReadPost(t.Context(), int64(1))
				noErr(t, err)

				startDate, err := time.Parse(time.RFC3339, pStart.Created)
				noErr(t, err)

				pEnd, err := db.ReadPost(t.Context(), int64(posts))
				noErr(t, err)

				endDate, err := time.Parse(time.RFC3339, pEnd.Created)
				noErr(t, err)

				startYear, startMonth, _ := startDate.Date()
				endYear, endMonth, _ := endDate.Date()

				return int(endYear-startYear)*12 + int(endMonth-startMonth) + 1
			} else {
				return 0
			}
		}

		// 1 row per month
		if n != months(posts) {
			t.Errorf("want n %d, got %d", months(posts), n)
		}
	})

	// Test the QueryJSON method.
	t.Run("JSON", func(t *testing.T) {
		n, err := db.QueryJSON(t.Context())
		noErr(t, err)

		// 1 row per day, 1 day per post
		if n != posts {
			t.Errorf("want n %d, got %d", posts, n)
		}
	})

	// Test the QueryOrderBy method.
	t.Run("OrderBy", func(t *testing.T) {
		n, err := db.QueryOrderBy(t.Context())
		noErr(t, err)

		// 1 row per comment
		if n != posts*comments {
			t.Errorf("want n %d, got %d", posts*comments, n)
		}
	})

	// Test the QueryRecursiveCTE method.
	t.Run("RecursiveCTE", func(t *testing.T) {
		n, err := db.QueryRecursiveCTE(t.Context())
		noErr(t, err)

		// 1 row per day, always 31 days
		if n != 31 {
			t.Errorf("want n %d, got %d", 31, n)
		}
	})

	// Test the QueryWindow method.
	t.Run("Window", func(t *testing.T) {
		n, err := db.QueryWindow(t.Context())
		noErr(t, err)

		// 1 row per day, 1 day per post
		if n != posts {
			t.Errorf("want n %d, got %d", posts, n)
		}
	})
}

// ===

// Run the Baseline benchmarks.
func BenchmarkBaseline(b *testing.B) {
	db := newDB(b, *defaultMaxReadConnections, *defaultMaxWriteConnections)
	defer db.Close()

	// Run the Conn Baseline benchmark.
	b.Run("Conn", func(b *testing.B) {
		for b.Loop() {
			err := db.Conn(b.Context())
			noErr(b, err)
		}
	})

	// Run the Conn Baseline benchmark in parallel.
	b.Run("ConnParallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				err := db.Conn(b.Context())
				noErr(b, err)
			}
		})
	})

	// Run the Select1 Baseline benchmark.
	b.Run("Select1", func(b *testing.B) {
		for b.Loop() {
			err := db.Select1(b.Context())
			noErr(b, err)
		}
	})

	// Run the Select1 Baseline benchmark in parallel.
	b.Run("Select1Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				err := db.Select1(b.Context())
				noErr(b, err)
			}
		})
	})

	// Run the Select1PrePrepared Baseline benchmark.
	b.Run("Select1PrePrepared", func(b *testing.B) {
		for b.Loop() {
			err := db.Select1PrePrepared(b.Context())
			noErr(b, err)
		}
	})

	// Run the Select1PrePrepared Baseline benchmark in parallel.
	b.Run("Select1PrePreparedParallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				err := db.Select1PrePrepared(b.Context())
				noErr(b, err)
			}
		})
	})
}

// ===

// Run the Populate benchmarks.
func BenchmarkPopulate(b *testing.B) {
	db := newPreparedDB(b, *defaultMaxReadConnections, *defaultMaxWriteConnections)
	defer db.Close()

	posts := *defaultPosts
	postParagraphs := *defaultPostParagraphs

	comments := *defaultComments
	commentParagraphs := *defaultCommentParagraphs

	if !testing.Short() {
		// Run the PopulateDB Populate benchmark.
		b.Run("PopulateDB", func(b *testing.B) {
			for b.Loop() {
				err := db.PopulateDB(b.Context(), posts, postParagraphs, comments, commentParagraphs)
				noErr(b, err)
			}
		})
	}

	if !testing.Short() {
		// Run the PopulateDBWithTx Populate benchmark.
		b.Run("PopulateDBWithTx", func(b *testing.B) {
			for b.Loop() {
				err := db.PopulateDBWithTx(b.Context(), posts, postParagraphs, comments, commentParagraphs)
				noErr(b, err)
			}
		})
	}

	// Run the PopulateDBWithTxs Populate benchmark.
	b.Run("PopulateDBWithTxs", func(b *testing.B) {
		for b.Loop() {
			err := db.PopulateDBWithTxs(b.Context(), posts, postParagraphs, comments, commentParagraphs)
			noErr(b, err)
		}
	})
}

// Run the ReadWrite benchmarks.
func BenchmarkReadWrite(b *testing.B) {
	db := newPopulatedDB(b, *defaultMaxReadConnections, *defaultMaxWriteConnections, *defaultPosts, *defaultPostParagraphs, *defaultComments, *defaultCommentParagraphs)
	defer db.Close()

	posts := *defaultPosts
	postParagraphs := *defaultPostParagraphs

	postID := max(int64(posts/2), 1)

	postTitle := LoremIpsum
	postContent := Paragraphs(LoremIpsum, postParagraphs)
	postStats := LoremIpsumJSON

	comments := *defaultComments
	commentParagraphs := *defaultCommentParagraphs

	commentName := LoremIpsum
	commentContent := Paragraphs(LoremIpsum, commentParagraphs)
	commentStats := LoremIpsumJSON

	if !testing.Short() {
		// Run the ReadPost ReadWrite benchmark.
		b.Run("ReadPost", func(b *testing.B) {
			if posts == 0 {
				b.Skipf("skip: no posts")
			}

			for b.Loop() {
				_, err := db.ReadPost(b.Context(), postID)
				noErr(b, err)
			}
		})
	}

	// Run the ReadPostWithTx ReadWrite benchmark.
	b.Run("ReadPostWithTx", func(b *testing.B) {
		if posts == 0 {
			b.Skipf("skip: no posts")
		}

		for b.Loop() {
			_, err := db.ReadPostWithTx(b.Context(), postID)
			noErr(b, err)
		}
	})

	if !testing.Short() {
		// Run the ReadPostAndComments ReadWrite benchmark.
		b.Run("ReadPostAndComments", func(b *testing.B) {
			if posts == 0 {
				b.Skipf("skip: no posts")
			}

			for b.Loop() {
				_, _, err := db.ReadPostAndComments(b.Context(), postID)
				noErr(b, err)
			}
		})
	}

	// Run the ReadPostAndCommentsWithTx ReadWrite benchmark.
	b.Run("ReadPostAndCommentsWithTx", func(b *testing.B) {
		if posts == 0 {
			b.Skipf("skip: no posts")
		}

		for b.Loop() {
			_, _, err := db.ReadPostAndCommentsWithTx(b.Context(), postID)
			noErr(b, err)
		}
	})

	if !testing.Short() {
		// Run the WritePost ReadWrite benchmark.
		b.Run("WritePost", func(b *testing.B) {
			for b.Loop() {
				_, err := db.WritePost(b.Context(), postTitle, postContent, postStats)
				noErr(b, err)
			}
		})
	}

	// Run the WritePostWithTx ReadWrite benchmark.
	b.Run("WritePostWithTx", func(b *testing.B) {
		for b.Loop() {
			_, err := db.WritePostWithTx(b.Context(), postTitle, postContent, postStats)
			noErr(b, err)
		}
	})

	if !testing.Short() {
		// Run the WritePostAndComments ReadWrite benchmark.
		b.Run("WritePostAndComments", func(b *testing.B) {
			for b.Loop() {
				_, err := db.WritePostAndComments(b.Context(), postTitle, postContent, postStats, comments, commentName, commentContent, commentStats)
				noErr(b, err)
			}
		})
	}

	// Run the WritePostAndCommentsWithTx ReadWrite benchmark.
	b.Run("WritePostAndCommentsWithTx", func(b *testing.B) {
		for b.Loop() {
			_, err := db.WritePostAndCommentsWithTx(b.Context(), postTitle, postContent, postStats, comments, commentName, commentContent, commentStats)
			noErr(b, err)
		}
	})

	if !testing.Short() {
		// Run the ReadOrWritePostAndComments ReadWrite benchmark.
		b.Run("ReadOrWritePostAndComments", func(b *testing.B) {
			if posts == 0 {
				b.Skipf("skip: no posts")
			}

			for _, writeRate := range []int{10, 90} {
				b.Run(fmt.Sprintf("write_rate=%d", writeRate), func(b *testing.B) {
					for b.Loop() {
						if rand.IntN(100) < writeRate {
							_, err := db.WritePostAndComments(b.Context(), postTitle, postContent, postStats, comments, commentName, commentContent, commentStats)
							noErr(b, err)
						} else {
							_, _, err := db.ReadPostAndComments(b.Context(), postID)
							noErr(b, err)
						}
					}
				})
			}
		})
	}

	if !testing.Short() {
		// Run the ReadOrWritePostAndCommentsParallel ReadWrite benchmark.
		b.Run("ReadOrWritePostAndCommentsParallel", func(b *testing.B) {
			if posts == 0 {
				b.Skipf("skip: no posts")
			}

			for _, writeRate := range []int{10, 90} {
				b.Run(fmt.Sprintf("write_rate=%d", writeRate), func(b *testing.B) {
					b.RunParallel(func(pb *testing.PB) {
						for pb.Next() {
							if rand.IntN(100) < writeRate {
								_, err := db.WritePostAndComments(b.Context(), postTitle, postContent, postStats, comments, commentName, commentContent, commentStats)
								noErr(b, err)
							} else {
								_, _, err := db.ReadPostAndComments(b.Context(), postID)
								noErr(b, err)
							}
						}
					})
				})
			}
		})
	}

	// Run the ReadOrWritePostAndCommentsWithTx ReadWrite benchmark.
	b.Run("ReadOrWritePostAndCommentsWithTx", func(b *testing.B) {
		if posts == 0 {
			b.Skipf("skip: no posts")
		}

		for _, writeRate := range []int{10, 90} {
			b.Run(fmt.Sprintf("write_rate=%d", writeRate), func(b *testing.B) {
				for b.Loop() {
					if rand.IntN(100) < writeRate {
						_, err := db.WritePostAndCommentsWithTx(b.Context(), postTitle, postContent, postStats, comments, commentName, commentContent, commentStats)
						noErr(b, err)
					} else {
						_, _, err := db.ReadPostAndCommentsWithTx(b.Context(), postID)
						noErr(b, err)
					}
				}
			})
		}
	})

	// Run the ReadOrWritePostAndCommentsWithTxParallel ReadWrite benchmark.
	b.Run("ReadOrWritePostAndCommentsWithTxParallel", func(b *testing.B) {
		if posts == 0 {
			b.Skipf("skip: no posts")
		}

		for _, writeRate := range []int{10, 90} {
			b.Run(fmt.Sprintf("write_rate=%d", writeRate), func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						if rand.IntN(100) < writeRate {
							_, err := db.WritePostAndCommentsWithTx(b.Context(), postTitle, postContent, postStats, comments, commentName, commentContent, commentStats)
							noErr(b, err)
						} else {
							_, _, err := db.ReadPostAndCommentsWithTx(b.Context(), postID)
							noErr(b, err)
						}
					}
				})
			})
		}
	})
}

// ===

// Run the Query benchmarks.
func BenchmarkQuery(b *testing.B) {
	db := newPopulatedDB(b, *defaultMaxReadConnections, *defaultMaxWriteConnections, *defaultPosts, *defaultPostParagraphs, *defaultComments, *defaultCommentParagraphs)
	defer db.Close()

	b.ResetTimer()

	// Run the QueryCorrelated Query benchmark.
	b.Run("Correlated", func(b *testing.B) {
		for b.Loop() {
			_, err := db.QueryCorrelated(b.Context())
			noErr(b, err)
		}
	})

	// Run the QueryCorrelated Query benchmark in parallel.
	b.Run("CorrelatedParallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := db.QueryCorrelated(b.Context())
				noErr(b, err)
			}
		})
	})

	// Run the QueryGroupBy Query benchmark.
	b.Run("GroupBy", func(b *testing.B) {
		for b.Loop() {
			_, err := db.QueryGroupBy(b.Context())
			noErr(b, err)
		}
	})

	// Run the QueryGroupBy Query benchmark in parallel.
	b.Run("GroupByParallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := db.QueryGroupBy(b.Context())
				noErr(b, err)
			}
		})
	})

	// Run the QueryJSON Query benchmark.
	b.Run("JSON", func(b *testing.B) {
		for b.Loop() {
			_, err := db.QueryJSON(b.Context())
			noErr(b, err)
		}
	})

	// Run the QueryJSON Query benchmark in parallel.
	b.Run("JSONParallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := db.QueryJSON(b.Context())
				noErr(b, err)
			}
		})
	})

	// Run the QueryOrderBy Query benchmark.
	b.Run("OrderBy", func(b *testing.B) {
		for b.Loop() {
			_, err := db.QueryOrderBy(b.Context())
			noErr(b, err)
		}
	})

	// Run the QueryOrderBy Query benchmark in parallel.
	b.Run("OrderByParallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := db.QueryOrderBy(b.Context())
				noErr(b, err)
			}
		})
	})

	// Run the QueryRecursiveCTE Query benchmark.
	b.Run("RecursiveCTE", func(b *testing.B) {
		for b.Loop() {
			_, err := db.QueryRecursiveCTE(b.Context())
			noErr(b, err)
		}
	})

	// Run the QueryRecursiveCTE Query benchmark in parallel.
	b.Run("RecursiveCTEParallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := db.QueryRecursiveCTE(b.Context())
				noErr(b, err)
			}
		})
	})

	// Run the QueryWindow Query benchmark.
	b.Run("Window", func(b *testing.B) {
		for b.Loop() {
			_, err := db.QueryWindow(b.Context())
			noErr(b, err)
		}
	})

	// Run the QueryWindow Query benchmark in parallel.
	b.Run("WindowParallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := db.QueryWindow(b.Context())
				noErr(b, err)
			}
		})
	})
}

// Returns a new DB.
func newDB(tb testing.TB, maxReadConnections, maxWriteConnections int) *DB {
	if !(maxReadConnections >= 0) {
		noErr(tb, errors.New("maxReadConnections must be >= 0"))
	}
	if !(maxWriteConnections >= 1) {
		noErr(tb, errors.New("maxWriteConnections must be >= 1"))
	}

	db, err := NewDB(tb.Context(), path.Join(tb.TempDir(), "go-sqlite-bench.db"), maxReadConnections, maxWriteConnections)
	noErr(tb, err)

	return db
}

// Returns a new prepared DB.
func newPreparedDB(tb testing.TB, maxReadConnections, maxWriteConnections int) *DB {
	db := newDB(tb, maxReadConnections, maxWriteConnections)

	err := db.PrepareDBWithTx(tb.Context())
	noErr(tb, err)

	err = db.Analyze(tb.Context())
	noErr(tb, err)

	return db
}

// Returns a new populated DB.
func newPopulatedDB(tb testing.TB, maxReadConnections, maxWriteConnections int, posts, postParagraphs, comments, commentParagraphs int) *DB {
	db := newPreparedDB(tb, maxReadConnections, maxWriteConnections)

	err := db.PopulateDBWithTxs(tb.Context(), posts, postParagraphs, comments, commentParagraphs)
	noErr(tb, err)

	err = db.Analyze(tb.Context())
	noErr(tb, err)

	return db
}

// Signals a fatal error if err is not nil.
func noErr(tb testing.TB, err error) {
	tb.Helper()

	if err != nil {
		tb.Fatal("Error is not nil:", err)
	}
}

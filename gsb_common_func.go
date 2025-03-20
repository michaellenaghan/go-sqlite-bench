package go_sqlite_bench

import (
	"strings"
)

func Paragraphs(s string, n int) string {
	paragraphs := make([]string, n)

	for i := range n {
		paragraphs[i] = s
	}

	return strings.Join(paragraphs, "\n\n")
}

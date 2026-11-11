// Package text holds small shared string helpers used by flatten and minify.
package text

import "strings"

// NormalizeNewlines converts CRLF/CR to LF and trims one trailing newline.
func NormalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.TrimSuffix(s, "\n")
}

// LineSpan returns the 1-based original line range for normalized content.
// Empty content yields (0, 0).
func LineSpan(s string) (start, end int) {
	if s == "" {
		return 0, 0
	}
	return 1, strings.Count(s, "\n") + 1
}

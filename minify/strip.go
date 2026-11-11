package minify

import (
	"strconv"
	"strings"
	"unicode"
)

// StripResult is minified text plus a map from output line to original line.
type StripResult struct {
	// Lines are the kept source lines after comment/blank stripping (trailing
	// whitespace trimmed). Index i corresponds to dump body line i+1.
	Lines []string

	// Map holds the original 1-based line number for each entry in Lines.
	// Map[i] is the original line for Lines[i] (dump line i+1).
	Map []int
}

// Strip removes comments and blank lines while preserving logic and identifiers.
//
// One kept source line becomes one output line. Map entries are 1-based original
// line numbers. Comment styles:
//   - // and /* */ for all extensions
//   - # for py, yaml, and yml (hashComments)
//
// String literals (including Go raw `...` and Python """ / ”') are preserved
// so comment markers inside them are not stripped.
func Strip(src string, ext string) StripResult {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	src = strings.ReplaceAll(src, "\r", "\n")

	hashComments := ext == "py" || ext == "yaml" || ext == "yml"
	cleaned := stripComments(src, hashComments)

	var out StripResult
	rawLines := strings.Split(cleaned, "\n")
	for i, line := range rawLines {
		orig := i + 1
		trimmedRight := strings.TrimRightFunc(line, unicode.IsSpace)
		if strings.TrimSpace(trimmedRight) == "" {
			continue
		}
		out.Lines = append(out.Lines, trimmedRight)
		out.Map = append(out.Map, orig)
	}
	return out
}

func stripComments(src string, hashComments bool) string {
	var b strings.Builder
	b.Grow(len(src))

	inLineComment := false
	inBlockComment := false
	inSingle := false
	inDouble := false
	inRaw := false // Go raw string `
	inTripleDouble := false
	inTripleSingle := false

	runes := []rune(src)
	n := len(runes)

	for i := 0; i < n; i++ {
		r := runes[i]
		next := func() rune {
			if i+1 < n {
				return runes[i+1]
			}
			return 0
		}
		peek2 := func() (rune, rune) {
			var r1, r2 rune
			if i+1 < n {
				r1 = runes[i+1]
			}
			if i+2 < n {
				r2 = runes[i+2]
			}
			return r1, r2
		}

		if inLineComment {
			if r == '\n' {
				inLineComment = false
				b.WriteRune('\n')
			}
			continue
		}

		if inBlockComment {
			if r == '*' && next() == '/' {
				inBlockComment = false
				i++
			}
			continue
		}

		if inTripleDouble {
			r1, r2 := peek2()
			if r == '"' && r1 == '"' && r2 == '"' {
				b.WriteString(`"""`)
				i += 2
				inTripleDouble = false
				continue
			}
			b.WriteRune(r)
			continue
		}

		if inTripleSingle {
			r1, r2 := peek2()
			if r == '\'' && r1 == '\'' && r2 == '\'' {
				b.WriteString(`'''`)
				i += 2
				inTripleSingle = false
				continue
			}
			b.WriteRune(r)
			continue
		}

		if inRaw {
			b.WriteRune(r)
			if r == '`' {
				inRaw = false
			}
			continue
		}

		if inSingle {
			b.WriteRune(r)
			if r == '\\' && i+1 < n {
				b.WriteRune(runes[i+1])
				i++
				continue
			}
			if r == '\'' {
				inSingle = false
			}
			continue
		}

		if inDouble {
			b.WriteRune(r)
			if r == '\\' && i+1 < n {
				b.WriteRune(runes[i+1])
				i++
				continue
			}
			if r == '"' {
				inDouble = false
			}
			continue
		}

		// not in string/comment
		if hashComments {
			r1, r2 := peek2()
			if r == '"' && r1 == '"' && r2 == '"' {
				inTripleDouble = true
				b.WriteString(`"""`)
				i += 2
				continue
			}
			if r == '\'' && r1 == '\'' && r2 == '\'' {
				inTripleSingle = true
				b.WriteString(`'''`)
				i += 2
				continue
			}
		}

		if r == '`' {
			inRaw = true
			b.WriteRune(r)
			continue
		}
		if r == '\'' {
			inSingle = true
			b.WriteRune(r)
			continue
		}
		if r == '"' {
			inDouble = true
			b.WriteRune(r)
			continue
		}

		if r == '/' && next() == '/' {
			inLineComment = true
			i++
			continue
		}
		if r == '/' && next() == '*' {
			inBlockComment = true
			i++
			continue
		}
		if hashComments && r == '#' {
			inLineComment = true
			continue
		}

		b.WriteRune(r)
	}

	return b.String()
}

// FormatMap builds an "@MAP" payload "1>3,2>4,..." from 1-based original line
// numbers. Index i in m becomes dump line i+1 mapped to original line m[i].
// An empty map yields an empty string.
func FormatMap(m []int) string {
	if len(m) == 0 {
		return ""
	}
	parts := make([]string, len(m))
	for i, orig := range m {
		parts[i] = strconv.Itoa(i+1) + ">" + strconv.Itoa(orig)
	}
	return strings.Join(parts, ",")
}

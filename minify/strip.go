package minify

import (
	"strconv"
	"strings"
	"unicode"
)

// StripResult is minified text plus a map from output line -> original line.
type StripResult struct {
	Lines []string
	Map   []int // Map[outLineIndex] = original 1-based line
}

// Strip removes comments and blank lines while preserving logic and identifiers.
// One kept source line becomes one output line. map entries are 1-based.
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
			var a, b rune
			if i+1 < n {
				a = runes[i+1]
			}
			if i+2 < n {
				b = runes[i+2]
			}
			return a, b
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
			a, c := peek2()
			if r == '"' && a == '"' && c == '"' {
				b.WriteString(`"""`)
				i += 2
				inTripleDouble = false
				continue
			}
			b.WriteRune(r)
			continue
		}

		if inTripleSingle {
			a, c := peek2()
			if r == '\'' && a == '\'' && c == '\'' {
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
			a, c := peek2()
			if r == '"' && a == '"' && c == '"' {
				inTripleDouble = true
				b.WriteString(`"""`)
				i += 2
				continue
			}
			if r == '\'' && a == '\'' && c == '\'' {
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

// FormatMap builds "1>3,2>4,..." from 1-based map entries.
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

package collect

import (
	"path/filepath"
	"strings"
)

// matcher applies simple gitignore-like path rules from config ignore entries.
//
// Supported shapes:
//   - "name/"     — directory (and its contents) by basename or path prefix
//   - "path/seg"  — exact relative path or prefix path/seg/
//   - "*.min.js"  — basename or full-path glob via filepath.Match
//
// Blank lines and rules starting with "#" are skipped. Matching uses
// slash-separated paths relative to the source root.
type matcher struct {
	patterns []pattern
}

type pattern struct {
	raw    string
	dir    bool
	glob   bool
	parts  []string
	suffix string
}

// newMatcher builds a matcher from config ignore rules. Invalid or empty rules
// are skipped; globs are detected by the presence of *, ?, or [.
func newMatcher(rules []string) *matcher {
	m := &matcher{}
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" || strings.HasPrefix(rule, "#") {
			continue
		}
		p := pattern{raw: rule}
		if strings.HasSuffix(rule, "/") {
			p.dir = true
			rule = strings.TrimSuffix(rule, "/")
		}
		rule = strings.ReplaceAll(rule, "\\", "/")
		if strings.ContainsAny(rule, "*?[") {
			p.glob = true
			p.suffix = rule
		} else {
			p.parts = strings.Split(rule, "/")
		}
		m.patterns = append(m.patterns, p)
	}
	return m
}

// match reports whether relPath should be ignored.
// relPath must be slash-separated and relative to the source root.
// isDir is true when the path is a directory (enables SkipDir for dir rules).
func (m *matcher) match(relPath string, isDir bool) bool {
	relPath = strings.ReplaceAll(relPath, "\\", "/")
	relPath = strings.TrimPrefix(relPath, "./")
	if relPath == "." || relPath == "" {
		return false
	}

	for _, p := range m.patterns {
		if p.match(relPath, isDir) {
			return true
		}
	}
	return false
}

func (p pattern) match(relPath string, isDir bool) bool {
	if p.dir && !isDir {
		if !p.glob {
			return p.matchDirFile(relPath)
		}
		// Dir-marked globs fall through to matchGlob below.
	}
	if p.glob {
		return p.matchGlob(relPath)
	}
	if p.matchExactOrPrefix(relPath) {
		return true
	}
	return p.matchSegment(relPath)
}

// matchDirFile applies a non-glob directory rule (trailing "/") to a file path:
// basename equals the rule's last segment, exact path match, or path prefix.
func (p pattern) matchDirFile(relPath string) bool {
	base := filepath.Base(relPath)
	base = strings.ReplaceAll(base, "\\", "/")
	name := p.raw
	name = strings.TrimSuffix(name, "/")
	name = strings.ReplaceAll(name, "\\", "/")
	return base == filepath.Base(name) || relPath == name || strings.HasPrefix(relPath, name+"/")
}

func (p pattern) matchGlob(relPath string) bool {
	name := filepath.Base(relPath)
	ok, err := filepath.Match(p.suffix, name)
	if err == nil && ok {
		return true
	}
	ok, err = filepath.Match(p.suffix, relPath)
	return err == nil && ok
}

func (p pattern) matchExactOrPrefix(relPath string) bool {
	target := strings.Join(p.parts, "/")
	return relPath == target || strings.HasPrefix(relPath, target+"/")
}

// matchSegment matches a single-segment rule against any path component
// (e.g. rule "vendor" matches "pkg/vendor/x.go").
func (p pattern) matchSegment(relPath string) bool {
	if len(p.parts) != 1 {
		return false
	}
	for _, seg := range strings.Split(relPath, "/") {
		if seg == p.parts[0] {
			return true
		}
	}
	return false
}

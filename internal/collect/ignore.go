package collect

import (
	"path/filepath"
	"strings"
)

// Matcher applies gitignore-like path rules.
type Matcher struct {
	patterns []pattern
}

type pattern struct {
	raw    string
	dir    bool
	glob   bool
	parts  []string
	suffix string
}

// NewMatcher builds a matcher from config ignore rules.
func NewMatcher(rules []string) *Matcher {
	m := &Matcher{}
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

// Match reports whether relPath (slash-separated, relative to source root) should be ignored.
// isDir is true when the path is a directory.
func (m *Matcher) Match(relPath string, isDir bool) bool {
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
		// Still match if any path segment equals the dir name (ancestor skipped via WalkDir skip).
		base := filepath.Base(relPath)
		base = strings.ReplaceAll(base, "\\", "/")
		name := p.raw
		name = strings.TrimSuffix(name, "/")
		name = strings.ReplaceAll(name, "\\", "/")
		if !p.glob {
			if base == filepath.Base(name) || relPath == name || strings.HasPrefix(relPath, name+"/") {
				return true
			}
			return false
		}
	}

	if p.glob {
		name := filepath.Base(relPath)
		ok, err := filepath.Match(p.suffix, name)
		if err == nil && ok {
			return true
		}
		ok, err = filepath.Match(p.suffix, relPath)
		return err == nil && ok
	}

	target := strings.Join(p.parts, "/")
	if relPath == target {
		return true
	}
	if strings.HasPrefix(relPath, target+"/") {
		return true
	}
	// Match a path segment equal to the last part (e.g. "vendor" anywhere).
	if len(p.parts) == 1 {
		for _, seg := range strings.Split(relPath, "/") {
			if seg == p.parts[0] {
				return true
			}
		}
	}
	return false
}

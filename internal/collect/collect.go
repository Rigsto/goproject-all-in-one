// Package collect walks a source tree and returns text files that match the
// configured extensions and are not excluded by ignore rules.
//
// It is shared by flatten and minify. Results are sorted by slash-separated
// relative path. Files containing a NUL byte are treated as binary and skipped.
package collect

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golanganalyzer/config"
)

// File is a collected source file ready for writing into a dump.
type File struct {
	// RelPath is the path relative to the source root, using forward slashes.
	RelPath string

	// Content is the raw file bytes as read from disk.
	Content []byte
}

// Collect walks sourcePath and returns matching text files in sorted RelPath order.
// sourcePath must be an existing directory. cfg.Extensions and cfg.Ignore control
// inclusion; directory ignore matches cause filepath.SkipDir.
func Collect(sourcePath string, cfg *config.Config) ([]File, error) {
	absSource, err := filepath.Abs(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("resolve source: %w", err)
	}

	info, err := os.Stat(absSource)
	if err != nil {
		return nil, fmt.Errorf("stat source: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("source is not a directory: %s", absSource)
	}

	extSet := make(map[string]struct{}, len(cfg.Extensions))
	for _, ext := range cfg.Extensions {
		extSet[ext] = struct{}{}
	}
	matcher := newMatcher(cfg.Ignore)

	var files []File
	err = filepath.WalkDir(absSource, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(absSource, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		if rel == "." {
			return nil
		}

		if d.IsDir() {
			if matcher.match(rel, true) {
				return filepath.SkipDir
			}
			return nil
		}

		if matcher.match(rel, false) {
			return nil
		}

		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
		if _, ok := extSet[ext]; !ok {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", rel, err)
		}
		if isBinary(data) {
			return nil
		}

		files = append(files, File{RelPath: rel, Content: data})
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].RelPath < files[j].RelPath
	})
	return files, nil
}

// isBinary reports whether data looks binary by presence of a NUL byte.
func isBinary(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return true
		}
	}
	return false
}

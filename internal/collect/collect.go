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

// File is a collected source file ready for writing.
type File struct {
	RelPath string
	Content []byte
}

// Collect walks sourcePath and returns matching text files in sorted order.
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
	matcher := NewMatcher(cfg.Ignore)

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
			if matcher.Match(rel, true) {
				return filepath.SkipDir
			}
			return nil
		}

		if matcher.Match(rel, false) {
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

func isBinary(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return true
		}
	}
	return false
}

// Package flatten writes a full-source project dump for AI review.
//
// Output format per file:
//
//	===== FILE <relpath> | LINES 1-<n>
//	<source body>
//
// Body line N is original source line N. Newlines are normalized to \n; a
// trailing newline is ensured after each file body. Paths are relative to the
// --source root and sorted by collect.Collect.
package flatten

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golanganalyzer/config"
	"golanganalyzer/internal/collect"
	"golanganalyzer/internal/text"
)

// FlattenWithConfig writes a full-source dump of sourcePath into targetFile
// using an already-loaded config (used by the CLI after parsing --config).
// Parent directories of targetFile are created if needed.
func FlattenWithConfig(sourcePath, targetFile string, cfg *config.Config) error {
	files, err := collect.Collect(sourcePath, cfg)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(targetFile), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	var b strings.Builder
	for i, f := range files {
		content := text.NormalizeNewlines(string(f.Content))
		startLine, endLine := text.LineSpan(content)

		if i > 0 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "===== FILE %s | LINES %d-%d\n", f.RelPath, startLine, endLine)
		if content != "" {
			b.WriteString(content)
			b.WriteByte('\n')
		}
	}

	if err := os.WriteFile(targetFile, []byte(b.String()), 0o644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}

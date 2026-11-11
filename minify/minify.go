// Package minify writes a compressed project dump with line maps for AI review.
//
// Output format per file:
//
//	@FILE <relpath>
//	@ORIG 1-<n>
//	@MAP 1>3,2>4,...
//	<stripped lines>
//
// Comments and blank lines are removed by Strip while keeping logic and
// identifiers. @MAP out>orig means dump body line "out" came from original
// source line "orig". Numbering resets for every @FILE block.
package minify

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Rigsto/goproject-all-in-one/config"
	"github.com/Rigsto/goproject-all-in-one/internal/collect"
	"github.com/Rigsto/goproject-all-in-one/internal/text"
)

// MinifyWithConfig writes a minified dump of sourcePath into targetFile using
// an already-loaded config (used by the CLI after parsing --config).
// Parent directories of targetFile are created if needed.
func MinifyWithConfig(sourcePath, targetFile string, cfg *config.Config) error {
	files, err := collect.Collect(sourcePath, cfg)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(targetFile), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	var b strings.Builder
	for i, f := range files {
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(f.RelPath)), ".")
		raw := text.NormalizeNewlines(string(f.Content))
		startLine, endLine := text.LineSpan(raw)

		stripped := Strip(raw, ext)

		if i > 0 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "@FILE %s\n", f.RelPath)
		fmt.Fprintf(&b, "@ORIG %d-%d\n", startLine, endLine)
		fmt.Fprintf(&b, "@MAP %s\n", FormatMap(stripped.Map))
		for _, line := range stripped.Lines {
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}

	if err := os.WriteFile(targetFile, []byte(b.String()), 0o644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}

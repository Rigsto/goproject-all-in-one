package flatten

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golanganalyzer/config"
	"golanganalyzer/internal/collect"
)

// Flatten writes a full-source dump of sourcePath into targetFile.
func Flatten(sourcePath, targetFile string) error {
	cfg, err := config.Load(defaultConfigPath())
	if err != nil {
		return err
	}
	return FlattenWithConfig(sourcePath, targetFile, cfg)
}

// FlattenWithConfig is Flatten using an already-loaded config.
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
		content := string(f.Content)
		content = strings.ReplaceAll(content, "\r\n", "\n")
		content = strings.TrimSuffix(content, "\n")
		endLine := 0
		if content != "" {
			endLine = strings.Count(content, "\n") + 1
		}
		startLine := 0
		if endLine > 0 {
			startLine = 1
		}

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

func defaultConfigPath() string {
	return "config.yaml"
}

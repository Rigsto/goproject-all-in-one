package minify

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golanganalyzer/config"
	"golanganalyzer/internal/collect"
)

// Minify writes a compressed dump of sourcePath into targetFile with line maps.
func Minify(sourcePath, targetFile string) error {
	cfg, err := config.Load(defaultConfigPath())
	if err != nil {
		return err
	}
	return MinifyWithConfig(sourcePath, targetFile, cfg)
}

// MinifyWithConfig is Minify using an already-loaded config.
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
		raw := string(f.Content)
		raw = strings.ReplaceAll(raw, "\r\n", "\n")
		raw = strings.TrimSuffix(raw, "\n")
		endLine := 0
		if raw != "" {
			endLine = strings.Count(raw, "\n") + 1
		}
		startLine := 0
		if endLine > 0 {
			startLine = 1
		}

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

func defaultConfigPath() string {
	return "config.yaml"
}

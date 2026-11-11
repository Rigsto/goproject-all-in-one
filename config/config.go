// Package config loads and validates the YAML settings that control which
// files flatten and minify include.
//
// The default file is config.yaml at the process working directory. Extensions
// are normalized (leading dots stripped, lowercased). Ignore rules are passed
// through to internal/collect for gitignore-like matching.
package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds flatten/minify settings loaded from config.yaml.
type Config struct {
	// Extensions is the allow-list of file suffixes to collect (e.g. "go", "py").
	// Dots are optional in YAML; Load strips them and lowercases values.
	Extensions []string `yaml:"extensions"`

	// Ignore is a list of gitignore-like path rules (directory names ending in "/",
	// basename globs like "*.min.js", or relative path prefixes).
	Ignore []string `yaml:"ignore"`
}

// Load reads path as YAML, validates that Extensions is non-empty, normalizes
// extension strings, and returns the config. The file must exist and parse as YAML.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	cfg.normalize()
	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Extensions) == 0 {
		return fmt.Errorf("config: extensions must not be empty")
	}
	return nil
}

func (c *Config) normalize() {
	for i, ext := range c.Extensions {
		ext = strings.TrimSpace(ext)
		ext = strings.TrimPrefix(ext, ".")
		c.Extensions[i] = strings.ToLower(ext)
	}
}

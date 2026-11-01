package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds flatten/minify settings loaded from config.yaml.
type Config struct {
	Extensions []string `yaml:"extensions"`
	Ignore     []string `yaml:"ignore"`
}

// Load reads and validates a YAML config file.
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

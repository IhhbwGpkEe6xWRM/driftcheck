package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level driftcheck configuration.
type Config struct {
	Statefile string        `yaml:"statefile"`
	Region    string        `yaml:"region"`
	Profile   string        `yaml:"profile"`
	Ignore    []IgnoreRule  `yaml:"ignore"`
	Output    OutputConfig  `yaml:"output"`
}

// IgnoreRule describes a resource/attribute pattern to skip during drift detection.
type IgnoreRule struct {
	Resource  string `yaml:"resource"`
	Attribute string `yaml:"attribute"`
}

// OutputConfig controls how drift reports are written.
type OutputConfig struct {
	Format string `yaml:"format"` // "text" or "json"
	File   string `yaml:"file"`   // empty means stdout
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}
	return parse(data)
}

// parse unmarshals raw YAML bytes into a Config, applying defaults.
func parse(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config YAML: %w", err)
	}
	if cfg.Output.Format == "" {
		cfg.Output.Format = "text"
	}
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	return &cfg, nil
}

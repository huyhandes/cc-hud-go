package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds all configuration options
type Config struct {
	Theme      string
	Colors     map[string]string
	LineLayout string
	Display    DisplayConfig
	Git        GitConfig
	Tools      ToolsConfig
}

type DisplayConfig struct {
	Model      bool
	Context    bool
	Git        bool
	Tools      bool
	Agents     bool
	Tasks      bool
	RateLimits bool
	Duration   bool
	FetchOAuth bool
	Caveman    bool
	Ponytail   bool
}

type GitConfig struct {
	ShowBranch      bool
	ShowDirty       bool
	ShowAheadBehind bool
	ShowFileStats   bool
}

type ToolsConfig struct {
	ShowSkills bool
	ShowMCP    bool
}

// Default returns the hardcoded default config
func Default() *Config {
	return &Config{
		Theme:      "macchiato",
		Colors:     make(map[string]string),
		LineLayout: "expanded",
		Display: DisplayConfig{
			Model:      true,
			Context:    true,
			Git:        true,
			Tools:      true,
			Agents:     true,
			Tasks:      true,
			RateLimits: true,
			Duration:   true,
			FetchOAuth: true,
			Caveman:    true,
			Ponytail:   true,
		},
		Git: GitConfig{
			ShowBranch:      true,
			ShowDirty:       true,
			ShowAheadBehind: true,
			ShowFileStats:   true,
		},
		Tools: ToolsConfig{
			ShowSkills: true,
			ShowMCP:    true,
		},
	}
}

// LoadFromFile loads config from JSON file, returns defaults on any error
func LoadFromFile(path string) (*Config, error) {
	// Start with defaults
	cfg := Default()

	// Try to read file
	data, err := os.ReadFile(path)
	if err != nil {
		// Missing file is OK, just use defaults
		if os.IsNotExist(err) {
			return cfg, nil
		}
		// Other read errors: log but continue with defaults
		fmt.Fprintf(os.Stderr, "warning: failed to read config: %v\n", err)
		return cfg, nil
	}

	// Try to parse JSON
	if err := json.Unmarshal(data, cfg); err != nil {
		// Invalid JSON: log but continue with defaults
		fmt.Fprintf(os.Stderr, "warning: failed to parse config: %v\n", err)
		return Default(), nil
	}

	return cfg, nil
}

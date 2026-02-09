package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds all configuration options
type Config struct {
	Preset             string
	LineLayout         string
	PathLevels         int
	ContextValue       string
	SevenDayThreshold  int
	Display            DisplayConfig
	Git                GitConfig
	Tools              ToolsConfig
}

type DisplayConfig struct {
	Model      bool
	Path       bool
	Context    bool
	Git        bool
	Tools      bool
	Agents     bool
	Tasks      bool
	RateLimits bool
	Duration   bool
	Speed      bool
}

type GitConfig struct {
	ShowBranch      bool
	ShowDirty       bool
	ShowAheadBehind bool
	ShowFileStats   bool
}

type ToolsConfig struct {
	GroupByCategory bool
	ShowTopN        int
	ShowSkills      bool
	ShowMCP         bool
}

// Default returns a config with sensible defaults (full preset)
func Default() *Config {
	return &Config{
		Preset:            "full",
		LineLayout:        "expanded",
		PathLevels:        2,
		ContextValue:      "percentage",
		SevenDayThreshold: 80,
		Display: DisplayConfig{
			Model:      true,
			Path:       true,
			Context:    true,
			Git:        true,
			Tools:      true,
			Agents:     true,
			Tasks:      true,
			RateLimits: true,
			Duration:   true,
			Speed:      true,
		},
		Git: GitConfig{
			ShowBranch:      true,
			ShowDirty:       true,
			ShowAheadBehind: true,
			ShowFileStats:   true,
		},
		Tools: ToolsConfig{
			GroupByCategory: true,
			ShowTopN:        5,
			ShowSkills:      true,
			ShowMCP:         true,
		},
	}
}

// Essential returns a config with core metrics only
func Essential() *Config {
	cfg := Default()
	cfg.Preset = "essential"
	cfg.LineLayout = "compact"
	cfg.Display.Tools = false
	cfg.Display.Agents = false
	cfg.Display.RateLimits = false
	cfg.Display.Duration = false
	cfg.Display.Speed = false
	return cfg
}

// Minimal returns a config with minimal information
func Minimal() *Config {
	cfg := Default()
	cfg.Preset = "minimal"
	cfg.LineLayout = "compact"
	cfg.PathLevels = 1
	cfg.Display.Path = false
	cfg.Display.Git = false
	cfg.Display.Tools = false
	cfg.Display.Agents = false
	cfg.Display.Tasks = false
	cfg.Display.RateLimits = false
	cfg.Display.Duration = false
	cfg.Display.Speed = false
	return cfg
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

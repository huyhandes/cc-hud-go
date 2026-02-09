package config

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

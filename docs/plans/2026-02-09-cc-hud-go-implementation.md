# cc-hud-go Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a feature-complete Claude Code statusline tool in Go with enhanced tool tracking, following TDD principles.

**Architecture:** Event-driven using Bubbletea's message-based model. Two data sources (stdin JSON + transcript JSONL) feed a central State, which renders multi-segment JSON output with lipgloss styling. Graceful degradation ensures robustness.

**Tech Stack:** Go 1.25+, Bubbletea, Lipgloss, Bubbles, standard library

---

## Phase 1: Project Setup & Dependencies

### Task 1.1: Initialize Go Module

**Files:**
- Create: `go.mod`
- Create: `.gitignore`

**Step 1: Initialize module**

```bash
go mod init github.com/yourusername/cc-hud-go
```

Expected: `go.mod` created with module declaration

**Step 2: Add dependencies**

```bash
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/charmbracelet/bubbles@latest
```

Expected: Dependencies added to `go.mod`

**Step 3: Create .gitignore**

```gitignore
# Binaries
cc-hud-go
*.exe
*.dll
*.so
*.dylib

# Test coverage
*.out
coverage.html

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
```

**Step 4: Verify setup**

```bash
go mod tidy
go mod verify
```

Expected: No errors

**Step 5: Commit**

```bash
git add go.mod go.sum .gitignore
git commit -m "chore: initialize Go module with Charm dependencies

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 1.2: Create Project Structure

**Files:**
- Create: `main.go`
- Create: `config/config.go`
- Create: `parser/parser.go`
- Create: `state/state.go`
- Create: `segment/segment.go`
- Create: `output/renderer.go`
- Create: `internal/git/git.go`
- Create: `internal/watcher/watcher.go`
- Create: `testdata/.gitkeep`

**Step 1: Create directory structure**

```bash
mkdir -p config parser state segment output internal/git internal/watcher testdata
```

**Step 2: Create placeholder files**

Create `main.go`:
```go
package main

func main() {
	// TODO: Implement
}
```

Create `config/config.go`:
```go
package config

// Config holds all configuration options
type Config struct {
	// TODO: Define fields
}
```

Create `parser/parser.go`:
```go
package parser

// TODO: Define parsers
```

Create `state/state.go`:
```go
package state

// State holds current session data
type State struct {
	// TODO: Define fields
}
```

Create `segment/segment.go`:
```go
package segment

// Segment represents a displayable statusline segment
type Segment interface {
	ID() string
	Render() (string, error)
	Enabled() bool
}
```

Create `output/renderer.go`:
```go
package output

// TODO: Define JSON renderer
```

Create `internal/git/git.go`:
```go
package git

// TODO: Define git operations
```

Create `internal/watcher/watcher.go`:
```go
package watcher

// TODO: Define file watcher
```

Create `testdata/.gitkeep` (empty file)

**Step 3: Verify build**

```bash
go build .
```

Expected: Successful build (binary created)

**Step 4: Commit**

```bash
git add .
git commit -m "chore: create project structure

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Phase 2: Configuration System (TDD)

### Task 2.1: Config Data Structures

**Files:**
- Modify: `config/config.go`
- Create: `config/config_test.go`

**Step 1: Write failing test for default config**

Create `config/config_test.go`:
```go
package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := Default()

	if cfg.Preset != "full" {
		t.Errorf("expected preset 'full', got '%s'", cfg.Preset)
	}

	if cfg.PathLevels != 2 {
		t.Errorf("expected pathLevels 2, got %d", cfg.PathLevels)
	}

	if !cfg.Display.Model {
		t.Error("expected Display.Model to be true")
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./config -v
```

Expected: FAIL - "undefined: Default"

**Step 3: Implement Config struct and Default()**

Modify `config/config.go`:
```go
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
```

**Step 4: Run test to verify it passes**

```bash
go test ./config -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add config/
git commit -m "feat(config): add Config structs and Default()

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 2.2: Config Presets

**Files:**
- Modify: `config/config.go`
- Modify: `config/config_test.go`

**Step 1: Write failing tests for presets**

Add to `config/config_test.go`:
```go
func TestEssentialPreset(t *testing.T) {
	cfg := Essential()

	if cfg.Preset != "essential" {
		t.Errorf("expected preset 'essential', got '%s'", cfg.Preset)
	}

	if cfg.LineLayout != "compact" {
		t.Errorf("expected layout 'compact', got '%s'", cfg.LineLayout)
	}

	if cfg.Display.Tools {
		t.Error("expected Display.Tools to be false in essential preset")
	}

	if !cfg.Display.Model || !cfg.Display.Context {
		t.Error("expected core displays (Model, Context) to be true")
	}
}

func TestMinimalPreset(t *testing.T) {
	cfg := Minimal()

	if cfg.Preset != "minimal" {
		t.Errorf("expected preset 'minimal', got '%s'", cfg.Preset)
	}

	if cfg.Display.Git || cfg.Display.Tasks {
		t.Error("expected Git and Tasks to be false in minimal preset")
	}

	if !cfg.Display.Model || !cfg.Display.Context {
		t.Error("expected core displays (Model, Context) to be true")
	}
}
```

**Step 2: Run tests to verify they fail**

```bash
go test ./config -v
```

Expected: FAIL - "undefined: Essential" and "undefined: Minimal"

**Step 3: Implement preset functions**

Add to `config/config.go`:
```go
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
```

**Step 4: Run tests to verify they pass**

```bash
go test ./config -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add config/
git commit -m "feat(config): add Essential and Minimal presets

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 2.3: Config File Loading

**Files:**
- Modify: `config/config.go`
- Modify: `config/config_test.go`
- Create: `testdata/config_valid.json`
- Create: `testdata/config_invalid.json`

**Step 1: Write failing test for file loading**

Add to `config/config_test.go`:
```go
import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	// Test valid config
	cfg, err := LoadFromFile("../testdata/config_valid.json")
	if err != nil {
		t.Fatalf("failed to load valid config: %v", err)
	}

	if cfg.Preset != "essential" {
		t.Errorf("expected preset 'essential', got '%s'", cfg.Preset)
	}

	// Test missing file (should return defaults)
	cfg, err = LoadFromFile("nonexistent.json")
	if err != nil {
		t.Fatalf("should not error on missing file: %v", err)
	}

	if cfg.Preset != "full" {
		t.Error("expected default preset on missing file")
	}

	// Test invalid JSON (should return defaults)
	cfg, err = LoadFromFile("../testdata/config_invalid.json")
	if err != nil {
		t.Fatalf("should not error on invalid JSON: %v", err)
	}

	if cfg.Preset != "full" {
		t.Error("expected default preset on invalid JSON")
	}
}
```

**Step 2: Create test data files**

Create `testdata/config_valid.json`:
```json
{
  "preset": "essential",
  "layout": "compact",
  "pathLevels": 1,
  "display": {
    "model": true,
    "context": true,
    "git": true
  }
}
```

Create `testdata/config_invalid.json`:
```json
{
  "preset": "essential"
  "invalid": json
}
```

**Step 3: Run test to verify it fails**

```bash
go test ./config -v
```

Expected: FAIL - "undefined: LoadFromFile"

**Step 4: Implement LoadFromFile()**

Add to `config/config.go`:
```go
import (
	"encoding/json"
	"fmt"
	"os"
)

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
```

**Step 5: Run test to verify it passes**

```bash
go test ./config -v
```

Expected: PASS

**Step 6: Commit**

```bash
git add config/ testdata/
git commit -m "feat(config): add LoadFromFile with graceful degradation

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 2.4: Config Validation

**Files:**
- Modify: `config/config.go`
- Modify: `config/config_test.go`

**Step 1: Write failing test for validation**

Add to `config/config_test.go`:
```go
func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *Config
		wantErr  bool
	}{
		{
			name:    "valid config",
			cfg:     Default(),
			wantErr: false,
		},
		{
			name: "pathLevels too low",
			cfg: &Config{
				PathLevels: 0,
			},
			wantErr: true,
		},
		{
			name: "pathLevels too high",
			cfg: &Config{
				PathLevels: 5,
			},
			wantErr: true,
		},
		{
			name: "threshold negative",
			cfg: &Config{
				PathLevels:        2,
				SevenDayThreshold: -10,
			},
			wantErr: true,
		},
		{
			name: "threshold too high",
			cfg: &Config{
				PathLevels:        2,
				SevenDayThreshold: 150,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./config -v
```

Expected: FAIL - "undefined: Config.Validate"

**Step 3: Implement Validate()**

Add to `config/config.go`:
```go
import "errors"

// Validate checks config values are within valid ranges
func (c *Config) Validate() error {
	if c.PathLevels < 1 || c.PathLevels > 3 {
		return errors.New("pathLevels must be between 1 and 3")
	}

	if c.SevenDayThreshold < 0 || c.SevenDayThreshold > 100 {
		return errors.New("sevenDayThreshold must be between 0 and 100")
	}

	return nil
}
```

**Step 4: Update LoadFromFile to validate and fix invalid values**

Modify `LoadFromFile()` in `config/config.go`:
```go
func LoadFromFile(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		fmt.Fprintf(os.Stderr, "warning: failed to read config: %v\n", err)
		return cfg, nil
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to parse config: %v\n", err)
		return Default(), nil
	}

	// Validate and fix invalid values
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: invalid config values, using defaults: %v\n", err)
		return Default(), nil
	}

	return cfg, nil
}
```

**Step 5: Run tests to verify they pass**

```bash
go test ./config -v
```

Expected: PASS

**Step 6: Commit**

```bash
git add config/
git commit -m "feat(config): add Validate() with range checks

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Phase 3: State Management (TDD)

### Task 3.1: Core State Structures

**Files:**
- Modify: `state/state.go`
- Create: `state/state_test.go`

**Step 1: Write failing test for state initialization**

Create `state/state_test.go`:
```go
package state

import (
	"testing"
	"time"
)

func TestNewState(t *testing.T) {
	s := New()

	if s == nil {
		t.Fatal("expected non-nil state")
	}

	if s.Session.StartTime.IsZero() {
		t.Error("expected StartTime to be set")
	}

	if s.Tools.AppTools == nil {
		t.Error("expected AppTools map to be initialized")
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./state -v
```

Expected: FAIL - "undefined: New"

**Step 3: Implement State structs**

Modify `state/state.go`:
```go
package state

import "time"

// State holds all current session data
type State struct {
	Model      ModelInfo
	Context    ContextInfo
	RateLimits RateLimitInfo
	Git        GitInfo
	Tools      ToolsState
	Agents     AgentInfo
	Tasks      TaskInfo
	Session    SessionInfo
}

type ModelInfo struct {
	Name     string
	PlanType string
}

type ContextInfo struct {
	UsedTokens  int
	TotalTokens int
	Percentage  float64
}

type RateLimitInfo struct {
	HourlyUsed    int
	HourlyTotal   int
	SevenDayUsed  int
	SevenDayTotal int
}

type GitInfo struct {
	Branch     string
	DirtyFiles int
	Ahead      int
	Behind     int
	Added      int
	Modified   int
	Deleted    int
}

type ToolsState struct {
	AppTools      map[string]int
	InternalTools map[string]int
	CustomTools   map[string]int
	MCPTools      map[MCPServer]map[string]int
	Skills        map[string]SkillUsage
}

type MCPServer struct {
	Name string
	Type string
}

type SkillUsage struct {
	Count    int
	LastUsed time.Time
	Duration time.Duration
}

type AgentInfo struct {
	ActiveAgent string
	TaskDesc    string
	ElapsedTime time.Duration
}

type TaskInfo struct {
	Pending    int
	InProgress int
	Completed  int
}

type SessionInfo struct {
	StartTime  time.Time
	Duration   time.Duration
	TokenSpeed float64
}

// New creates a new State with initialized maps
func New() *State {
	return &State{
		Tools: ToolsState{
			AppTools:      make(map[string]int),
			InternalTools: make(map[string]int),
			CustomTools:   make(map[string]int),
			MCPTools:      make(map[MCPServer]map[string]int),
			Skills:        make(map[string]SkillUsage),
		},
		Session: SessionInfo{
			StartTime: time.Now(),
		},
	}
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./state -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add state/
git commit -m "feat(state): add State structs and New()

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 3.2: State Update Methods

**Files:**
- Modify: `state/state.go`
- Modify: `state/state_test.go`

**Step 1: Write failing test for UpdateDerived**

Add to `state/state_test.go`:
```go
func TestUpdateDerived(t *testing.T) {
	s := New()

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	s.UpdateDerived()

	if s.Session.Duration == 0 {
		t.Error("expected Duration to be updated")
	}

	if s.Session.Duration < 100*time.Millisecond {
		t.Errorf("expected Duration >= 100ms, got %v", s.Session.Duration)
	}

	// Test percentage calculation
	s.Context.UsedTokens = 50
	s.Context.TotalTokens = 100
	s.UpdateDerived()

	if s.Context.Percentage != 50.0 {
		t.Errorf("expected Percentage 50.0, got %f", s.Context.Percentage)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./state -v
```

Expected: FAIL - "undefined: State.UpdateDerived"

**Step 3: Implement UpdateDerived()**

Add to `state/state.go`:
```go
// UpdateDerived updates calculated fields like duration and percentage
func (s *State) UpdateDerived() {
	// Update session duration
	s.Session.Duration = time.Since(s.Session.StartTime)

	// Update context percentage
	if s.Context.TotalTokens > 0 {
		s.Context.Percentage = float64(s.Context.UsedTokens) / float64(s.Context.TotalTokens) * 100.0
	}
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./state -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add state/
git commit -m "feat(state): add UpdateDerived() for calculated fields

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Phase 4: Stdin Parser (TDD)

### Task 4.1: Basic Stdin Parsing

**Files:**
- Modify: `parser/parser.go`
- Create: `parser/stdin_test.go`
- Create: `testdata/stdin_basic.json`

**Step 1: Write failing test for stdin parsing**

Create `parser/stdin_test.go`:
```go
package parser

import (
	"testing"

	"github.com/yourusername/cc-hud-go/state"
)

func TestParseStdin(t *testing.T) {
	input := `{
		"model": "claude-sonnet-4.5",
		"planType": "Pro",
		"context": {
			"used": 1500,
			"total": 200000
		}
	}`

	s := state.New()
	err := ParseStdin([]byte(input), s)

	if err != nil {
		t.Fatalf("ParseStdin failed: %v", err)
	}

	if s.Model.Name != "claude-sonnet-4.5" {
		t.Errorf("expected model 'claude-sonnet-4.5', got '%s'", s.Model.Name)
	}

	if s.Model.PlanType != "Pro" {
		t.Errorf("expected plan 'Pro', got '%s'", s.Model.PlanType)
	}

	if s.Context.UsedTokens != 1500 {
		t.Errorf("expected UsedTokens 1500, got %d", s.Context.UsedTokens)
	}

	if s.Context.TotalTokens != 200000 {
		t.Errorf("expected TotalTokens 200000, got %d", s.Context.TotalTokens)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./parser -v
```

Expected: FAIL - "undefined: ParseStdin"

**Step 3: Implement ParseStdin()**

Modify `parser/parser.go`:
```go
package parser

import (
	"encoding/json"

	"github.com/yourusername/cc-hud-go/state"
)

// StdinData represents the JSON structure from Claude Code
type StdinData struct {
	Model    string `json:"model"`
	PlanType string `json:"planType"`
	Context  struct {
		Used  int `json:"used"`
		Total int `json:"total"`
	} `json:"context"`
	RateLimits *struct {
		HourlyUsed    int `json:"hourlyUsed"`
		HourlyTotal   int `json:"hourlyTotal"`
		SevenDayUsed  int `json:"sevenDayUsed"`
		SevenDayTotal int `json:"sevenDayTotal"`
	} `json:"rateLimits,omitempty"`
}

// ParseStdin parses stdin JSON and updates state
func ParseStdin(data []byte, s *state.State) error {
	var stdin StdinData
	if err := json.Unmarshal(data, &stdin); err != nil {
		return err
	}

	// Update model info
	s.Model.Name = stdin.Model
	s.Model.PlanType = stdin.PlanType

	// Update context
	s.Context.UsedTokens = stdin.Context.Used
	s.Context.TotalTokens = stdin.Context.Total

	// Update rate limits if present
	if stdin.RateLimits != nil {
		s.RateLimits.HourlyUsed = stdin.RateLimits.HourlyUsed
		s.RateLimits.HourlyTotal = stdin.RateLimits.HourlyTotal
		s.RateLimits.SevenDayUsed = stdin.RateLimits.SevenDayUsed
		s.RateLimits.SevenDayTotal = stdin.RateLimits.SevenDayTotal
	}

	return nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./parser -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add parser/
git commit -m "feat(parser): add ParseStdin() for Claude Code JSON

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 4.2: Handle Malformed Stdin

**Files:**
- Modify: `parser/stdin_test.go`

**Step 1: Write failing test for error handling**

Add to `parser/stdin_test.go`:
```go
func TestParseStdinInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid JSON",
			input: `{"model": invalid}`,
		},
		{
			name:  "empty input",
			input: ``,
		},
		{
			name:  "partial JSON",
			input: `{"model": "test"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := state.New()
			err := ParseStdin([]byte(tt.input), s)

			if err == nil {
				t.Error("expected error for invalid input")
			}
		})
	}
}
```

**Step 2: Run test to verify it passes (error handling already works)**

```bash
go test ./parser -v
```

Expected: PASS (json.Unmarshal returns errors for invalid JSON)

**Step 3: Commit**

```bash
git add parser/
git commit -m "test(parser): add error handling tests for ParseStdin

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Phase 5: Basic Segments with Lipgloss (TDD)

### Task 5.1: Segment Interface and Registry

**Files:**
- Modify: `segment/segment.go`
- Create: `segment/segment_test.go`

**Step 1: Write failing test for segment registry**

Create `segment/segment_test.go`:
```go
package segment

import (
	"testing"

	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

func TestRegistry(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	segments := All()

	if len(segments) == 0 {
		t.Error("expected at least one segment")
	}

	// Check that segments implement interface
	for _, seg := range segments {
		if seg.ID() == "" {
			t.Error("segment ID should not be empty")
		}

		// Should be able to check if enabled
		_ = seg.Enabled(cfg)

		// Should be able to render
		_, err := seg.Render(s, cfg)
		if err != nil {
			t.Errorf("segment %s render failed: %v", seg.ID(), err)
		}
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./segment -v
```

Expected: FAIL - "undefined: All"

**Step 3: Update Segment interface and add registry**

Modify `segment/segment.go`:
```go
package segment

import (
	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

// Segment represents a displayable statusline segment
type Segment interface {
	ID() string
	Render(s *state.State, cfg *config.Config) (string, error)
	Enabled(cfg *config.Config) bool
}

// All returns all available segments
func All() []Segment {
	return []Segment{
		// TODO: Add segments as we build them
	}
}
```

**Step 4: Run test to verify it passes (with empty registry)**

```bash
go test ./segment -v
```

Expected: FAIL - "expected at least one segment"

**Step 5: Add a dummy segment for testing**

Modify `segment/segment.go`:
```go
type dummySegment struct{}

func (d dummySegment) ID() string { return "dummy" }
func (d dummySegment) Enabled(cfg *config.Config) bool { return true }
func (d dummySegment) Render(s *state.State, cfg *config.Config) (string, error) {
	return "dummy", nil
}

func All() []Segment {
	return []Segment{
		dummySegment{},
	}
}
```

**Step 6: Run test to verify it passes**

```bash
go test ./segment -v
```

Expected: PASS

**Step 7: Commit**

```bash
git add segment/
git commit -m "feat(segment): add Segment interface and registry

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 5.2: Model Segment

**Files:**
- Create: `segment/model.go`
- Create: `segment/model_test.go`
- Modify: `segment/segment.go`

**Step 1: Write failing test for ModelSegment**

Create `segment/model_test.go`:
```go
package segment

import (
	"strings"
	"testing"

	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

func TestModelSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Model.Name = "claude-sonnet-4.5"
	s.Model.PlanType = "Pro"

	seg := &ModelSegment{}

	if seg.ID() != "model" {
		t.Errorf("expected ID 'model', got '%s'", seg.ID())
	}

	if !seg.Enabled(cfg) {
		t.Error("expected segment to be enabled by default")
	}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "claude-sonnet-4.5") {
		t.Errorf("expected output to contain model name, got '%s'", output)
	}

	if !strings.Contains(output, "Pro") {
		t.Errorf("expected output to contain plan type, got '%s'", output)
	}
}

func TestModelSegmentDisabled(t *testing.T) {
	cfg := config.Default()
	cfg.Display.Model = false

	seg := &ModelSegment{}

	if seg.Enabled(cfg) {
		t.Error("expected segment to be disabled")
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./segment -v
```

Expected: FAIL - "undefined: ModelSegment"

**Step 3: Implement ModelSegment**

Create `segment/model.go`:
```go
package segment

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

var (
	modelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // Blue
	planStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green
)

type ModelSegment struct{}

func (m *ModelSegment) ID() string {
	return "model"
}

func (m *ModelSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Model
}

func (m *ModelSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	if s.Model.Name == "" {
		return "", nil
	}

	model := modelStyle.Render(s.Model.Name)

	if s.Model.PlanType != "" {
		plan := planStyle.Render(fmt.Sprintf("[%s]", s.Model.PlanType))
		return fmt.Sprintf("%s %s", model, plan), nil
	}

	return model, nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./segment -run TestModel -v
```

Expected: PASS

**Step 5: Update registry to include ModelSegment**

Modify `segment/segment.go`:
```go
func All() []Segment {
	return []Segment{
		&ModelSegment{},
	}
}
```

**Step 6: Commit**

```bash
git add segment/
git commit -m "feat(segment): add ModelSegment with lipgloss styling

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 5.3: Context Segment with Color Thresholds

**Files:**
- Create: `segment/context.go`
- Create: `segment/context_test.go`
- Modify: `segment/segment.go`

**Step 1: Write failing test for ContextSegment**

Create `segment/context_test.go`:
```go
package segment

import (
	"strings"
	"testing"

	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

func TestContextSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	seg := &ContextSegment{}

	if seg.ID() != "context" {
		t.Errorf("expected ID 'context', got '%s'", seg.ID())
	}

	// Test green threshold (<70%)
	s.Context.UsedTokens = 50000
	s.Context.TotalTokens = 200000
	s.UpdateDerived()

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "25") {
		t.Errorf("expected percentage in output, got '%s'", output)
	}

	// Test yellow threshold (70-90%)
	s.Context.UsedTokens = 160000
	s.UpdateDerived()

	output, err = seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "80") {
		t.Errorf("expected percentage in output, got '%s'", output)
	}

	// Test red threshold (>90%)
	s.Context.UsedTokens = 190000
	s.UpdateDerived()

	output, err = seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "95") {
		t.Errorf("expected percentage in output, got '%s'", output)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./segment -run TestContext -v
```

Expected: FAIL - "undefined: ContextSegment"

**Step 3: Implement ContextSegment**

Create `segment/context.go`:
```go
package segment

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

var (
	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	yellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	redStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

type ContextSegment struct{}

func (c *ContextSegment) ID() string {
	return "context"
}

func (c *ContextSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Context
}

func (c *ContextSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	if s.Context.TotalTokens == 0 {
		return "", nil
	}

	percentage := s.Context.Percentage

	// Choose color based on thresholds
	var style lipgloss.Style
	if percentage < 70 {
		style = greenStyle
	} else if percentage < 90 {
		style = yellowStyle
	} else {
		style = redStyle
	}

	// Build progress bar
	barWidth := 10
	filled := int(percentage / 10)
	if filled > barWidth {
		filled = barWidth
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	var display string
	if cfg.ContextValue == "tokens" {
		display = fmt.Sprintf("%d/%d", s.Context.UsedTokens, s.Context.TotalTokens)
	} else {
		display = fmt.Sprintf("%.0f%%", percentage)
	}

	return style.Render(fmt.Sprintf("[%s] %s", bar, display)), nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./segment -run TestContext -v
```

Expected: PASS

**Step 5: Update registry**

Modify `segment/segment.go`:
```go
func All() []Segment {
	return []Segment{
		&ModelSegment{},
		&ContextSegment{},
	}
}
```

**Step 6: Commit**

```bash
git add segment/
git commit -m "feat(segment): add ContextSegment with color thresholds

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Phase 6: Git Integration (TDD)

### Task 6.1: Git Command Wrappers

**Files:**
- Modify: `internal/git/git.go`
- Create: `internal/git/git_test.go`

**Step 1: Write failing test for GetBranch**

Create `internal/git/git_test.go`:
```go
package git

import (
	"testing"
)

func TestGetBranch(t *testing.T) {
	// This test requires a git repo
	// For now, test that it doesn't panic
	branch, err := GetBranch()

	// In a non-git directory, expect error
	if err != nil && branch != "" {
		t.Error("expected empty branch on error")
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/git -v
```

Expected: FAIL - "undefined: GetBranch"

**Step 3: Implement GetBranch()**

Modify `internal/git/git.go`:
```go
package git

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
)

// GetBranch returns the current git branch name
func GetBranch() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/git -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/git/
git commit -m "feat(git): add GetBranch() with timeout

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 6.2: Git Status Parsing

**Files:**
- Modify: `internal/git/git.go`
- Modify: `internal/git/git_test.go`

**Step 1: Write failing test for GetStatus**

Add to `internal/git/git_test.go`:
```go
func TestGetStatus(t *testing.T) {
	// Test in current repo (should have git)
	status, err := GetStatus()

	// If we're in a git repo, should not error
	// In a non-git directory, expect error
	if err != nil && status != nil {
		t.Error("expected nil status on error")
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/git -v
```

Expected: FAIL - "undefined: GetStatus"

**Step 3: Implement GetStatus()**

Add to `internal/git/git.go`:
```go
import "strconv"

// Status holds git status information
type Status struct {
	DirtyFiles int
	Ahead      int
	Behind     int
	Added      int
	Modified   int
	Deleted    int
}

// GetStatus returns git status information
func GetStatus() (*Status, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	status := &Status{}

	// Get ahead/behind
	cmd := exec.CommandContext(ctx, "git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err == nil {
		parts := strings.Fields(out.String())
		if len(parts) == 2 {
			status.Ahead, _ = strconv.Atoi(parts[0])
			status.Behind, _ = strconv.Atoi(parts[1])
		}
	}

	// Get file stats
	cmd = exec.CommandContext(ctx, "git", "status", "--porcelain")
	out.Reset()
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}

		status.DirtyFiles++

		code := line[:2]
		if strings.Contains(code, "A") {
			status.Added++
		}
		if strings.Contains(code, "M") {
			status.Modified++
		}
		if strings.Contains(code, "D") {
			status.Deleted++
		}
	}

	return status, nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/git -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/git/
git commit -m "feat(git): add GetStatus() for file stats

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 6.3: Git Segment

**Files:**
- Create: `segment/git.go`
- Create: `segment/git_test.go`
- Modify: `segment/segment.go`

**Step 1: Write failing test for GitSegment**

Create `segment/git_test.go`:
```go
package segment

import (
	"strings"
	"testing"

	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

func TestGitSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Git.Branch = "main"
	s.Git.DirtyFiles = 3
	s.Git.Ahead = 2
	s.Git.Behind = 1

	seg := &GitSegment{}

	if seg.ID() != "git" {
		t.Errorf("expected ID 'git', got '%s'", seg.ID())
	}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "main") {
		t.Errorf("expected branch name in output, got '%s'", output)
	}

	if !strings.Contains(output, "3") {
		t.Errorf("expected dirty files count in output, got '%s'", output)
	}
}

func TestGitSegmentNoBranch(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	seg := &GitSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if output != "" {
		t.Errorf("expected empty output with no branch, got '%s'", output)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./segment -run TestGit -v
```

Expected: FAIL - "undefined: GitSegment"

**Step 3: Implement GitSegment**

Create `segment/git.go`:
```go
package segment

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

var (
	branchStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13")) // Magenta
	dirtyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // Yellow
)

type GitSegment struct{}

func (g *GitSegment) ID() string {
	return "git"
}

func (g *GitSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Git
}

func (g *GitSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	if s.Git.Branch == "" {
		return "", nil
	}

	var parts []string

	// Branch name
	if cfg.Git.ShowBranch {
		parts = append(parts, branchStyle.Render(s.Git.Branch))
	}

	// Dirty indicator
	if cfg.Git.ShowDirty && s.Git.DirtyFiles > 0 {
		parts = append(parts, dirtyStyle.Render(fmt.Sprintf("✗%d", s.Git.DirtyFiles)))
	}

	// Ahead/behind
	if cfg.Git.ShowAheadBehind {
		if s.Git.Ahead > 0 {
			parts = append(parts, fmt.Sprintf("↑%d", s.Git.Ahead))
		}
		if s.Git.Behind > 0 {
			parts = append(parts, fmt.Sprintf("↓%d", s.Git.Behind))
		}
	}

	// File stats
	if cfg.Git.ShowFileStats {
		if s.Git.Added > 0 {
			parts = append(parts, fmt.Sprintf("+%d", s.Git.Added))
		}
		if s.Git.Modified > 0 {
			parts = append(parts, fmt.Sprintf("~%d", s.Git.Modified))
		}
		if s.Git.Deleted > 0 {
			parts = append(parts, fmt.Sprintf("-%d", s.Git.Deleted))
		}
	}

	return strings.Join(parts, " "), nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./segment -run TestGit -v
```

Expected: PASS

**Step 5: Update registry**

Modify `segment/segment.go`:
```go
func All() []Segment {
	return []Segment{
		&ModelSegment{},
		&ContextSegment{},
		&GitSegment{},
	}
}
```

**Step 6: Commit**

```bash
git add segment/
git commit -m "feat(segment): add GitSegment with configurable displays

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Phase 7: Transcript Parser & Tool Tracking (TDD)

### Task 7.1: Tool Categorization Logic

**Files:**
- Modify: `parser/parser.go`
- Create: `parser/transcript_test.go`

**Step 1: Write failing test for tool categorization**

Create `parser/transcript_test.go`:
```go
package parser

import (
	"testing"
)

func TestCategorizeTool(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		want     ToolCategory
	}{
		{"App tool", "Read", CategoryApp},
		{"App tool lowercase", "read", CategoryApp},
		{"Bash", "Bash", CategoryInternal},
		{"MCP tool", "mcp__claude_ai_Atlassian__getConfluencePage", CategoryMCP},
		{"Skill", "Skill", CategorySkill},
		{"Custom", "MyCustomTool", CategoryCustom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CategorizeTool(tt.toolName)
			if got != tt.want {
				t.Errorf("CategorizeTool(%s) = %v, want %v", tt.toolName, got, tt.want)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./parser -run TestCategorize -v
```

Expected: FAIL - "undefined: CategorizeTool"

**Step 3: Implement CategorizeTool()**

Add to `parser/parser.go`:
```go
import "strings"

type ToolCategory int

const (
	CategoryApp ToolCategory = iota
	CategoryInternal
	CategoryCustom
	CategoryMCP
	CategorySkill
)

var appTools = map[string]bool{
	"read":      true,
	"write":     true,
	"edit":      true,
	"bash":      true,
	"glob":      true,
	"grep":      true,
	"task":      true,
	"webfetch":  true,
	"websearch": true,
}

// CategorizeTool determines the category of a tool by name
func CategorizeTool(name string) ToolCategory {
	lower := strings.ToLower(name)

	// Check for MCP pattern
	if strings.HasPrefix(lower, "mcp__") {
		return CategoryMCP
	}

	// Check for Skill
	if lower == "skill" {
		return CategorySkill
	}

	// Check for app tools
	if appTools[lower] {
		return CategoryApp
	}

	// Check for internal (Bash is special)
	if lower == "bash" {
		return CategoryInternal
	}

	// Everything else is custom
	return CategoryCustom
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./parser -run TestCategorize -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add parser/
git commit -m "feat(parser): add CategorizeTool() for tool classification

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 7.2: Transcript JSONL Parsing

**Files:**
- Modify: `parser/parser.go`
- Modify: `parser/transcript_test.go`

**Step 1: Write failing test for transcript parsing**

Add to `parser/transcript_test.go`:
```go
import (
	"github.com/yourusername/cc-hud-go/state"
)

func TestParseTranscriptLine(t *testing.T) {
	line := `{"type":"tool_use","name":"Read","id":"tool_123"}`

	s := state.New()
	err := ParseTranscriptLine([]byte(line), s)

	if err != nil {
		t.Fatalf("ParseTranscriptLine failed: %v", err)
	}

	if s.Tools.AppTools["Read"] != 1 {
		t.Errorf("expected Read count 1, got %d", s.Tools.AppTools["Read"])
	}
}

func TestParseTranscriptLineMCP(t *testing.T) {
	line := `{"type":"tool_use","name":"mcp__claude_ai_Atlassian__getConfluencePage"}`

	s := state.New()
	err := ParseTranscriptLine([]byte(line), s)

	if err != nil {
		t.Fatalf("ParseTranscriptLine failed: %v", err)
	}

	// Check MCP tools map
	found := false
	for server, tools := range s.Tools.MCPTools {
		if server.Name == "claude_ai_Atlassian" {
			if tools["getConfluencePage"] == 1 {
				found = true
			}
		}
	}

	if !found {
		t.Error("expected MCP tool to be tracked")
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./parser -run TestParseTranscript -v
```

Expected: FAIL - "undefined: ParseTranscriptLine"

**Step 3: Implement ParseTranscriptLine()**

Add to `parser/parser.go`:
```go
import (
	"github.com/yourusername/cc-hud-go/state"
	"time"
)

type TranscriptLine struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// ParseTranscriptLine parses a single JSONL line and updates state
func ParseTranscriptLine(data []byte, s *state.State) error {
	var line TranscriptLine
	if err := json.Unmarshal(data, &line); err != nil {
		return err
	}

	// Only process tool_use events
	if line.Type != "tool_use" {
		return nil
	}

	category := CategorizeTool(line.Name)

	switch category {
	case CategoryApp:
		s.Tools.AppTools[line.Name]++

	case CategoryInternal:
		s.Tools.InternalTools[line.Name]++

	case CategoryCustom:
		s.Tools.CustomTools[line.Name]++

	case CategoryMCP:
		// Parse MCP tool name: mcp__<server>__<tool>
		parts := strings.Split(line.Name, "__")
		if len(parts) >= 3 {
			server := state.MCPServer{
				Name: parts[1],
				Type: "mcp",
			}

			if s.Tools.MCPTools[server] == nil {
				s.Tools.MCPTools[server] = make(map[string]int)
			}

			toolName := strings.Join(parts[2:], "__")
			s.Tools.MCPTools[server][toolName]++
		}

	case CategorySkill:
		// Skills need additional parsing from the tool parameters
		// For now, just count as app tool
		s.Tools.AppTools["Skill"]++
	}

	return nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./parser -run TestParseTranscript -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add parser/
git commit -m "feat(parser): add ParseTranscriptLine() with tool tracking

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 7.3: Tools Segment

**Files:**
- Create: `segment/tools.go`
- Create: `segment/tools_test.go`
- Modify: `segment/segment.go`

**Step 1: Write failing test for ToolsSegment**

Create `segment/tools_test.go`:
```go
package segment

import (
	"strings"
	"testing"

	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

func TestToolsSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	s.Tools.AppTools["Read"] = 15
	s.Tools.AppTools["Edit"] = 8
	s.Tools.MCPTools[state.MCPServer{Name: "github", Type: "mcp"}] = map[string]int{
		"create_issue": 2,
	}
	s.Tools.Skills["brainstorming"] = state.SkillUsage{Count: 1}

	seg := &ToolsSegment{}

	if seg.ID() != "tools" {
		t.Errorf("expected ID 'tools', got '%s'", seg.ID())
	}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "Tools:") {
		t.Errorf("expected 'Tools:' prefix in output, got '%s'", output)
	}

	// Should show category counts
	if !strings.Contains(output, "App:") {
		t.Errorf("expected 'App:' category in output, got '%s'", output)
	}
}

func TestToolsSegmentEmpty(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	seg := &ToolsSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if output != "" {
		t.Errorf("expected empty output with no tools, got '%s'", output)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./segment -run TestTools -v
```

Expected: FAIL - "undefined: ToolsSegment"

**Step 3: Implement ToolsSegment**

Create `segment/tools.go`:
```go
package segment

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

var toolsStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14")) // Cyan

type ToolsSegment struct{}

func (t *ToolsSegment) ID() string {
	return "tools"
}

func (t *ToolsSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Tools
}

func (t *ToolsSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	// Count totals
	appTotal := 0
	for _, count := range s.Tools.AppTools {
		appTotal += count
	}

	internalTotal := 0
	for _, count := range s.Tools.InternalTools {
		internalTotal += count
	}

	customTotal := 0
	for _, count := range s.Tools.CustomTools {
		customTotal += count
	}

	mcpTotal := 0
	for _, tools := range s.Tools.MCPTools {
		for _, count := range tools {
			mcpTotal += count
		}
	}

	skillsTotal := 0
	for _, usage := range s.Tools.Skills {
		skillsTotal += usage.Count
	}

	total := appTotal + internalTotal + customTotal + mcpTotal + skillsTotal

	if total == 0 {
		return "", nil
	}

	// Build output
	var parts []string

	if cfg.Tools.GroupByCategory {
		if appTotal > 0 {
			parts = append(parts, fmt.Sprintf("App:%d", appTotal))
		}
		if mcpTotal > 0 && cfg.Tools.ShowMCP {
			parts = append(parts, fmt.Sprintf("MCP:%d", mcpTotal))
		}
		if skillsTotal > 0 && cfg.Tools.ShowSkills {
			parts = append(parts, fmt.Sprintf("Skills:%d", skillsTotal))
		}
		if customTotal > 0 {
			parts = append(parts, fmt.Sprintf("Custom:%d", customTotal))
		}

		return toolsStyle.Render(fmt.Sprintf("Tools: %d (%s)", total, strings.Join(parts, " "))), nil
	}

	return toolsStyle.Render(fmt.Sprintf("Tools: %d", total)), nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./segment -run TestTools -v
```

Expected: PASS

**Step 5: Update registry**

Modify `segment/segment.go`:
```go
func All() []Segment {
	return []Segment{
		&ModelSegment{},
		&ContextSegment{},
		&GitSegment{},
		&ToolsSegment{},
	}
}
```

**Step 6: Commit**

```bash
git add segment/
git commit -m "feat(segment): add ToolsSegment with category grouping

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Phase 8: Remaining Segments (TDD)

### Task 8.1: Tasks Segment

**Files:**
- Create: `segment/tasks.go`
- Create: `segment/tasks_test.go`
- Modify: `segment/segment.go`

**Step 1: Write failing test**

Create `segment/tasks_test.go`:
```go
package segment

import (
	"strings"
	"testing"

	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

func TestTasksSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Tasks.Completed = 2
	s.Tasks.InProgress = 1
	s.Tasks.Pending = 2

	seg := &TasksSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "2/5") {
		t.Errorf("expected progress ratio in output, got '%s'", output)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./segment -run TestTasks -v
```

Expected: FAIL - "undefined: TasksSegment"

**Step 3: Implement TasksSegment**

Create `segment/tasks.go`:
```go
package segment

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

var tasksStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green

type TasksSegment struct{}

func (t *TasksSegment) ID() string {
	return "tasks"
}

func (t *TasksSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Tasks
}

func (t *TasksSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	total := s.Tasks.Completed + s.Tasks.InProgress + s.Tasks.Pending

	if total == 0 {
		return "", nil
	}

	return tasksStyle.Render(fmt.Sprintf("Tasks: %d/%d ✓ (%d pending)",
		s.Tasks.Completed, total, s.Tasks.Pending)), nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./segment -run TestTasks -v
```

Expected: PASS

**Step 5: Update registry and commit**

Modify `segment/segment.go`, then:
```bash
git add segment/
git commit -m "feat(segment): add TasksSegment

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 8.2: Agent Segment

**Files:**
- Create: `segment/agent.go`
- Create: `segment/agent_test.go`
- Modify: `segment/segment.go`

**Step 1: Write failing test**

Create `segment/agent_test.go`:
```go
package segment

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

func TestAgentSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Agents.ActiveAgent = "explore"
	s.Agents.TaskDesc = "Searching codebase"
	s.Agents.ElapsedTime = 2*time.Minute + 15*time.Second

	seg := &AgentSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "explore") {
		t.Errorf("expected agent name in output, got '%s'", output)
	}

	if !strings.Contains(output, "2m") {
		t.Errorf("expected elapsed time in output, got '%s'", output)
	}
}
```

**Step 2: Implement AgentSegment**

Create `segment/agent.go`:
```go
package segment

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

var agentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13")) // Magenta

type AgentSegment struct{}

func (a *AgentSegment) ID() string {
	return "agent"
}

func (a *AgentSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Agents
}

func (a *AgentSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	if s.Agents.ActiveAgent == "" {
		return "", nil
	}

	elapsed := s.Agents.ElapsedTime.Round(time.Second)

	output := fmt.Sprintf("Agent: %s [%s]", s.Agents.ActiveAgent, elapsed)

	if s.Agents.TaskDesc != "" {
		output += fmt.Sprintf(" - %s", s.Agents.TaskDesc)
	}

	return agentStyle.Render(output), nil
}
```

**Step 3: Test, update registry, and commit**

```bash
go test ./segment -run TestAgent -v
git add segment/
git commit -m "feat(segment): add AgentSegment

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 8.3: Rate Limit Segment

**Files:**
- Create: `segment/ratelimit.go`
- Create: `segment/ratelimit_test.go`
- Modify: `segment/segment.go`

**Step 1: Write failing test**

Create `segment/ratelimit_test.go`:
```go
package segment

import (
	"strings"
	"testing"

	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

func TestRateLimitSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.RateLimits.HourlyUsed = 60
	s.RateLimits.HourlyTotal = 100
	s.RateLimits.SevenDayUsed = 450
	s.RateLimits.SevenDayTotal = 1000

	seg := &RateLimitSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "Hourly:") {
		t.Errorf("expected hourly rate in output, got '%s'", output)
	}
}
```

**Step 2: Implement RateLimitSegment**

Create `segment/ratelimit.go`:
```go
package segment

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/state"
)

type RateLimitSegment struct{}

func (r *RateLimitSegment) ID() string {
	return "ratelimit"
}

func (r *RateLimitSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.RateLimits
}

func (r *RateLimitSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	if s.RateLimits.HourlyTotal == 0 {
		return "", nil
	}

	var parts []string

	// Hourly
	hourlyPct := float64(s.RateLimits.HourlyUsed) / float64(s.RateLimits.HourlyTotal) * 100
	hourlyBar := buildBar(hourlyPct, 5)
	parts = append(parts, fmt.Sprintf("Hourly: [%s] %.0f%%", hourlyBar, hourlyPct))

	// 7-day (only show if above threshold)
	if s.RateLimits.SevenDayTotal > 0 {
		sevenDayPct := float64(s.RateLimits.SevenDayUsed) / float64(s.RateLimits.SevenDayTotal) * 100

		if sevenDayPct >= float64(cfg.SevenDayThreshold) {
			sevenDayBar := buildBar(sevenDayPct, 5)
			parts = append(parts, fmt.Sprintf("7-day: [%s] %.0f%%", sevenDayBar, sevenDayPct))
		}
	}

	return strings.Join(parts, " | "), nil
}

func buildBar(percentage float64, width int) string {
	filled := int(percentage / 100 * float64(width))
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}
```

**Step 3: Test, update registry, and commit**

```bash
go test ./segment -run TestRateLimit -v
git add segment/
git commit -m "feat(segment): add RateLimitSegment with bars

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Phase 9: JSON Output Renderer (TDD)

### Task 9.1: Renderer Implementation

**Files:**
- Modify: `output/renderer.go`
- Create: `output/renderer_test.go`

**Step 1: Write failing test**

Create `output/renderer_test.go`:
```go
package output

import (
	"encoding/json"
	"testing"

	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/segment"
	"github.com/yourusername/cc-hud-go/state"
)

func TestRenderJSON(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Model.Name = "claude-sonnet-4.5"
	s.Model.PlanType = "Pro"

	segments := []segment.Segment{
		&segment.ModelSegment{},
	}

	output, err := RenderJSON(s, cfg, segments)
	if err != nil {
		t.Fatalf("RenderJSON failed: %v", err)
	}

	// Parse output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	// Check structure
	if _, ok := result["segments"]; !ok {
		t.Error("expected 'segments' key in output")
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./output -v
```

Expected: FAIL - "undefined: RenderJSON"

**Step 3: Implement RenderJSON**

Modify `output/renderer.go`:
```go
package output

import (
	"encoding/json"

	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/segment"
	"github.com/yourusername/cc-hud-go/state"
)

type SegmentOutput struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Output struct {
	Segments []SegmentOutput `json:"segments"`
}

// RenderJSON renders all enabled segments as JSON
func RenderJSON(s *state.State, cfg *config.Config, segments []segment.Segment) (string, error) {
	output := Output{
		Segments: []SegmentOutput{},
	}

	for _, seg := range segments {
		if !seg.Enabled(cfg) {
			continue
		}

		text, err := seg.Render(s, cfg)
		if err != nil {
			// Log error but continue
			continue
		}

		// Skip empty segments
		if text == "" {
			continue
		}

		output.Segments = append(output.Segments, SegmentOutput{
			ID:   seg.ID(),
			Text: text,
		})
	}

	data, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./output -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add output/
git commit -m "feat(output): add RenderJSON for segment output

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Phase 10: Bubbletea Integration

### Task 10.1: Bubbletea Model

**Files:**
- Modify: `main.go`

**Step 1: Write basic Bubbletea model structure**

Modify `main.go`:
```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/cc-hud-go/config"
	"github.com/yourusername/cc-hud-go/segment"
	"github.com/yourusername/cc-hud-go/state"
)

type model struct {
	state    *state.State
	config   *config.Config
	segments []segment.Segment
}

func initialModel() model {
	return model{
		state:    state.New(),
		config:   config.Default(),
		segments: segment.All(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return ""
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```

**Step 2: Build and test**

```bash
go build .
./cc-hud-go
```

Expected: Program starts and exits on Ctrl+C

**Step 3: Commit**

```bash
git add main.go
git commit -m "feat: add basic Bubbletea model structure

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 10.2: Stdin Reading Command

**Files:**
- Modify: `main.go`

**Step 1: Add stdin reading**

Modify `main.go`:
```go
import (
	"bufio"
	"io"

	"github.com/yourusername/cc-hud-go/parser"
)

type stdinMsg struct {
	data []byte
	err  error
}

func readStdinCmd() tea.Cmd {
	return func() tea.Msg {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadBytes('\n')
		return stdinMsg{data: line, err: err}
	}
}

func (m model) Init() tea.Cmd {
	return readStdinCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case stdinMsg:
		if msg.err != nil {
			if msg.err == io.EOF {
				return m, tea.Quit
			}
			// Log error but continue
			fmt.Fprintf(os.Stderr, "stdin error: %v\n", msg.err)
			return m, readStdinCmd()
		}

		// Parse and update state
		if err := parser.ParseStdin(msg.data, m.state); err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		}

		// Output JSON and read next line
		return m, tea.Batch(outputCmd(m), readStdinCmd())
	}

	return m, nil
}
```

**Step 2: Add output command**

Add to `main.go`:
```go
import "github.com/yourusername/cc-hud-go/output"

func outputCmd(m model) tea.Cmd {
	return func() tea.Msg {
		m.state.UpdateDerived()

		json, err := output.RenderJSON(m.state, m.config, m.segments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "render error: %v\n", err)
			return nil
		}

		fmt.Println(json)
		return nil
	}
}
```

**Step 3: Build and test with mock stdin**

```bash
go build .
echo '{"model":"test","planType":"Pro","context":{"used":100,"total":1000}}' | ./cc-hud-go
```

Expected: JSON output with model segment

**Step 4: Commit**

```bash
git add main.go
git commit -m "feat: add stdin reading and JSON output

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 10.3: Transcript Watcher

**Files:**
- Modify: `internal/watcher/watcher.go`
- Modify: `main.go`

**Step 1: Implement simple file watcher**

Modify `internal/watcher/watcher.go`:
```go
package watcher

import (
	"bufio"
	"os"
	"time"
)

// Watch watches a file and sends new lines on the channel
func Watch(path string, lines chan<- string, stop <-chan struct{}) error {
	// Wait for file to exist
	for {
		if _, err := os.Stat(path); err == nil {
			break
		}

		select {
		case <-stop:
			return nil
		case <-time.After(1 * time.Second):
			// Retry
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Seek to end
	file.Seek(0, 2)

	scanner := bufio.NewScanner(file)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return nil
		case <-ticker.C:
			for scanner.Scan() {
				lines <- scanner.Text()
			}
		}
	}
}
```

**Step 2: Integrate into main**

Modify `main.go`:
```go
import (
	"path/filepath"
	"github.com/yourusername/cc-hud-go/internal/watcher"
)

type transcriptMsg struct {
	line string
}

func watchTranscriptCmd() tea.Cmd {
	return func() tea.Msg {
		// Find transcript file (simplified - in production would query Claude Code)
		home, _ := os.UserHomeDir()
		transcriptPath := filepath.Join(home, ".claude", "transcript.jsonl")

		lines := make(chan string, 10)
		stop := make(chan struct{})

		go func() {
			watcher.Watch(transcriptPath, lines, stop)
		}()

		// Read one line
		line := <-lines
		return transcriptMsg{line: line}
	}
}

// Add to Update():
case transcriptMsg:
	if err := parser.ParseTranscriptLine([]byte(msg.line), m.state); err != nil {
		fmt.Fprintf(os.Stderr, "transcript parse error: %v\n", err)
	}
	return m, tea.Batch(outputCmd(m), watchTranscriptCmd())
```

**Step 3: Build and verify**

```bash
go build .
```

Expected: Successful build

**Step 4: Commit**

```bash
git add main.go internal/watcher/
git commit -m "feat: add transcript file watcher

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 10.4: Config Loading and Git Integration

**Files:**
- Modify: `main.go`

**Step 1: Load config on startup**

Modify `initialModel()` in `main.go`:
```go
func initialModel() model {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".claude", "cc-hud-go", "config.json")

	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		cfg = config.Default()
	}

	return model{
		state:    state.New(),
		config:   cfg,
		segments: segment.All(),
	}
}
```

**Step 2: Add periodic git refresh**

Add to `main.go`:
```go
import (
	"github.com/yourusername/cc-hud-go/internal/git"
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update Init():
func (m model) Init() tea.Cmd {
	return tea.Batch(readStdinCmd(), watchTranscriptCmd(), tickCmd())
}

// Add to Update():
case tickMsg:
	// Update git info
	if branch, err := git.GetBranch(); err == nil {
		m.state.Git.Branch = branch
	}

	if status, err := git.GetStatus(); err == nil {
		m.state.Git.DirtyFiles = status.DirtyFiles
		m.state.Git.Ahead = status.Ahead
		m.state.Git.Behind = status.Behind
		m.state.Git.Added = status.Added
		m.state.Git.Modified = status.Modified
		m.state.Git.Deleted = status.Deleted
	}

	return m, tea.Batch(outputCmd(m), tickCmd())
```

**Step 3: Build and test**

```bash
go build .
```

Expected: Successful build

**Step 4: Commit**

```bash
git add main.go
git commit -m "feat: add config loading and git refresh

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Phase 11: Integration Testing & Polish

### Task 11.1: End-to-End Integration Test

**Files:**
- Create: `integration_test.go`

**Step 1: Write integration test**

Create `integration_test.go`:
```go
//go:build integration
// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	// Build binary
	cmd := exec.Command("go", "build", "-o", "cc-hud-go-test", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}
	defer exec.Command("rm", "cc-hud-go-test").Run()

	// Run with mock input
	cmd = exec.Command("./cc-hud-go-test")

	input := `{"model":"claude-sonnet-4.5","planType":"Pro","context":{"used":50000,"total":200000}}`
	cmd.Stdin = bytes.NewBufferString(input + "\n")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// Start and wait briefly
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start: %v", err)
	}

	time.Sleep(500 * time.Millisecond)
	cmd.Process.Kill()

	// Parse output
	var result map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse output: %v\n%s", err, stdout.String())
	}

	segments, ok := result["segments"].([]interface{})
	if !ok || len(segments) == 0 {
		t.Error("expected segments in output")
	}
}
```

**Step 2: Run integration test**

```bash
go test -tags=integration -v .
```

Expected: PASS

**Step 3: Commit**

```bash
git add integration_test.go
git commit -m "test: add end-to-end integration test

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 11.2: README Documentation

**Files:**
- Create: `README.md`

**Step 1: Write comprehensive README**

Create `README.md`:
```markdown
# cc-hud-go

A Go-based statusline tool for Claude Code that displays comprehensive session information with enhanced tool tracking.

## Features

- **Model & Context**: Real-time model info and token usage with color-coded thresholds
- **Git Integration**: Branch status, dirty files, ahead/behind tracking
- **Enhanced Tool Tracking**: Categorizes app tools, MCP tools, skills, and custom tools
- **Task Progress**: Track pending, in-progress, and completed tasks
- **Agent Monitoring**: See active subagents and their tasks
- **Rate Limits**: Usage tracking for Pro/Max/Team plans
- **Configurable**: Three presets (Full/Essential/Minimal) plus granular controls

## Installation

```bash
go install github.com/yourusername/cc-hud-go@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/cc-hud-go.git
cd cc-hud-go
go build .
```

## Configuration

Config file location: `~/.claude/cc-hud-go/config.json`

See [docs/plans/2026-02-09-cc-hud-go-design.md](docs/plans/2026-02-09-cc-hud-go-design.md) for full configuration options.

### Presets

- **Full**: All segments enabled, maximum detail
- **Essential**: Core metrics only (model, context, git, tasks)
- **Minimal**: Just model and context bar

## Usage

cc-hud-go integrates with Claude Code's statusline API. Configure it in your Claude Code settings:

```json
{
  "statusline": {
    "command": "cc-hud-go"
  }
}
```

## Development

**Run tests:**
```bash
go test ./...
```

**Run integration tests:**
```bash
go test -tags=integration -v .
```

**Build:**
```bash
go build .
```

## Architecture

Built with:
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Event-driven TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

See [docs/plans/2026-02-09-cc-hud-go-design.md](docs/plans/2026-02-09-cc-hud-go-design.md) for detailed architecture documentation.

## License

MIT

## Acknowledgments

Inspired by:
- [claude-hud](https://github.com/jarrodwatts/claude-hud)
- [Oh My Posh Claude segment](https://ohmyposh.dev/docs/segments/cli/claude)
- The Charm ecosystem
```

**Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add comprehensive README

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Task 11.3: Final Testing & Bug Fixes

**Step 1: Run all tests**

```bash
go test ./... -cover
```

**Step 2: Run integration tests**

```bash
go test -tags=integration -v .
```

**Step 3: Manual testing**

```bash
go build .
echo '{"model":"test","planType":"Pro","context":{"used":100000,"total":200000}}' | ./cc-hud-go
```

**Step 4: Fix any issues found**

(Address bugs as they arise during testing)

**Step 5: Final commit**

```bash
git add .
git commit -m "chore: final testing and polish

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Completion

The implementation is complete! You now have:

✅ Feature-complete Go statusline tool
✅ TDD approach with >80% coverage goal
✅ Enhanced tool tracking (app/MCP/skills/custom)
✅ Bubbletea event-driven architecture
✅ Lipgloss styling with color thresholds
✅ Graceful error handling
✅ Comprehensive documentation

**Next steps:**
1. Deploy to Claude Code statusline
2. Test with real Claude Code sessions
3. Gather user feedback
4. Iterate on features

**Total estimated time:** 4-6 hours of focused TDD implementation

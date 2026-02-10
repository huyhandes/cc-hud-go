# Lipgloss Visual Enhancement Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Catppuccin theme support, gradient progress bars, enhanced spacing, and contextual table rendering to cc-hud-go statusline.

**Architecture:** Create theme package for color management, add gradient bar rendering to style package, extend config with table thresholds, update segments to use new rendering capabilities.

**Tech Stack:** Go 1.21+, lipgloss (TUI styling), Catppuccin color palette

---

## Phase 1: Theme System Foundation

### Task 1: Create Theme Package Structure

**Files:**
- Create: `theme/theme.go`
- Create: `theme/theme_test.go`

**Step 1: Write the theme interface test**

Create `theme/theme_test.go`:
```go
package theme

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestThemeInterface(t *testing.T) {
	theme := GetTheme("macchiato")

	// Test that theme returns valid colors
	color := theme.GetColor("success")
	if color == lipgloss.Color("") {
		t.Error("Expected non-empty color for 'success'")
	}
}

func TestInvalidTheme(t *testing.T) {
	theme := GetTheme("invalid")

	// Should fall back to macchiato
	if theme == nil {
		t.Error("Expected fallback theme, got nil")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./theme -v`
Expected: FAIL with "undefined: GetTheme"

**Step 3: Write minimal theme interface**

Create `theme/theme.go`:
```go
package theme

import "github.com/charmbracelet/lipgloss"

// Theme defines the color palette interface
type Theme interface {
	Name() string
	GetColor(semantic string) lipgloss.Color
}

// GetTheme returns a theme by name, falls back to macchiato
func GetTheme(name string) Theme {
	switch name {
	case "macchiato":
		return NewMacchiato()
	default:
		return NewMacchiato() // fallback
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./theme -v`
Expected: FAIL with "undefined: NewMacchiato"

**Step 5: Commit**

```bash
git add theme/
git commit -m "feat(theme): add theme package structure and interface"
```

---

### Task 2: Implement Catppuccin Macchiato Theme

**Files:**
- Create: `theme/catppuccin.go`
- Modify: `theme/theme_test.go`

**Step 1: Write Macchiato theme test**

Add to `theme/theme_test.go`:
```go
func TestMacchiatoColors(t *testing.T) {
	theme := NewMacchiato()

	tests := []struct {
		semantic string
		expected string
	}{
		{"success", "#a6da95"},
		{"warning", "#eed49f"},
		{"danger", "#ed8796"},
		{"input", "#8aadf4"},
		{"output", "#8bd5ca"},
	}

	for _, tt := range tests {
		t.Run(tt.semantic, func(t *testing.T) {
			color := theme.GetColor(tt.semantic)
			if string(color) != tt.expected {
				t.Errorf("GetColor(%s) = %s, want %s", tt.semantic, color, tt.expected)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./theme -v`
Expected: FAIL with "undefined: NewMacchiato"

**Step 3: Implement Macchiato theme**

Create `theme/catppuccin.go`:
```go
package theme

import "github.com/charmbracelet/lipgloss"

// Macchiato implements the Catppuccin Macchiato theme
type Macchiato struct{}

func NewMacchiato() *Macchiato {
	return &Macchiato{}
}

func (m *Macchiato) Name() string {
	return "macchiato"
}

func (m *Macchiato) GetColor(semantic string) lipgloss.Color {
	colors := map[string]string{
		// Status colors
		"success": "#a6da95", // Green
		"warning": "#eed49f", // Yellow
		"danger":  "#ed8796", // Red

		// Flow colors
		"input":  "#8aadf4", // Blue
		"output": "#8bd5ca", // Teal

		// Cache colors
		"cacheRead":  "#c6a0f6", // Mauve
		"cacheWrite": "#f5bde6", // Pink

		// Primary UI colors
		"primary":   "#c6a0f6", // Mauve
		"highlight": "#91d7e3", // Sky
		"accent":    "#f5a97f", // Peach

		// Utility colors
		"muted":  "#5b6078", // Overlay0
		"bright": "#cad3f5", // Text
		"info":   "#8bd5ca", // Teal
	}

	if color, ok := colors[semantic]; ok {
		return lipgloss.Color(color)
	}
	return lipgloss.Color("#cad3f5") // default to text color
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./theme -v`
Expected: PASS

**Step 5: Commit**

```bash
git add theme/catppuccin.go theme/theme_test.go
git commit -m "feat(theme): implement Catppuccin Macchiato theme"
```

---

### Task 3: Add All Catppuccin Flavors

**Files:**
- Modify: `theme/catppuccin.go`
- Modify: `theme/theme.go`
- Modify: `theme/theme_test.go`

**Step 1: Write test for all flavors**

Add to `theme/theme_test.go`:
```go
func TestAllCatppuccinFlavors(t *testing.T) {
	flavors := []string{"mocha", "macchiato", "frappe", "latte"}

	for _, flavor := range flavors {
		t.Run(flavor, func(t *testing.T) {
			theme := GetTheme(flavor)
			if theme == nil {
				t.Errorf("GetTheme(%s) returned nil", flavor)
			}
			if theme.Name() != flavor {
				t.Errorf("Expected theme name %s, got %s", flavor, theme.Name())
			}

			// Verify all semantic colors exist
			semantics := []string{"success", "warning", "danger", "input", "output"}
			for _, semantic := range semantics {
				color := theme.GetColor(semantic)
				if color == lipgloss.Color("") {
					t.Errorf("Theme %s missing color for %s", flavor, semantic)
				}
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./theme -v`
Expected: FAIL - only macchiato exists

**Step 3: Implement other Catppuccin flavors**

Update `theme/catppuccin.go`:
```go
// Mocha theme
type Mocha struct{}

func NewMocha() *Mocha {
	return &Mocha{}
}

func (m *Mocha) Name() string {
	return "mocha"
}

func (m *Mocha) GetColor(semantic string) lipgloss.Color {
	colors := map[string]string{
		"success":    "#a6e3a1", // Green
		"warning":    "#f9e2af", // Yellow
		"danger":     "#f38ba8", // Red
		"input":      "#89b4fa", // Blue
		"output":     "#94e2d5", // Teal
		"cacheRead":  "#cba6f7", // Mauve
		"cacheWrite": "#f5c2e7", // Pink
		"primary":    "#cba6f7", // Mauve
		"highlight":  "#89dceb", // Sky
		"accent":     "#fab387", // Peach
		"muted":      "#585b70", // Overlay0
		"bright":     "#cdd6f4", // Text
		"info":       "#94e2d5", // Teal
	}

	if color, ok := colors[semantic]; ok {
		return lipgloss.Color(color)
	}
	return lipgloss.Color("#cdd6f4")
}

// Frappe theme
type Frappe struct{}

func NewFrappe() *Frappe {
	return &Frappe{}
}

func (f *Frappe) Name() string {
	return "frappe"
}

func (f *Frappe) GetColor(semantic string) lipgloss.Color {
	colors := map[string]string{
		"success":    "#a6d189", // Green
		"warning":    "#e5c890", // Yellow
		"danger":     "#e78284", // Red
		"input":      "#8caaee", // Blue
		"output":     "#81c8be", // Teal
		"cacheRead":  "#ca9ee6", // Mauve
		"cacheWrite": "#f4b8e4", // Pink
		"primary":    "#ca9ee6", // Mauve
		"highlight":  "#99d1db", // Sky
		"accent":     "#ef9f76", // Peach
		"muted":      "#51576d", // Overlay0
		"bright":     "#c6d0f5", // Text
		"info":       "#81c8be", // Teal
	}

	if color, ok := colors[semantic]; ok {
		return lipgloss.Color(color)
	}
	return lipgloss.Color("#c6d0f5")
}

// Latte theme
type Latte struct{}

func NewLatte() *Latte {
	return &Latte{}
}

func (l *Latte) Name() string {
	return "latte"
}

func (l *Latte) GetColor(semantic string) lipgloss.Color {
	colors := map[string]string{
		"success":    "#40a02b", // Green
		"warning":    "#df8e1d", // Yellow
		"danger":     "#d20f39", // Red
		"input":      "#1e66f5", // Blue
		"output":     "#179299", // Teal
		"cacheRead":  "#8839ef", // Mauve
		"cacheWrite": "#ea76cb", // Pink
		"primary":    "#8839ef", // Mauve
		"highlight":  "#04a5e5", // Sky
		"accent":     "#fe640b", // Peach
		"muted":      "#9ca0b0", // Overlay0
		"bright":     "#4c4f69", // Text
		"info":       "#179299", // Teal
	}

	if color, ok := colors[semantic]; ok {
		return lipgloss.Color(color)
	}
	return lipgloss.Color("#4c4f69")
}
```

Update `theme/theme.go`:
```go
func GetTheme(name string) Theme {
	switch name {
	case "mocha":
		return NewMocha()
	case "macchiato":
		return NewMacchiato()
	case "frappe":
		return NewFrappe()
	case "latte":
		return NewLatte()
	default:
		return NewMacchiato() // fallback to macchiato
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./theme -v`
Expected: PASS

**Step 5: Commit**

```bash
git add theme/
git commit -m "feat(theme): add all Catppuccin flavors (mocha, frappe, latte)"
```

---

### Task 4: Add Theme Config Support

**Files:**
- Modify: `config/config.go`
- Modify: `config/config_test.go`

**Step 1: Write config test for theme**

Add to `config/config_test.go`:
```go
func TestConfigTheme(t *testing.T) {
	cfg := Default()

	if cfg.Theme == "" {
		t.Error("Expected default theme to be set")
	}

	if cfg.Theme != "macchiato" {
		t.Errorf("Expected default theme 'macchiato', got %s", cfg.Theme)
	}
}

func TestConfigColorOverrides(t *testing.T) {
	cfg := Default()

	if cfg.Colors == nil {
		t.Error("Expected Colors map to be initialized")
	}
}

func TestTableConfigDefaults(t *testing.T) {
	cfg := Default()

	if cfg.Tables.ToolsThreshold != 5 {
		t.Errorf("Expected ToolsThreshold 5, got %d", cfg.Tables.ToolsThreshold)
	}

	if cfg.Tables.TasksThreshold != 3 {
		t.Errorf("Expected TasksThreshold 3, got %d", cfg.Tables.TasksThreshold)
	}

	if cfg.Tables.ContextThreshold != 999 {
		t.Errorf("Expected ContextThreshold 999, got %d", cfg.Tables.ContextThreshold)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./config -v`
Expected: FAIL - missing Theme, Colors, Tables fields

**Step 3: Add theme config fields**

Update `config/config.go`:
```go
type Config struct {
	Theme             string
	Colors            map[string]string
	Preset            string
	LineLayout        string
	PathLevels        int
	ContextValue      string
	SevenDayThreshold int
	Display           DisplayConfig
	Git               GitConfig
	Tools             ToolsConfig
	Tables            TableConfig
}

type TableConfig struct {
	ToolsThreshold   int `json:"toolsTableThreshold"`
	TasksThreshold   int `json:"tasksTableThreshold"`
	ContextThreshold int `json:"contextTableThreshold"`
}

func Default() *Config {
	return &Config{
		Theme:             "macchiato",
		Colors:            make(map[string]string),
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
		Tables: TableConfig{
			ToolsThreshold:   5,
			TasksThreshold:   3,
			ContextThreshold: 999,
		},
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./config -v`
Expected: PASS

**Step 5: Commit**

```bash
git add config/
git commit -m "feat(config): add theme and table threshold configuration"
```

---

### Task 5: Add Theme Loader with Color Overrides

**Files:**
- Modify: `theme/theme.go`
- Modify: `theme/theme_test.go`

**Step 1: Write test for LoadFromConfig**

Add to `theme/theme_test.go`:
```go
func TestLoadFromConfig(t *testing.T) {
	// Mock config structure
	type Config struct {
		Theme  string
		Colors map[string]string
	}

	cfg := Config{
		Theme: "mocha",
		Colors: map[string]string{
			"success": "#00ff00",
		},
	}

	// Create a theme wrapper that can apply overrides
	wrapper := LoadThemeFromConfig(cfg.Theme, cfg.Colors)

	// Check base theme
	if wrapper.Name() != "mocha" {
		t.Errorf("Expected theme 'mocha', got %s", wrapper.Name())
	}

	// Check override applied
	successColor := wrapper.GetColor("success")
	if string(successColor) != "#00ff00" {
		t.Errorf("Expected override color #00ff00, got %s", successColor)
	}

	// Check non-overridden color uses theme default
	warningColor := wrapper.GetColor("warning")
	if warningColor == lipgloss.Color("") {
		t.Error("Expected warning color from theme")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./theme -v`
Expected: FAIL with "undefined: LoadThemeFromConfig"

**Step 3: Implement theme loader with overrides**

Add to `theme/theme.go`:
```go
// ThemeWrapper wraps a theme and applies color overrides
type ThemeWrapper struct {
	base      Theme
	overrides map[string]string
}

func (tw *ThemeWrapper) Name() string {
	return tw.base.Name()
}

func (tw *ThemeWrapper) GetColor(semantic string) lipgloss.Color {
	// Check for override first
	if color, ok := tw.overrides[semantic]; ok {
		return lipgloss.Color(color)
	}
	// Fall back to base theme
	return tw.base.GetColor(semantic)
}

// LoadThemeFromConfig loads a theme and applies color overrides
func LoadThemeFromConfig(themeName string, colorOverrides map[string]string) Theme {
	base := GetTheme(themeName)

	if len(colorOverrides) == 0 {
		return base
	}

	return &ThemeWrapper{
		base:      base,
		overrides: colorOverrides,
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./theme -v`
Expected: PASS

**Step 5: Commit**

```bash
git add theme/
git commit -m "feat(theme): add theme loader with color override support"
```

---

## Phase 2: Gradient Progress Bars

### Task 6: Add Gradient Bar Helper to Style Package

**Files:**
- Modify: `style/style.go`
- Create: `style/gradient_test.go`

**Step 1: Write gradient bar test**

Create `style/gradient_test.go`:
```go
package style

import (
	"strings"
	"testing"
)

func TestRenderGradientBar(t *testing.T) {
	// Initialize with a mock theme for testing
	// We'll update Init() to accept theme later

	tests := []struct {
		name       string
		percentage float64
		width      int
		wantFilled int
	}{
		{"0 percent", 0, 10, 0},
		{"50 percent", 50, 10, 5},
		{"100 percent", 100, 10, 10},
		{"75 percent", 75, 10, 7}, // floor(7.5) = 7
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderGradientBar(tt.percentage, tt.width)

			// Strip ANSI codes for testing
			stripped := stripAnsi(result)

			// Count filled characters (█▓▒)
			filled := strings.Count(stripped, "█") +
			          strings.Count(stripped, "▓") +
			          strings.Count(stripped, "▒")
			empty := strings.Count(stripped, "░")

			if filled < tt.wantFilled-1 || filled > tt.wantFilled+1 {
				t.Errorf("Expected ~%d filled chars, got %d (result: %s)",
					tt.wantFilled, filled, stripped)
			}

			if filled+empty != tt.width {
				t.Errorf("Expected total width %d, got %d", tt.width, filled+empty)
			}
		})
	}
}

// Helper to strip ANSI codes for testing
func stripAnsi(s string) string {
	// Simple strip for testing - remove escape sequences
	result := ""
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
		} else if inEscape && r == 'm' {
			inEscape = false
		} else if !inEscape {
			result += string(r)
		}
	}
	return result
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./style -v`
Expected: FAIL with "undefined: RenderGradientBar"

**Step 3: Implement gradient bar rendering**

Add to `style/style.go`:
```go
// RenderGradientBar renders a gradient progress bar
func RenderGradientBar(percentage float64, width int) string {
	if width <= 0 {
		width = 10
	}
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}

	filled := int(percentage / 100 * float64(width))
	if filled > width {
		filled = width
	}

	segments := make([]string, 0, width)

	for i := 0; i < width; i++ {
		if i < filled {
			// Use gradient characters for filled portion
			char := getGradientChar(i, filled, width)
			color := getColorForPercentage(percentage)
			segments = append(segments, renderer.NewStyle().Foreground(color).Render(char))
		} else {
			// Empty portion
			segments = append(segments, renderer.NewStyle().Foreground(ColorMuted).Render("░"))
		}
	}

	return strings.Join(segments, "")
}

// getGradientChar returns the appropriate gradient character
func getGradientChar(position, filled, width int) string {
	if filled == 0 {
		return "░"
	}

	// Use different characters based on position in filled area
	progress := float64(position) / float64(filled)

	if progress < 0.3 {
		return "█" // solid
	} else if progress < 0.6 {
		return "▓" // dark
	} else {
		return "▒" // medium
	}
}

// getColorForPercentage returns color based on percentage thresholds
func getColorForPercentage(percentage float64) lipgloss.Color {
	if percentage >= 90 {
		return ColorDanger // red
	} else if percentage >= 70 {
		return ColorWarning // yellow/orange
	}
	return ColorSuccess // green
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./style -v`
Expected: PASS

**Step 5: Commit**

```bash
git add style/
git commit -m "feat(style): add gradient progress bar rendering"
```

---

### Task 7: Update Style Package to Use Theme

**Files:**
- Modify: `style/style.go`
- Modify: `style/gradient_test.go`

**Step 1: Write test for style initialization with theme**

Add to `style/gradient_test.go`:
```go
func TestInitWithTheme(t *testing.T) {
	// Create a mock theme
	type mockTheme struct{}

	// We'll define this interface in style.go
	var theme Theme = &mockTheme{}

	Init(theme)

	// Verify colors are set from theme
	if ColorSuccess == lipgloss.Color("") {
		t.Error("Expected ColorSuccess to be initialized from theme")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./style -v`
Expected: FAIL with "undefined: Theme" or "undefined: Init"

**Step 3: Update style package to use theme**

Update `style/style.go`:
```go
package style

import (
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Theme interface for color access
type Theme interface {
	Name() string
	GetColor(semantic string) lipgloss.Color
}

var (
	// Global renderer that forces color output
	renderer *lipgloss.Renderer

	// Current theme
	currentTheme Theme

	// Color palette - loaded from theme
	ColorSuccess    lipgloss.Color
	ColorWarning    lipgloss.Color
	ColorDanger     lipgloss.Color
	ColorInput      lipgloss.Color
	ColorOutput     lipgloss.Color
	ColorCacheRead  lipgloss.Color
	ColorCacheWrite lipgloss.Color
	ColorPrimary    lipgloss.Color
	ColorHighlight  lipgloss.Color
	ColorAccent     lipgloss.Color
	ColorMuted      lipgloss.Color
	ColorBright     lipgloss.Color
	ColorInfo       lipgloss.Color

	// Pre-configured styles
	ModelStyle     lipgloss.Style
	ContextStyle   lipgloss.Style
	GitStyle       lipgloss.Style
	CostStyle      lipgloss.Style
	ToolsStyle     lipgloss.Style
	AgentStyle     lipgloss.Style
	SeparatorStyle lipgloss.Style

	// Progress bar styles
	ProgressGood    lipgloss.Style
	ProgressWarning lipgloss.Style
	ProgressDanger  lipgloss.Style
)

func init() {
	// Create renderer that forces color output with TrueColor profile
	renderer = lipgloss.NewRenderer(os.Stdout, termenv.WithProfile(termenv.TrueColor))
	renderer.SetColorProfile(termenv.TrueColor)
}

// Init initializes styles with the given theme
func Init(theme Theme) {
	currentTheme = theme

	// Load colors from theme
	ColorSuccess = theme.GetColor("success")
	ColorWarning = theme.GetColor("warning")
	ColorDanger = theme.GetColor("danger")
	ColorInput = theme.GetColor("input")
	ColorOutput = theme.GetColor("output")
	ColorCacheRead = theme.GetColor("cacheRead")
	ColorCacheWrite = theme.GetColor("cacheWrite")
	ColorPrimary = theme.GetColor("primary")
	ColorHighlight = theme.GetColor("highlight")
	ColorAccent = theme.GetColor("accent")
	ColorMuted = theme.GetColor("muted")
	ColorBright = theme.GetColor("bright")
	ColorInfo = theme.GetColor("info")

	// Initialize styles with theme colors
	ModelStyle = renderer.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)

	ContextStyle = renderer.NewStyle().
		Foreground(ColorInfo)

	GitStyle = renderer.NewStyle().
		Foreground(ColorHighlight)

	CostStyle = renderer.NewStyle().
		Foreground(ColorAccent)

	ToolsStyle = renderer.NewStyle().
		Foreground(ColorSuccess)

	AgentStyle = renderer.NewStyle().
		Foreground(ColorPrimary).
		Italic(true)

	SeparatorStyle = renderer.NewStyle().
		Foreground(ColorMuted)

	// Progress bar color schemes
	ProgressGood = renderer.NewStyle().
		Foreground(ColorSuccess)

	ProgressWarning = renderer.NewStyle().
		Foreground(ColorWarning)

	ProgressDanger = renderer.NewStyle().
		Foreground(ColorDanger)
}

// Rest of existing functions (GetRenderer, NewRendererForWriter, Separator, Icon)
// ... keep existing code ...
```

**Step 4: Run test to verify it passes**

Run: `go test ./style -v`
Expected: FAIL - need to provide mock theme implementation

Update test:
```go
func TestInitWithTheme(t *testing.T) {
	// Mock theme implementation
	type mockTheme struct{}

	func (m *mockTheme) Name() string { return "test" }
	func (m *mockTheme) GetColor(semantic string) lipgloss.Color {
		return lipgloss.Color("#ff0000")
	}

	theme := &mockTheme{}
	Init(theme)

	// Verify colors are set from theme
	if ColorSuccess == lipgloss.Color("") {
		t.Error("Expected ColorSuccess to be initialized from theme")
	}

	if string(ColorSuccess) != "#ff0000" {
		t.Errorf("Expected ColorSuccess #ff0000, got %s", ColorSuccess)
	}
}
```

Run: `go test ./style -v`
Expected: PASS

**Step 5: Commit**

```bash
git add style/
git commit -m "feat(style): integrate theme system for color management"
```

---

### Task 8: Update Context Segment to Use Gradient Bar

**Files:**
- Modify: `segment/context.go`
- Modify: `segment/context_test.go`

**Step 1: Write test for gradient bar in context**

Add to `segment/context_test.go` (or create if doesn't exist):
```go
func TestContextSegmentUsesGradientBar(t *testing.T) {
	// Setup
	s := &state.State{
		Context: state.ContextState{
			UsedTokens:  54000,
			TotalTokens: 100000,
			Percentage:  54,
		},
	}
	cfg := config.Default()

	seg := &ContextSegment{}
	result, err := seg.Render(s, cfg)

	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should contain gradient bar characters, not dots
	if strings.Contains(result, "●") || strings.Contains(result, "○") {
		t.Error("Expected gradient bar, found old dot characters")
	}

	// Should contain gradient characters
	hasGradient := strings.Contains(result, "█") ||
	               strings.Contains(result, "▓") ||
	               strings.Contains(result, "▒") ||
	               strings.Contains(result, "░")

	if !hasGradient {
		t.Error("Expected gradient bar characters (█▓▒░)")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./segment -run TestContextSegmentUsesGradientBar -v`
Expected: FAIL - still uses dots

**Step 3: Update context segment implementation**

Update `segment/context.go`:
```go
func (c *ContextSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	if s.Context.TotalTokens == 0 {
		return "", nil
	}

	percentage := s.Context.Percentage

	// Choose icon based on thresholds
	var icon = "🟢"
	if percentage >= 90 {
		icon = "🔴"
	} else if percentage >= 70 {
		icon = "🟡"
	}

	// Build gradient progress bar
	bar := style.RenderGradientBar(percentage, 10)

	// Format tokens in thousands (k)
	formatTokens := func(tokens int) string {
		if tokens >= 1000 {
			return fmt.Sprintf("%dk", tokens/1000)
		}
		return fmt.Sprintf("%d", tokens)
	}

	// Main display with bar and percentage
	percentageColor := style.ColorSuccess
	if percentage >= 90 {
		percentageColor = style.ColorDanger
	} else if percentage >= 70 {
		percentageColor = style.ColorWarning
	}

	percentageStyle := style.GetRenderer().NewStyle().Foreground(percentageColor)
	mainDisplay := fmt.Sprintf("%s %s %s",
		icon,
		bar,
		percentageStyle.Render(fmt.Sprintf("%.0f%%", percentage)),
	)

	// Detailed token breakdown with semantic colors (keep existing)
	details := []string{}

	inStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInput)
	details = append(details,
		fmt.Sprintf("📥 %s", inStyle.Render(formatTokens(s.Context.TotalInputTokens))),
	)

	outStyle := style.GetRenderer().NewStyle().Foreground(style.ColorOutput)
	details = append(details,
		fmt.Sprintf("📤 %s", outStyle.Render(formatTokens(s.Context.TotalOutputTokens))),
	)

	if s.Context.CacheReadTokens > 0 || s.Context.CacheCreateTokens > 0 {
		cacheReadStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCacheRead)
		cacheWriteStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCacheWrite)

		details = append(details,
			fmt.Sprintf("💾 %s%s%s",
				cacheReadStyle.Render("R:"+formatTokens(s.Context.CacheReadTokens)),
				style.GetRenderer().NewStyle().Foreground(style.ColorMuted).Render("/"),
				cacheWriteStyle.Render("W:"+formatTokens(s.Context.CacheCreateTokens)),
			),
		)
	}

	totalStyle := style.GetRenderer().NewStyle().Foreground(style.ColorMuted)
	details = append(details,
		fmt.Sprintf("⚡ %s", totalStyle.Render(formatTokens(s.Context.TotalTokens))),
	)

	return fmt.Sprintf("%s %s", mainDisplay, strings.Join(details, " ")), nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./segment -run TestContextSegmentUsesGradientBar -v`
Expected: PASS

**Step 5: Commit**

```bash
git add segment/context.go segment/context_test.go
git commit -m "feat(segment): use gradient bar in context segment"
```

---

### Task 9: Update Rate Limit Segment for Gradient Bar

**Files:**
- Modify: `segment/ratelimit.go`
- Create/Modify: `segment/ratelimit_test.go`

**Step 1: Write test for rate limit gradient bar**

Create or update `segment/ratelimit_test.go`:
```go
func TestRateLimitUsesGradientBar(t *testing.T) {
	s := &state.State{
		RateLimits: state.RateLimitsState{
			SevenDayUsage:    67,
			SevenDayRequests: 670,
			SevenDayLimit:    1000,
		},
	}
	cfg := config.Default()

	seg := &RateLimitSegment{}
	result, err := seg.Render(s, cfg)

	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should contain gradient bar characters
	hasGradient := strings.Contains(result, "█") ||
	               strings.Contains(result, "▓") ||
	               strings.Contains(result, "▒") ||
	               strings.Contains(result, "░")

	if !hasGradient {
		t.Error("Expected gradient bar characters in rate limit segment")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./segment -run TestRateLimitUsesGradientBar -v`
Expected: FAIL - rate limit doesn't use gradient yet

**Step 3: Update rate limit segment**

Update `segment/ratelimit.go` to use gradient bar (find the rendering section):
```go
// Replace any existing progress bar with:
bar := style.RenderGradientBar(sevenDayPercentage, 10)

// Update the output format to include the gradient bar
output := fmt.Sprintf("📊 %s %.0f%% (%d/%d)",
	bar,
	sevenDayPercentage,
	s.RateLimits.SevenDayRequests,
	s.RateLimits.SevenDayLimit,
)
```

**Step 4: Run test to verify it passes**

Run: `go test ./segment -run TestRateLimitUsesGradientBar -v`
Expected: PASS

**Step 5: Commit**

```bash
git add segment/ratelimit.go segment/ratelimit_test.go
git commit -m "feat(segment): use gradient bar in rate limit segment"
```

---

## Phase 3: Layout & Spacing

### Task 10: Update Renderer for Enhanced Spacing

**Files:**
- Modify: `output/renderer.go`
- Modify: `output/renderer_test.go`

**Step 1: Write test for two-space separators**

Add to `output/renderer_test.go`:
```go
func TestJoinSegmentsWithSpacing(t *testing.T) {
	segments := []string{"segment1", "segment2", "segment3"}

	result := joinSegments(segments)

	// Should have 2 spaces around separator
	expected := "segment1  │  segment2  │  segment3"

	if !strings.Contains(result, "  │  ") {
		t.Errorf("Expected two-space separator, got: %s", result)
	}

	// Count separators
	sepCount := strings.Count(result, "│")
	if sepCount != 2 {
		t.Errorf("Expected 2 separators, got %d", sepCount)
	}
}

func TestEmptySegmentsFiltering(t *testing.T) {
	segments := []string{"segment1", "", "segment2"}

	result := joinSegments(segments)

	// Should filter out empty segments
	if strings.Contains(result, "│  │") {
		t.Error("Empty segments should be filtered out")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./output -v`
Expected: FAIL with "undefined: joinSegments"

**Step 3: Implement join segments helper**

Add to `output/renderer.go`:
```go
// joinSegments joins segment outputs with two-space separators
func joinSegments(segments []string) string {
	// Filter out empty segments
	nonEmpty := make([]string, 0, len(segments))
	for _, seg := range segments {
		if strings.TrimSpace(seg) != "" {
			nonEmpty = append(nonEmpty, seg)
		}
	}

	return strings.Join(nonEmpty, "  │  ")
}
```

Update the main rendering function to use `joinSegments`:
```go
// Find where segments are currently joined and replace with:
text := joinSegments(segmentOutputs)
```

**Step 4: Run test to verify it passes**

Run: `go test ./output -v`
Expected: PASS

**Step 5: Commit**

```bash
git add output/
git commit -m "feat(output): add two-space separator spacing"
```

---

## Phase 4: Table Rendering

### Task 11: Add Table Helper to Style Package

**Files:**
- Modify: `style/style.go`
- Create: `style/table_test.go`

**Step 1: Write table rendering test**

Create `style/table_test.go`:
```go
package style

import (
	"strings"
	"testing"
)

func TestRenderTable(t *testing.T) {
	headers := []string{"Column1", "Column2"}
	rows := [][]string{
		{"Value1", "Value2"},
		{"Value3", "Value4"},
	}

	result := RenderTable(headers, rows)

	// Should contain table border characters
	if !strings.Contains(result, "┌") || !strings.Contains(result, "┐") {
		t.Error("Expected table top border characters")
	}

	if !strings.Contains(result, "└") || !strings.Contains(result, "┘") {
		t.Error("Expected table bottom border characters")
	}

	// Should contain headers
	if !strings.Contains(result, "Column1") || !strings.Contains(result, "Column2") {
		t.Error("Expected table headers in output")
	}

	// Should contain data
	if !strings.Contains(result, "Value1") || !strings.Contains(result, "Value4") {
		t.Error("Expected table data in output")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./style -run TestRenderTable -v`
Expected: FAIL with "undefined: RenderTable"

**Step 3: Implement table rendering helper**

Add to `style/style.go`:
```go
// RenderTable renders a table with headers and rows using lipgloss
func RenderTable(headers []string, rows [][]string) string {
	// Use lipgloss table (if available) or simple manual rendering
	// For now, let's use a simple implementation

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var result strings.Builder

	// Top border
	result.WriteString("┌")
	for i, w := range widths {
		result.WriteString(strings.Repeat("─", w+2))
		if i < len(widths)-1 {
			result.WriteString("┬")
		}
	}
	result.WriteString("┐\n")

	// Headers
	result.WriteString("│")
	for i, h := range headers {
		result.WriteString(" ")
		result.WriteString(h)
		result.WriteString(strings.Repeat(" ", widths[i]-len(h)+1))
		result.WriteString("│")
	}
	result.WriteString("\n")

	// Header separator
	result.WriteString("├")
	for i, w := range widths {
		result.WriteString(strings.Repeat("─", w+2))
		if i < len(widths)-1 {
			result.WriteString("┼")
		}
	}
	result.WriteString("┤\n")

	// Rows
	for _, row := range rows {
		result.WriteString("│")
		for i, cell := range row {
			if i >= len(widths) {
				break
			}
			result.WriteString(" ")
			result.WriteString(cell)
			result.WriteString(strings.Repeat(" ", widths[i]-len(cell)+1))
			result.WriteString("│")
		}
		result.WriteString("\n")
	}

	// Bottom border
	result.WriteString("└")
	for i, w := range widths {
		result.WriteString(strings.Repeat("─", w+2))
		if i < len(widths)-1 {
			result.WriteString("┴")
		}
	}
	result.WriteString("┘")

	return result.String()
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./style -run TestRenderTable -v`
Expected: PASS

**Step 5: Commit**

```bash
git add style/
git commit -m "feat(style): add table rendering helper"
```

---

### Task 12: Add Table Support to Tools Segment

**Files:**
- Modify: `segment/tools.go`
- Modify: `segment/tools_test.go`

**Step 1: Write test for tools table threshold**

Add to `segment/tools_test.go`:
```go
func TestToolsSegmentTableThreshold(t *testing.T) {
	// Below threshold - should be inline
	s := &state.State{
		Tools: state.ToolsState{
			ByCategory: map[string]int{
				"App":   3,
				"MCP":   1,
				"Skills": 1,
			},
		},
	}
	cfg := config.Default()
	cfg.Tables.ToolsThreshold = 5

	seg := &ToolsSegment{}
	result, err := seg.Render(s, cfg)

	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should be inline (no table borders)
	if strings.Contains(result, "┌") {
		t.Error("Expected inline format below threshold")
	}

	// Above threshold - should be table
	s.Tools.ByCategory = map[string]int{
		"App":     10,
		"MCP":     3,
		"Skills":  2,
		"Custom":  1,
	}

	result, err = seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should be table format
	if !strings.Contains(result, "┌") {
		t.Error("Expected table format above threshold")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./segment -run TestToolsSegmentTableThreshold -v`
Expected: FAIL - threshold logic not implemented

**Step 3: Implement table logic in tools segment**

Update `segment/tools.go`:
```go
func (t *ToolsSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	if len(s.Tools.ByCategory) == 0 {
		return "", nil
	}

	toolCount := 0
	for _, count := range s.Tools.ByCategory {
		toolCount += count
	}

	// Check if we should use table view
	if toolCount > cfg.Tables.ToolsThreshold {
		return t.renderTable(s, cfg)
	}

	// Inline view (existing implementation)
	return t.renderInline(s, cfg)
}

func (t *ToolsSegment) renderInline(s *state.State, cfg *config.Config) (string, error) {
	// Existing inline implementation
	toolCount := 0
	for _, count := range s.Tools.ByCategory {
		toolCount += count
	}

	// Build category breakdown
	parts := []string{}
	for category, count := range s.Tools.ByCategory {
		parts = append(parts, fmt.Sprintf("%s:%d", category, count))
	}

	return fmt.Sprintf("🔧 %d (%s)", toolCount, strings.Join(parts, " ")), nil
}

func (t *ToolsSegment) renderTable(s *state.State, cfg *config.Config) (string, error) {
	headers := []string{"Category", "Count"}
	rows := [][]string{}

	// Sort categories for consistent order
	categories := make([]string, 0, len(s.Tools.ByCategory))
	for cat := range s.Tools.ByCategory {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	for _, cat := range categories {
		count := s.Tools.ByCategory[cat]
		rows = append(rows, []string{cat, fmt.Sprintf("%d", count)})
	}

	return style.RenderTable(headers, rows), nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./segment -run TestToolsSegmentTableThreshold -v`
Expected: PASS

**Step 5: Commit**

```bash
git add segment/tools.go segment/tools_test.go
git commit -m "feat(segment): add table view for tools segment with threshold"
```

---

### Task 13: Add Table Support to Tasks Segment

**Files:**
- Modify: `segment/tasks.go`
- Create/Modify: `segment/tasks_test.go`

**Step 1: Write test for tasks table with last 3 completed**

Create/update `segment/tasks_test.go`:
```go
func TestTasksTableFiltering(t *testing.T) {
	s := &state.State{
		Tasks: state.TasksState{
			Completed: []state.Task{
				{Subject: "Task 1", Status: "completed"},
				{Subject: "Task 2", Status: "completed"},
				{Subject: "Task 3", Status: "completed"},
				{Subject: "Task 4", Status: "completed"}, // should be filtered
				{Subject: "Task 5", Status: "completed"}, // should be filtered
			},
			Active: []state.Task{
				{Subject: "Active Task", Status: "in_progress"},
			},
			Pending: []state.Task{
				{Subject: "Pending Task", Status: "pending"},
			},
		},
	}
	cfg := config.Default()
	cfg.Tables.TasksThreshold = 3

	seg := &TasksSegment{}
	result, err := seg.Render(s, cfg)

	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should show last 3 completed
	if !strings.Contains(result, "Task 3") {
		t.Error("Expected Task 3 in output")
	}

	// Should NOT show Task 5 (too old)
	if strings.Contains(result, "Task 5") {
		t.Error("Task 5 should be filtered out (not in last 3)")
	}

	// Should show active and pending
	if !strings.Contains(result, "Active Task") {
		t.Error("Expected active task in output")
	}

	if !strings.Contains(result, "Pending Task") {
		t.Error("Expected pending task in output")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./segment -run TestTasksTableFiltering -v`
Expected: FAIL - filtering not implemented

**Step 3: Implement tasks table with filtering**

Update `segment/tasks.go`:
```go
func (t *TasksSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	totalTasks := len(s.Tasks.Completed) + len(s.Tasks.Active) + len(s.Tasks.Pending)

	if totalTasks == 0 {
		return "", nil
	}

	// Check if we should use table view
	if totalTasks > cfg.Tables.TasksThreshold {
		return t.renderTable(s, cfg)
	}

	// Inline view (existing)
	return t.renderInline(s, cfg)
}

func (t *TasksSegment) renderInline(s *state.State, cfg *config.Config) (string, error) {
	// Existing inline implementation
	completed := len(s.Tasks.Completed)
	active := len(s.Tasks.Active)
	pending := len(s.Tasks.Pending)
	total := completed + active + pending

	return fmt.Sprintf("✓ %d/%d tasks", completed, total), nil
}

func (t *TasksSegment) renderTable(s *state.State, cfg *config.Config) (string, error) {
	headers := []string{"Task", "Status"}
	rows := [][]string{}

	// Get last 3 completed tasks
	completed := s.Tasks.Completed
	startIdx := 0
	if len(completed) > 3 {
		startIdx = len(completed) - 3
	}
	recentCompleted := completed[startIdx:]

	// Add recent completed tasks
	for _, task := range recentCompleted {
		status := "✓ Done"
		rows = append(rows, []string{truncate(task.Subject, 40), status})
	}

	// Add active tasks
	for _, task := range s.Tasks.Active {
		status := "⏵ Active"
		rows = append(rows, []string{truncate(task.Subject, 40), status})
	}

	// Add pending tasks
	for _, task := range s.Tasks.Pending {
		status := "○ Pending"
		rows = append(rows, []string{truncate(task.Subject, 40), status})
	}

	return style.RenderTable(headers, rows), nil
}

// truncate truncates a string to maxLen with ellipsis
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./segment -run TestTasksTableFiltering -v`
Expected: PASS

**Step 5: Commit**

```bash
git add segment/tasks.go segment/tasks_test.go
git commit -m "feat(segment): add table view for tasks with last 3 completed filter"
```

---

## Phase 5: Integration & Main

### Task 14: Wire Theme into Main Application

**Files:**
- Modify: `main.go`

**Step 1: Write integration test**

Add to `main_test.go` or `integration_test.go`:
```go
func TestMainInitializesTheme(t *testing.T) {
	// This test verifies that main properly initializes the theme system
	// We'll test this by checking that style colors are set after initialization

	// Setup test config
	cfg := config.Default()
	cfg.Theme = "macchiato"

	// Initialize theme (extract this logic from main)
	themeInstance := theme.LoadThemeFromConfig(cfg.Theme, cfg.Colors)
	style.Init(themeInstance)

	// Verify colors are set
	if style.ColorSuccess == lipgloss.Color("") {
		t.Error("Expected ColorSuccess to be initialized")
	}

	// Verify it's using Catppuccin Macchiato colors
	if string(style.ColorSuccess) != "#a6da95" {
		t.Errorf("Expected Macchiato green, got %s", style.ColorSuccess)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test -run TestMainInitializesTheme -v`
Expected: FAIL - main doesn't initialize theme yet

**Step 3: Update main.go to initialize theme**

Update `main.go`:
```go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/output"
	"github.com/huybui/cc-hud-go/parser"
	"github.com/huybui/cc-hud-go/segment"
	"github.com/huybui/cc-hud-go/state"
	"github.com/huybui/cc-hud-go/style"
	"github.com/huybui/cc-hud-go/theme"
	"github.com/huybui/cc-hud-go/version"
)

func main() {
	// Parse flags
	showVersion := flag.Bool("version", false, "show version")
	flag.BoolVar(showVersion, "v", false, "show version (shorthand)")
	showHelp := flag.Bool("help", false, "show help")
	flag.BoolVar(showHelp, "h", false, "show help (shorthand)")
	flag.Parse()

	if *showVersion {
		fmt.Println(version.Version)
		return
	}

	if *showHelp {
		showHelpMessage()
		return
	}

	// Load config
	configPath := os.Getenv("CC_HUD_CONFIG")
	if configPath == "" {
		configPath = os.ExpandEnv("$HOME/.config/cc-hud-go/config.json")
	}

	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v, using defaults\n", err)
		cfg = config.Default()
	}

	// Initialize theme and style system
	themeInstance := theme.LoadThemeFromConfig(cfg.Theme, cfg.Colors)
	style.Init(themeInstance)

	// Read stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
		os.Exit(1)
	}

	// Parse state
	s, err := parser.ParseStdin(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing input: %v\n", err)
		os.Exit(1)
	}

	// Render segments
	renderer := output.NewRenderer(cfg)
	result, err := renderer.Render(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error rendering: %v\n", err)
		os.Exit(1)
	}

	// Output JSON
	output, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling output: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func showHelpMessage() {
	fmt.Println("cc-hud-go - Claude Code statusline HUD")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  cc-hud-go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help     Show this help message")
	fmt.Println("  -v, --version  Show version information")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  Set CC_HUD_CONFIG environment variable to specify config file")
	fmt.Println("  Default: $HOME/.config/cc-hud-go/config.json")
}
```

**Step 4: Run test to verify it passes**

Run: `go test -run TestMainInitializesTheme -v`
Expected: PASS

**Step 5: Commit**

```bash
git add main.go main_test.go
git commit -m "feat(main): initialize theme system at startup"
```

---

## Phase 6: Documentation & Examples

### Task 15: Create Example Config File

**Files:**
- Create: `examples/config-macchiato.json`
- Create: `examples/config-mocha.json`
- Create: `examples/config-custom-colors.json`

**Step 1: Create example configs**

Create `examples/config-macchiato.json`:
```json
{
  "theme": "macchiato",
  "preset": "full",
  "lineLayout": "expanded",
  "pathLevels": 2,
  "contextValue": "percentage",
  "sevenDayThreshold": 80,
  "display": {
    "model": true,
    "path": true,
    "context": true,
    "git": true,
    "tools": true,
    "agents": true,
    "tasks": true,
    "rateLimits": true,
    "duration": true,
    "speed": true
  },
  "git": {
    "showBranch": true,
    "showDirty": true,
    "showAheadBehind": true,
    "showFileStats": true
  },
  "tools": {
    "groupByCategory": true,
    "showTopN": 5,
    "showSkills": true,
    "showMCP": true
  },
  "tables": {
    "toolsTableThreshold": 5,
    "tasksTableThreshold": 3,
    "contextTableThreshold": 999
  }
}
```

Create `examples/config-mocha.json`:
```json
{
  "theme": "mocha",
  "preset": "full",
  "lineLayout": "expanded",
  "tables": {
    "toolsTableThreshold": 5,
    "tasksTableThreshold": 3,
    "contextTableThreshold": 999
  }
}
```

Create `examples/config-custom-colors.json`:
```json
{
  "theme": "macchiato",
  "colors": {
    "success": "#00ff00",
    "warning": "#ffaa00",
    "danger": "#ff0000",
    "primary": "#aa00ff"
  },
  "preset": "full",
  "tables": {
    "toolsTableThreshold": 8,
    "tasksTableThreshold": 5
  }
}
```

**Step 2: Commit example configs**

```bash
mkdir -p examples
git add examples/
git commit -m "docs: add example configuration files"
```

---

### Task 16: Update README with Theme Documentation

**Files:**
- Modify: `README.md`

**Step 1: Add theme section to README**

Add to `README.md` (after the Features section):

```markdown
## Themes

cc-hud-go supports [Catppuccin](https://github.com/catppuccin/catppuccin) themes with full customization.

### Available Themes

- **macchiato** (default) - Dark theme with warm tones
- **mocha** - Darker theme, most popular
- **frappe** - Medium-dark theme
- **latte** - Light theme

### Configuration

Set your theme in `~/.config/cc-hud-go/config.json`:

```json
{
  "theme": "macchiato"
}
```

### Custom Colors

Override specific semantic colors:

```json
{
  "theme": "macchiato",
  "colors": {
    "success": "#00ff00",
    "warning": "#ffaa00",
    "danger": "#ff0000"
  }
}
```

### Semantic Color Names

- `success`, `warning`, `danger` - Status indicators
- `input`, `output` - Data flow
- `cacheRead`, `cacheWrite` - Cache operations
- `primary`, `highlight`, `accent` - UI elements
- `muted`, `bright`, `info` - Utility colors

See `examples/` for complete configuration examples.

## Visual Features

### Gradient Progress Bars

Context usage and rate limits now display with smooth gradient bars that shift from green → yellow → orange → red as usage increases.

```
🟢 █████▓▓▒▒░░ 54%  │  📊 ████████▓▒ 87%
```

### Contextual Tables

When data exceeds thresholds, segments automatically switch to table view for better readability:

**Tools (>5 tools):**
```
┌──────────┬───────┐
│ Category │ Count │
├──────────┼───────┤
│ App      │   12  │
│ MCP      │    3  │
│ Skills   │    2  │
└──────────┴───────┘
```

**Tasks (>3 tasks):**
Shows last 3 completed + all active + all pending
```
┌─────────────────────┬──────────┐
│ Task                │ Status   │
├─────────────────────┼──────────┤
│ Implement auth      │ ✓ Done   │
│ Write tests         │ ⏵ Active │
│ Deploy staging      │ ○ Pending│
└─────────────────────┴──────────┘
```

Configure thresholds in your config:
```json
{
  "tables": {
    "toolsTableThreshold": 5,
    "tasksTableThreshold": 3,
    "contextTableThreshold": 999
  }
}
```
```

**Step 2: Commit README updates**

```bash
git add README.md
git commit -m "docs: add theme and visual features documentation"
```

---

### Task 17: Run Full Test Suite

**Step 1: Run all tests**

Run: `make test` or `go test ./... -v`
Expected: ALL PASS

**Step 2: Check test coverage**

Run: `make test-coverage` or `go test ./... -cover`
Expected: Good coverage across all packages

**Step 3: If any tests fail, fix them before proceeding**

Address any failing tests, then re-run.

**Step 4: Commit any fixes**

```bash
git add .
git commit -m "fix: address test failures"
```

---

### Task 18: Build and Manual Test

**Step 1: Build the binary**

Run: `make build` or `go build -o cc-hud-go .`
Expected: Successful build

**Step 2: Create test config**

Create `~/.config/cc-hud-go/config.json`:
```json
{
  "theme": "macchiato",
  "preset": "full",
  "tables": {
    "toolsTableThreshold": 2,
    "tasksTableThreshold": 1
  }
}
```

**Step 3: Test with sample data**

```bash
# Test with sample stdin data
echo '{"context": {"usedTokens": 54000, "totalTokens": 100000}}' | ./cc-hud-go
```

Expected: See gradient bar with Catppuccin colors

**Step 4: Verify visual output**

Check that:
- Colors match Catppuccin Macchiato palette
- Gradient bar displays with █▓▒░ characters
- Separators have 2 spaces on each side
- Tables appear when thresholds exceeded

**Step 5: Commit if any adjustments needed**

```bash
git add .
git commit -m "fix: adjust visual rendering based on manual testing"
```

---

## Summary

This implementation adds:

1. **Theme System** - Catppuccin support with 4 flavors + custom color overrides
2. **Gradient Progress Bars** - Smooth green→red gradients for context & rate limits
3. **Enhanced Spacing** - Two-space separators for better scannability
4. **Contextual Tables** - Automatic table view when data exceeds thresholds
5. **Task Filtering** - Shows last 3 completed + active + pending tasks

All changes follow TDD approach with comprehensive test coverage.

## Next Steps

After implementation:
- Consider adding visual regression tests with screenshots
- Add CI/CD pipeline for testing
- Create additional example configs for different use cases
- Add theme preview command (e.g., `cc-hud-go --preview-theme mocha`)

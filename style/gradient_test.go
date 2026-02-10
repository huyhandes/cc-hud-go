package style

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
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

// Mock theme for testing
type mockTheme struct{}

func (m *mockTheme) Name() string { return "test" }
func (m *mockTheme) GetColor(semantic string) lipgloss.Color {
	return lipgloss.Color("#ff0000")
}

func TestInitWithTheme(t *testing.T) {
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

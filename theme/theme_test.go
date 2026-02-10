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

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

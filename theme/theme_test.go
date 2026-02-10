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

func TestLoadFromConfig(t *testing.T) {
	wrapper := LoadThemeFromConfig("mocha", map[string]string{"success": "#00ff00"})

	if wrapper.Name() != "mocha" {
		t.Errorf("Expected theme 'mocha', got %s", wrapper.Name())
	}

	successColor := wrapper.GetColor("success")
	if string(successColor) != "#00ff00" {
		t.Errorf("Expected override color #00ff00, got %s", successColor)
	}

	warningColor := wrapper.GetColor("warning")
	if warningColor == lipgloss.Color("") {
		t.Error("Expected warning color from theme")
	}
}

func TestLoadFromConfigNoOverrides(t *testing.T) {
	th := LoadThemeFromConfig("frappe", nil)

	if th.Name() != "frappe" {
		t.Errorf("expected 'frappe', got %s", th.Name())
	}

	// Should be the base theme directly (not a wrapper)
	if _, ok := th.(*ThemeWrapper); ok {
		t.Error("expected base theme, not wrapper, when no overrides")
	}
}

func TestLoadFromConfigEmptyOverrides(t *testing.T) {
	th := LoadThemeFromConfig("latte", map[string]string{})

	if th.Name() != "latte" {
		t.Errorf("expected 'latte', got %s", th.Name())
	}

	if _, ok := th.(*ThemeWrapper); ok {
		t.Error("expected base theme with empty overrides map")
	}
}

func TestAllSemanticColorsExist(t *testing.T) {
	semantics := []string{
		"success", "warning", "danger", "input", "output",
		"cacheRead", "cacheWrite", "primary", "highlight",
		"accent", "muted", "bright", "info",
	}

	for _, flavor := range []string{"macchiato", "mocha", "frappe", "latte"} {
		th := GetTheme(flavor)
		for _, key := range semantics {
			color := th.GetColor(key)
			if color == lipgloss.Color("") {
				t.Errorf("%s: missing color for %q", flavor, key)
			}
		}
	}
}

func TestThemeWrapperOverrideFallback(t *testing.T) {
	wrapper := &ThemeWrapper{
		base:      NewMacchiato(),
		overrides: map[string]string{"custom_key": "#123456"},
	}

	// Override hit
	if string(wrapper.GetColor("custom_key")) != "#123456" {
		t.Error("expected override color")
	}

	// Base fallback
	baseColor := NewMacchiato().GetColor("success")
	if wrapper.GetColor("success") != baseColor {
		t.Error("expected base theme color for non-overridden key")
	}
}

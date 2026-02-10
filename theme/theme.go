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

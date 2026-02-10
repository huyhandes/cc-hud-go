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

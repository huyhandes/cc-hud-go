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

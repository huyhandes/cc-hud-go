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

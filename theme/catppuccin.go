package theme

import "github.com/charmbracelet/lipgloss"

// catppuccin is the single concrete theme type. All four flavors share the
// same shape; only the color map differs.
type catppuccin struct {
	name   string
	colors map[string]string
}

func (c *catppuccin) Name() string { return c.name }

func (c *catppuccin) GetColor(semantic string) lipgloss.Color {
	if color, ok := c.colors[semantic]; ok {
		return lipgloss.Color(color)
	}
	return lipgloss.Color(c.colors["bright"])
}

func NewMacchiato() *catppuccin {
	return &catppuccin{name: "macchiato", colors: catppuccinMacchiato}
}

func NewMocha() *catppuccin {
	return &catppuccin{name: "mocha", colors: catppuccinMocha}
}

func NewFrappe() *catppuccin {
	return &catppuccin{name: "frappe", colors: catppuccinFrappe}
}

func NewLatte() *catppuccin {
	return &catppuccin{name: "latte", colors: catppuccinLatte}
}

var catppuccinMacchiato = map[string]string{
	"success":    "#a6da95",
	"warning":    "#eed49f",
	"danger":     "#ed8796",
	"input":      "#8aadf4",
	"output":     "#8bd5ca",
	"cacheRead":  "#c6a0f6",
	"cacheWrite": "#f5bde6",
	"primary":    "#c6a0f6",
	"highlight":  "#91d7e3",
	"accent":     "#f5a97f",
	"muted":      "#5b6078",
	"bright":     "#cad3f5",
	"info":       "#8bd5ca",
}

var catppuccinMocha = map[string]string{
	"success":    "#a6e3a1",
	"warning":    "#f9e2af",
	"danger":     "#f38ba8",
	"input":      "#89b4fa",
	"output":     "#94e2d5",
	"cacheRead":  "#cba6f7",
	"cacheWrite": "#f5c2e7",
	"primary":    "#cba6f7",
	"highlight":  "#89dceb",
	"accent":     "#fab387",
	"muted":      "#585b70",
	"bright":     "#cdd6f4",
	"info":       "#94e2d5",
}

var catppuccinFrappe = map[string]string{
	"success":    "#a6d189",
	"warning":    "#e5c890",
	"danger":     "#e78284",
	"input":      "#8caaee",
	"output":     "#81c8be",
	"cacheRead":  "#ca9ee6",
	"cacheWrite": "#f4b8e4",
	"primary":    "#ca9ee6",
	"highlight":  "#99d1db",
	"accent":     "#ef9f76",
	"muted":      "#51576d",
	"bright":     "#c6d0f5",
	"info":       "#81c8be",
}

var catppuccinLatte = map[string]string{
	"success":    "#40a02b",
	"warning":    "#df8e1d",
	"danger":     "#d20f39",
	"input":      "#1e66f5",
	"output":     "#179299",
	"cacheRead":  "#8839ef",
	"cacheWrite": "#ea76cb",
	"primary":    "#8839ef",
	"highlight":  "#04a5e5",
	"accent":     "#fe640b",
	"muted":      "#9ca0b0",
	"bright":     "#4c4f69",
	"info":       "#179299",
}

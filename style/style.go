package style

import (
	"fmt"
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

// GetRenderer returns the global renderer
func GetRenderer() *lipgloss.Renderer {
	return renderer
}

// NewRendererForWriter creates a new renderer for a specific writer
func NewRendererForWriter(w io.Writer) *lipgloss.Renderer {
	return lipgloss.NewRenderer(w)
}

// Separator renders a styled separator
func Separator() string {
	return SeparatorStyle.Render("│")
}

// Icon renders a styled icon
func Icon(icon string, style lipgloss.Style) string {
	return style.Render(icon)
}

// RenderGradientBar renders a static gradient progress bar
// The gradient is always green → yellow → orange → red (0-100%)
// Only the filled portion (based on percentage) is displayed
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
			// Calculate position in the full 0-100% gradient
			positionPercent := (float64(i) / float64(width)) * 100

			// Get color for this position in the static gradient
			color := getStaticGradientColor(positionPercent)
			segments = append(segments, renderer.NewStyle().Foreground(color).Render("█"))
		} else {
			// Empty portion
			segments = append(segments, renderer.NewStyle().Foreground(ColorMuted).Render("░"))
		}
	}

	return strings.Join(segments, "")
}

// getStaticGradientColor returns a color from the static gradient (0-100%)
// Gradient: green (0%) → yellow (50%) → orange (75%) → red (100%)
func getStaticGradientColor(position float64) lipgloss.Color {
	// Define gradient stops with RGB colors
	// Using TrueColor values for smooth transitions
	var r, g, b uint8

	if position < 50 {
		// Green (0-50%): #a6da95 → #eed49f (green to yellow)
		// Interpolate between green and yellow
		t := position / 50
		r = lerp(0xa6, 0xee, t)
		g = lerp(0xda, 0xd4, t)
		b = lerp(0x95, 0x9f, t)
	} else if position < 75 {
		// Yellow to Orange (50-75%): #eed49f → #f5a97f
		t := (position - 50) / 25
		r = lerp(0xee, 0xf5, t)
		g = lerp(0xd4, 0xa9, t)
		b = lerp(0x9f, 0x7f, t)
	} else {
		// Orange to Red (75-100%): #f5a97f → #ed8796
		t := (position - 75) / 25
		r = lerp(0xf5, 0xed, t)
		g = lerp(0xa9, 0x87, t)
		b = lerp(0x7f, 0x96, t)
	}

	return lipgloss.Color(formatRGB(r, g, b))
}

// lerp performs linear interpolation between two values
func lerp(start, end uint8, t float64) uint8 {
	return uint8(float64(start) + (float64(end)-float64(start))*t)
}

// formatRGB formats RGB values as a hex color string
func formatRGB(r, g, b uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// RenderTable renders a table with headers and rows using lipgloss
func RenderTable(headers []string, rows [][]string) string {
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

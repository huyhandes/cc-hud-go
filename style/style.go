package style

import (
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

// RenderGradientBar renders a gradient progress bar
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
			// Use gradient characters for filled portion
			char := getGradientChar(i, filled, width)
			color := getColorForPercentage(percentage)
			segments = append(segments, renderer.NewStyle().Foreground(color).Render(char))
		} else {
			// Empty portion
			segments = append(segments, renderer.NewStyle().Foreground(ColorMuted).Render("░"))
		}
	}

	return strings.Join(segments, "")
}

// getGradientChar returns the appropriate gradient character
func getGradientChar(position, filled, width int) string {
	if filled == 0 {
		return "░"
	}

	// Use different characters based on position in filled area
	progress := float64(position) / float64(filled)

	if progress < 0.3 {
		return "█" // solid
	} else if progress < 0.6 {
		return "▓" // dark
	} else {
		return "▒" // medium
	}
}

// getColorForPercentage returns color based on percentage thresholds
func getColorForPercentage(percentage float64) lipgloss.Color {
	if percentage >= 90 {
		return ColorDanger // red
	} else if percentage >= 70 {
		return ColorWarning // yellow/orange
	}
	return ColorSuccess // green
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

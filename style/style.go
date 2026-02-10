package style

import (
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	// Global renderer that forces color output
	renderer *lipgloss.Renderer

	// Color palette - organized by semantic meaning

	// Status colors (usage levels)
	ColorSuccess = lipgloss.Color("#10B981") // Green - healthy/good
	ColorWarning = lipgloss.Color("#F59E0B") // Orange - caution
	ColorDanger  = lipgloss.Color("#EF4444") // Red - critical

	// Flow colors (data movement)
	ColorInput  = lipgloss.Color("#3B82F6") // Blue - incoming data
	ColorOutput = lipgloss.Color("#10B981") // Emerald - outgoing data

	// Cache colors (storage layer)
	ColorCacheRead  = lipgloss.Color("#8B5CF6") // Purple - cache read
	ColorCacheWrite = lipgloss.Color("#EC4899") // Pink - cache write

	// Primary UI colors
	ColorPrimary   = lipgloss.Color("#7C3AED") // Purple - model/agent
	ColorHighlight = lipgloss.Color("#06B6D4") // Cyan - git/highlights
	ColorAccent    = lipgloss.Color("#F59E0B") // Orange - cost/emphasis

	// Utility colors
	ColorMuted  = lipgloss.Color("#6B7280") // Gray - separators/static
	ColorBright = lipgloss.Color("#F3F4F6") // Light gray
	ColorInfo   = lipgloss.Color("#14B8A6") // Teal - information

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

	// Initialize styles with renderer
	ModelStyle = renderer.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)

	ContextStyle = renderer.NewStyle().
		Foreground(ColorInfo)

	GitStyle = renderer.NewStyle().
		Foreground(ColorHighlight)

	CostStyle = renderer.NewStyle().
		Foreground(ColorWarning)

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

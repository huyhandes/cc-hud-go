package style

import (
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	// Global renderer that forces color output
	renderer *lipgloss.Renderer

	// Color palette
	ColorPrimary   = lipgloss.Color("#7C3AED") // Purple
	ColorSuccess   = lipgloss.Color("#10B981") // Green
	ColorWarning   = lipgloss.Color("#F59E0B") // Orange
	ColorDanger    = lipgloss.Color("#EF4444") // Red
	ColorInfo      = lipgloss.Color("#3B82F6") // Blue
	ColorCyan      = lipgloss.Color("#06B6D4") // Cyan
	ColorMuted     = lipgloss.Color("#6B7280") // Gray
	ColorBright    = lipgloss.Color("#F3F4F6") // Light gray

	// Pre-configured styles
	ModelStyle    lipgloss.Style
	ContextStyle  lipgloss.Style
	GitStyle      lipgloss.Style
	CostStyle     lipgloss.Style
	ToolsStyle    lipgloss.Style
	AgentStyle    lipgloss.Style
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
		Foreground(ColorCyan)

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

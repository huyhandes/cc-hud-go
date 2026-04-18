package style

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/huyhandes/cc-hud-go/theme"
	"github.com/muesli/termenv"
)

var (
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
	renderer = lipgloss.NewRenderer(os.Stdout, termenv.WithProfile(termenv.TrueColor))
	renderer.SetColorProfile(termenv.TrueColor)
}

// Init initializes styles with the given theme
func Init(th theme.Theme) {
	ColorSuccess = th.GetColor("success")
	ColorWarning = th.GetColor("warning")
	ColorDanger = th.GetColor("danger")
	ColorInput = th.GetColor("input")
	ColorOutput = th.GetColor("output")
	ColorCacheRead = th.GetColor("cacheRead")
	ColorCacheWrite = th.GetColor("cacheWrite")
	ColorPrimary = th.GetColor("primary")
	ColorHighlight = th.GetColor("highlight")
	ColorAccent = th.GetColor("accent")
	ColorMuted = th.GetColor("muted")
	ColorBright = th.GetColor("bright")
	ColorInfo = th.GetColor("info")

	ModelStyle = renderer.NewStyle().Foreground(ColorPrimary).Bold(true)
	ContextStyle = renderer.NewStyle().Foreground(ColorInfo)
	GitStyle = renderer.NewStyle().Foreground(ColorHighlight)
	CostStyle = renderer.NewStyle().Foreground(ColorAccent)
	ToolsStyle = renderer.NewStyle().Foreground(ColorSuccess)
	AgentStyle = renderer.NewStyle().Foreground(ColorPrimary).Italic(true)
	SeparatorStyle = renderer.NewStyle().Foreground(ColorMuted)

	ProgressGood = renderer.NewStyle().Foreground(ColorSuccess)
	ProgressWarning = renderer.NewStyle().Foreground(ColorWarning)
	ProgressDanger = renderer.NewStyle().Foreground(ColorDanger)
}

// GetRenderer returns the global renderer
func GetRenderer() *lipgloss.Renderer {
	return renderer
}

// Separator renders a styled separator
func Separator() string {
	return SeparatorStyle.Render("│")
}

// Icon renders a styled icon
func Icon(icon string, s lipgloss.Style) string {
	return s.Render(icon)
}

// ThresholdColor returns a color based on percentage thresholds (green/yellow/red)
func ThresholdColor(percentage float64) lipgloss.Color {
	if percentage >= 90 {
		return ColorDanger
	}
	if percentage >= 70 {
		return ColorWarning
	}
	return ColorSuccess
}

// ContextTokenColor returns a color for absolute token usage.
// Thresholds tuned for context cliff: green < 150k, warn 150-175k,
// red 175-230k, critical >= 230k.
func ContextTokenColor(tokens int) lipgloss.Color {
	switch {
	case tokens >= 230_000:
		return lipgloss.Color("#ff0000") // CRITICAL
	case tokens >= 175_000:
		return ColorDanger
	case tokens >= 150_000:
		return ColorWarning
	default:
		return ColorSuccess
	}
}

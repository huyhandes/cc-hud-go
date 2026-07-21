package style

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/huyhandes/cc-hud-go/theme"
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

	// Pre-configured styles that are actually used
	ModelStyle lipgloss.Style
	AgentStyle lipgloss.Style
)

func init() {
	renderer = lipgloss.NewRenderer(os.Stdout)
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
	AgentStyle = renderer.NewStyle().Foreground(ColorPrimary).Italic(true)
}

// GetRenderer returns the global renderer
func GetRenderer() *lipgloss.Renderer {
	return renderer
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

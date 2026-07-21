package segment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
)

type CavemanSegment struct{}

func (c *CavemanSegment) ID() string { return "caveman" }

func (c *CavemanSegment) Render(_ *state.State) (string, error) {
	return renderModeBadge(".caveman-active", "🦴", "CAVEMAN"), nil
}

// renderModeBadge renders an "emoji LABEL:LEVEL" badge from a flag file under
// the Claude config dir. Returns "" when the file is absent, unreadable, or
// holds an unrecognized level.
func renderModeBadge(flagFile, emoji, label string) string {
	level := readModeLevel(flagFile)
	if level == "" {
		return ""
	}

	var color lipgloss.Color
	switch level {
	case "ultra":
		color = style.ColorDanger
	case "full":
		color = style.ColorWarning
	case "lite":
		color = style.ColorInfo
	}

	badge := style.GetRenderer().NewStyle().Foreground(color).Bold(true).
		Render(fmt.Sprintf("%s:%s", label, strings.ToUpper(level)))
	return emoji + " " + badge
}

// readModeLevel reads a canonical mode level from flagFile under the Claude
// config dir. Empty string for any non-canonical value.
func readModeLevel(flagFile string) string {
	dir := claudeConfigDir()
	data, err := os.ReadFile(filepath.Join(dir, flagFile))
	if err != nil {
		return ""
	}
	level := strings.ToLower(strings.TrimSpace(string(data)))
	switch level {
	case "full", "lite", "ultra":
		return level
	}
	return ""
}

// claudeConfigDir resolves the Claude config directory.
//
// ponytail: CLAUDE_CONFIG_DIR is honored (not just $HOME/.claude) so this
// segment and ponytail's own statusline script read state from the same
// place — they can never disagree about where the flag file lives.
func claudeConfigDir() string {
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".claude")
}

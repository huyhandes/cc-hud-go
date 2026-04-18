package segment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
)

type CavemanSegment struct{}

func (c *CavemanSegment) ID() string { return "caveman" }

func (c *CavemanSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Caveman
}

func (c *CavemanSegment) Render(_ *state.State, _ *config.Config) (string, error) {
	level := readCavemanLevel()
	if level == "" {
		return "", nil
	}

	var color lipgloss.Color
	switch level {
	case "ultra":
		color = style.ColorDanger
	case "full":
		color = style.ColorWarning
	case "lite":
		color = style.ColorInfo
	default:
		color = style.ColorMuted
	}

	badge := style.GetRenderer().NewStyle().Foreground(color).Bold(true).
		Render(fmt.Sprintf("CAVEMAN:%s", strings.ToUpper(level)))
	return "🦴 " + badge, nil
}

func readCavemanLevel() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(home, ".claude", ".caveman-active"))
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

package segment

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

var toolsStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14")) // Cyan

type ToolsSegment struct{}

func (t *ToolsSegment) ID() string {
	return "tools"
}

func (t *ToolsSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Tools
}

func (t *ToolsSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	// Count totals
	appTotal := 0
	for _, count := range s.Tools.AppTools {
		appTotal += count
	}

	internalTotal := 0
	for _, count := range s.Tools.InternalTools {
		internalTotal += count
	}

	customTotal := 0
	for _, count := range s.Tools.CustomTools {
		customTotal += count
	}

	mcpTotal := 0
	for _, tools := range s.Tools.MCPTools {
		for _, count := range tools {
			mcpTotal += count
		}
	}

	skillsTotal := 0
	for _, usage := range s.Tools.Skills {
		skillsTotal += usage.Count
	}

	total := appTotal + internalTotal + customTotal + mcpTotal + skillsTotal

	if total == 0 {
		return "", nil
	}

	// Build output
	var parts []string

	if cfg.Tools.GroupByCategory {
		if appTotal > 0 {
			parts = append(parts, fmt.Sprintf("App:%d", appTotal))
		}
		if mcpTotal > 0 && cfg.Tools.ShowMCP {
			parts = append(parts, fmt.Sprintf("MCP:%d", mcpTotal))
		}
		if skillsTotal > 0 && cfg.Tools.ShowSkills {
			parts = append(parts, fmt.Sprintf("Skills:%d", skillsTotal))
		}
		if customTotal > 0 {
			parts = append(parts, fmt.Sprintf("Custom:%d", customTotal))
		}

		return toolsStyle.Render(fmt.Sprintf("Tools: %d (%s)", total, strings.Join(parts, " "))), nil
	}

	return toolsStyle.Render(fmt.Sprintf("Tools: %d", total)), nil
}

package segment

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

var agentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13")) // Purple

type AgentSegment struct{}

func (a *AgentSegment) ID() string {
	return "agent"
}

func (a *AgentSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Agents
}

func (a *AgentSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	if s.Agents.ActiveAgent == "" {
		return "", nil
	}

	// Basic format: "Agent: <name>"
	output := fmt.Sprintf("Agent: %s", s.Agents.ActiveAgent)

	// Add task description if available
	if s.Agents.TaskDesc != "" {
		output = fmt.Sprintf("%s (%s)", output, s.Agents.TaskDesc)
	}

	return agentStyle.Render(output), nil
}

package segment

import (
	"fmt"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
	"github.com/huybui/cc-hud-go/style"
)

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

	// Add agent icon
	icon := "👤"
	output := fmt.Sprintf("%s %s", icon, s.Agents.ActiveAgent)

	// Add task description if available
	if s.Agents.TaskDesc != "" {
		output = fmt.Sprintf("%s (%s)", output, s.Agents.TaskDesc)
	}

	return style.AgentStyle.Render(output), nil
}

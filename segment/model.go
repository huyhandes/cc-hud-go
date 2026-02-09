package segment

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

var (
	modelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // Blue
	planStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green
)

type ModelSegment struct{}

func (m *ModelSegment) ID() string {
	return "model"
}

func (m *ModelSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Model
}

func (m *ModelSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	if s.Model.Name == "" {
		return "", nil
	}

	model := modelStyle.Render(s.Model.Name)

	if s.Model.PlanType != "" {
		plan := planStyle.Render(fmt.Sprintf("[%s]", s.Model.PlanType))
		return fmt.Sprintf("%s %s", model, plan), nil
	}

	return model, nil
}

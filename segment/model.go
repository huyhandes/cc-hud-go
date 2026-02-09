package segment

import (
	"fmt"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
	"github.com/huybui/cc-hud-go/style"
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

	// Add AI/robot icon for model
	icon := "🤖"
	model := style.ModelStyle.Render(s.Model.Name)

	if s.Model.PlanType != "" {
		return fmt.Sprintf("%s %s", icon, model), nil
	}

	return fmt.Sprintf("%s %s", icon, model), nil
}

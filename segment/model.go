package segment

import (
	"fmt"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
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

	model := style.ModelStyle.Render(s.Model.Name)
	return fmt.Sprintf("🤖 %s", model), nil
}

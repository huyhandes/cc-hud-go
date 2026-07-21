package segment

import (
	"fmt"

	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
)

type ModelSegment struct{}

func (m *ModelSegment) ID() string {
	return "model"
}

func (m *ModelSegment) Render(s *state.State) string {
	if s.Model.Name == "" {
		return ""
	}

	model := style.ModelStyle.Render(s.Model.Name)
	return fmt.Sprintf("🤖 %s", model)
}

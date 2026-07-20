package segment

import (
	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
)

type PonytailSegment struct{}

func (p *PonytailSegment) ID() string { return "ponytail" }

func (p *PonytailSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Ponytail
}

func (p *PonytailSegment) Render(_ *state.State, _ *config.Config) (string, error) {
	return renderModeBadge(".ponytail-active", "🐴", "PONYTAIL"), nil
}

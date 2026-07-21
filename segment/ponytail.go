package segment

import "github.com/huyhandes/cc-hud-go/state"

type PonytailSegment struct{}

func (p *PonytailSegment) ID() string { return "ponytail" }

func (p *PonytailSegment) Render(_ *state.State) (string, error) {
	return renderModeBadge(".ponytail-active", "🐴", "PONYTAIL"), nil
}

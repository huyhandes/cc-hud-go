package segment

import (
	"testing"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
)

func TestRegistry(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	segments := All()

	if len(segments) == 0 {
		t.Error("expected at least one segment")
	}

	// Check that segments implement interface
	for _, seg := range segments {
		if seg.ID() == "" {
			t.Error("segment ID should not be empty")
		}

		// Should be able to check if enabled
		_ = seg.Enabled(cfg)

		// Should be able to render
		_, err := seg.Render(s, cfg)
		if err != nil {
			t.Errorf("segment %s render failed: %v", seg.ID(), err)
		}
	}
}

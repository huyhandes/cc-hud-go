package segment

import (
	"testing"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
)

func TestByID(t *testing.T) {
	m := ByID()
	all := All()

	if len(m) != len(all) {
		t.Errorf("ByID() returned %d segments, All() returned %d", len(m), len(all))
	}

	for _, seg := range all {
		got, ok := m[seg.ID()]
		if !ok {
			t.Errorf("ByID() missing segment %q", seg.ID())
			continue
		}
		if got.ID() != seg.ID() {
			t.Errorf("ByID()[%q].ID() = %q", seg.ID(), got.ID())
		}
	}
}

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

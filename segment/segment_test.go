package segment

import (
	"testing"

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
	s := state.New()

	segments := All()

	wantIDs := []string{"model", "caveman", "ponytail", "git", "fivehour", "ratelimit"}
	if len(segments) != len(wantIDs) {
		t.Fatalf("expected %d segments, got %d", len(wantIDs), len(segments))
	}
	for i, want := range wantIDs {
		if segments[i].ID() != want {
			t.Errorf("segment[%d]: want ID %q, got %q", i, want, segments[i].ID())
		}
	}

	// Check that segments implement interface
	for _, seg := range segments {
		if seg.ID() == "" {
			t.Error("segment ID should not be empty")
		}

		// Should be able to render
		_, err := seg.Render(s)
		if err != nil {
			t.Errorf("segment %s render failed: %v", seg.ID(), err)
		}
	}
}

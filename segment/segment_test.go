package segment

import (
	"testing"

	"github.com/huyhandes/cc-hud-go/state"
)

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

		// Should be able to render (self-gates to "" when no data)
		_ = seg.Render(s)
	}
}

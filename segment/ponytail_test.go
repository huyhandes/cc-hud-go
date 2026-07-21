package segment

import (
	"strings"
	"testing"
)

// TestPonytailRegistry covers wiring only. The renderModeBadge helper itself
// is covered by caveman_test.go — no re-test here.
func TestPonytailRegistry(t *testing.T) {
	if _, ok := ByID()["ponytail"]; !ok {
		t.Errorf(`ByID() missing key "ponytail"`)
	}
}

func TestPonytailSegmentID(t *testing.T) {
	seg := &PonytailSegment{}
	if got := seg.ID(); got != "ponytail" {
		t.Errorf("PonytailSegment.ID() = %q, want %q", got, "ponytail")
	}
}

// TestPonytailSegmentDelegates confirms Render wires through to
// renderModeBadge with the right flag/emoji/label.
func TestPonytailSegmentDelegates(t *testing.T) {
	dir := withConfigDir(t)
	writeFlag(t, dir, ".ponytail-active", "ultra")

	seg := &PonytailSegment{}
	out, err := seg.Render(nil)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	if got := stripANSI(out); !strings.Contains(got, "🐴 PONYTAIL:ULTRA") {
		t.Errorf("Render = %q, want it to contain %q", got, "🐴 PONYTAIL:ULTRA")
	}
}

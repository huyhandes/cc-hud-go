package state

import "testing"

func TestNewState(t *testing.T) {
	s := New()

	if s == nil {
		t.Fatal("expected non-nil state")
	}

	if s.Session.StartTime.IsZero() {
		t.Error("expected StartTime to be set")
	}

	if s.Tools.AppTools == nil {
		t.Error("expected AppTools map to be initialized")
	}
}

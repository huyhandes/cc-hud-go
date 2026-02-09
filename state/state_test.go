package state

import (
	"testing"
	"time"
)

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

func TestUpdateDerived(t *testing.T) {
	s := New()

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	s.UpdateDerived()

	if s.Session.Duration == 0 {
		t.Error("expected Duration to be updated")
	}

	if s.Session.Duration < 100*time.Millisecond {
		t.Errorf("expected Duration >= 100ms, got %v", s.Session.Duration)
	}

	// Test percentage calculation
	s.Context.UsedTokens = 50
	s.Context.TotalTokens = 100
	s.UpdateDerived()

	if s.Context.Percentage != 50.0 {
		t.Errorf("expected Percentage 50.0, got %f", s.Context.Percentage)
	}
}

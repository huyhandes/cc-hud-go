package segment

import (
	"strings"
	"testing"

	"github.com/huyhandes/cc-hud-go/state"
)

func TestRateLimitSegment(t *testing.T) {
	s := state.New()
	s.RateLimits.SevenDayPercent = 75

	seg := &RateLimitSegment{}

	output := seg.Render(s)
	if !strings.Contains(output, "75%") {
		t.Errorf("expected percentage in output, got '%s'", output)
	}
}

func TestRateLimitSegmentEmpty(t *testing.T) {
	s := state.New()

	seg := &RateLimitSegment{}

	output := seg.Render(s)
	if output != "" {
		t.Errorf("expected empty output with no rate limit data, got '%s'", output)
	}
}

func TestRateLimitSegmentHighUsage(t *testing.T) {
	s := state.New()
	s.RateLimits.SevenDayPercent = 85

	seg := &RateLimitSegment{}

	output := seg.Render(s)
	if !strings.Contains(output, "85%") {
		t.Errorf("expected percentage in output, got '%s'", output)
	}
}

func TestRateLimitUsesGradientBar(t *testing.T) {
	s := state.New()
	s.RateLimits.SevenDayPercent = 67

	seg := &RateLimitSegment{}
	result := seg.Render(s)

	// Should contain gradient bar characters
	hasGradient := strings.Contains(result, "█") ||
		strings.Contains(result, "▓") ||
		strings.Contains(result, "▒") ||
		strings.Contains(result, "░")

	if !hasGradient {
		t.Error("Expected gradient bar characters in rate limit segment")
	}
}

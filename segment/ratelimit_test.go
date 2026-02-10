package segment

import (
	"strings"
	"testing"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
)

func TestRateLimitSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.RateLimits.SevenDayUsed = 75
	s.RateLimits.SevenDayTotal = 100

	seg := &RateLimitSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "75%") {
		t.Errorf("expected percentage in output, got '%s'", output)
	}
}

func TestRateLimitSegmentEmpty(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	seg := &RateLimitSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if output != "" {
		t.Errorf("expected empty output with no rate limit data, got '%s'", output)
	}
}

func TestRateLimitSegmentHighUsage(t *testing.T) {
	cfg := config.Default()
	cfg.SevenDayThreshold = 80
	s := state.New()
	s.RateLimits.SevenDayUsed = 85
	s.RateLimits.SevenDayTotal = 100

	seg := &RateLimitSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Should show warning when usage exceeds threshold
	if !strings.Contains(output, "85%") {
		t.Errorf("expected percentage in output, got '%s'", output)
	}
}

func TestRateLimitUsesGradientBar(t *testing.T) {
	s := state.New()
	s.RateLimits.SevenDayUsed = 67
	s.RateLimits.SevenDayTotal = 100

	cfg := config.Default()

	seg := &RateLimitSegment{}
	result, err := seg.Render(s, cfg)

	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should contain gradient bar characters
	hasGradient := strings.Contains(result, "█") ||
		strings.Contains(result, "▓") ||
		strings.Contains(result, "▒") ||
		strings.Contains(result, "░")

	if !hasGradient {
		t.Error("Expected gradient bar characters in rate limit segment")
	}
}

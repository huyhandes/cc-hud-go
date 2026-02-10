package segment

import (
	"strings"
	"testing"
	"time"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

func TestFiveHourSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.RateLimits.FiveHourPercent = 15.0

	// Set reset time to 2 hours from now
	resetTime := time.Now().Add(2 * time.Hour)
	s.RateLimits.FiveHourResetsAt = resetTime.Format(time.RFC3339)

	seg := &FiveHourSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "15%") {
		t.Errorf("expected percentage in output, got '%s'", output)
	}

	if !strings.Contains(output, "⏱️") {
		t.Errorf("expected timer emoji in output, got '%s'", output)
	}

	// Should contain time remaining (format: 1h59m or similar)
	if !strings.Contains(output, "h") && !strings.Contains(output, "m") {
		t.Errorf("expected time remaining in output, got '%s'", output)
	}
}

func TestFiveHourSegmentEmpty(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	seg := &FiveHourSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if output != "" {
		t.Errorf("expected empty output with no 5h data, got '%s'", output)
	}
}

func TestFiveHourSegmentHighUsage(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.RateLimits.FiveHourPercent = 85.0

	// Set reset time to 30 minutes from now
	resetTime := time.Now().Add(30 * time.Minute)
	s.RateLimits.FiveHourResetsAt = resetTime.Format(time.RFC3339)

	seg := &FiveHourSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "85%") {
		t.Errorf("expected percentage in output, got '%s'", output)
	}

	// Should show minutes only when less than 1 hour
	if !strings.Contains(output, "m)") {
		t.Errorf("expected minutes in output, got '%s'", output)
	}
}

func TestFiveHourSegmentUsesGradientBar(t *testing.T) {
	s := state.New()
	s.RateLimits.FiveHourPercent = 45.0

	cfg := config.Default()

	seg := &FiveHourSegment{}
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
		t.Error("Expected gradient bar characters in 5h segment")
	}
}

func TestFiveHourSegmentTimeFormatting(t *testing.T) {
	tests := []struct {
		name          string
		remainingTime time.Duration
		expectHours   bool
		expectMinutes bool
	}{
		{
			name:          "hours and minutes",
			remainingTime: 2*time.Hour + 30*time.Minute,
			expectHours:   true,
			expectMinutes: true,
		},
		{
			name:          "minutes only",
			remainingTime: 45 * time.Minute,
			expectHours:   false,
			expectMinutes: true,
		},
		{
			name:          "almost done",
			remainingTime: 5 * time.Minute,
			expectHours:   false,
			expectMinutes: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Default()
			s := state.New()
			s.RateLimits.FiveHourPercent = 50.0

			resetTime := time.Now().Add(tt.remainingTime)
			s.RateLimits.FiveHourResetsAt = resetTime.Format(time.RFC3339)

			seg := &FiveHourSegment{}
			output, err := seg.Render(s, cfg)

			if err != nil {
				t.Fatalf("render failed: %v", err)
			}

			if tt.expectHours && !strings.Contains(output, "h") {
				t.Errorf("expected 'h' in output for %s, got '%s'", tt.name, output)
			}

			if !tt.expectHours && strings.Contains(output, "h") {
				t.Errorf("did not expect 'h' in output for %s, got '%s'", tt.name, output)
			}

			if tt.expectMinutes && !strings.Contains(output, "m") {
				t.Errorf("expected 'm' in output for %s, got '%s'", tt.name, output)
			}
		})
	}
}

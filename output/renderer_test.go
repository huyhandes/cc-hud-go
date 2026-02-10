package output

import (
	"strings"
	"testing"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
	"github.com/huyhandes/cc-hud-go/theme"
)

func init() {
	style.Init(theme.NewMacchiato())
}

func TestRender(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Model.Name = "Sonnet 4.5"
	s.Context.UsedTokens = 5000
	s.Context.TotalTokens = 10000

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if output == "" {
		t.Error("expected non-empty output")
	}

	if !strings.Contains(output, "Sonnet") {
		t.Errorf("expected output to contain model name, got: %s", output)
	}

	if !strings.Contains(output, "│") {
		t.Errorf("expected output to contain separator '│', got: %s", output)
	}
}

func TestRenderWithDisabledSegments(t *testing.T) {
	cfg := config.Minimal()
	s := state.New()
	s.Model.Name = "Sonnet 4.5"

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if output == "" {
		t.Error("expected non-empty output")
	}

	if !strings.Contains(output, "Sonnet") {
		t.Errorf("expected model in output, got: %s", output)
	}
}

func TestRenderEmptyState(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	_, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
}

func TestRenderMultiLine(t *testing.T) {
	cfg := config.Default()
	cfg.LineLayout = "multiline"
	s := state.New()
	s.Model.Name = "Opus 4.6"
	s.Context.UsedTokens = 50000
	s.Context.TotalTokens = 200000
	s.Context.TotalInputTokens = 30000
	s.Context.TotalOutputTokens = 10000
	s.Context.CacheReadTokens = 5000
	s.Context.CacheCreateTokens = 3000
	s.Cost.TotalUSD = 0.0567
	s.Cost.DurationMs = 154000
	s.Cost.LinesAdded = 45
	s.Cost.LinesRemoved = 12

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("renderMultiLine failed: %v", err)
	}

	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(lines))
	}

	if !strings.Contains(output, "Opus") {
		t.Error("expected model name in output")
	}
	if !strings.Contains(output, "📥") {
		t.Error("expected input token icon")
	}
	if !strings.Contains(output, "📤") {
		t.Error("expected output token icon")
	}
	if !strings.Contains(output, "💾") {
		t.Error("expected cache icon")
	}
	if !strings.Contains(output, "💰") {
		t.Error("expected cost icon")
	}
	if !strings.Contains(output, "⏱") {
		t.Error("expected duration icon")
	}
	if !strings.Contains(output, "📝") {
		t.Error("expected file changes icon")
	}
}

func TestRenderMultiLineExpanded(t *testing.T) {
	cfg := config.Default()
	cfg.LineLayout = "expanded"
	s := state.New()
	s.Model.Name = "Sonnet 4.5"
	s.Context.UsedTokens = 100000
	s.Context.TotalTokens = 200000

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("render expanded failed: %v", err)
	}

	if !strings.Contains(output, "Sonnet") {
		t.Error("expected model name in expanded layout")
	}
}

func TestRenderMultiLineMinimalState(t *testing.T) {
	cfg := config.Default()
	cfg.LineLayout = "multiline"
	s := state.New()
	s.Model.Name = "Haiku 4.5"

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("render multiline empty state failed: %v", err)
	}

	if !strings.Contains(output, "Haiku") {
		t.Error("expected model name even with minimal state")
	}
	if strings.Contains(output, "💾") {
		t.Error("should not show cache tokens when they're zero")
	}
	if strings.Contains(output, "💰") {
		t.Error("should not show cost when zero")
	}
}

func TestRenderMultiLineNoCacheTokens(t *testing.T) {
	cfg := config.Default()
	cfg.LineLayout = "multiline"
	s := state.New()
	s.Model.Name = "Test"
	s.Context.UsedTokens = 5000
	s.Context.TotalTokens = 200000
	s.Context.TotalInputTokens = 3000
	s.Context.TotalOutputTokens = 2000

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(output, "💾") {
		t.Error("should not render cache tokens when both are zero")
	}
}

func TestRenderMultiLineWithRateLimits(t *testing.T) {
	cfg := config.Default()
	cfg.LineLayout = "multiline"
	s := state.New()
	s.Model.Name = "Test"
	s.Context.UsedTokens = 5000
	s.Context.TotalTokens = 200000
	s.RateLimits.HourlyUsed = 10
	s.RateLimits.HourlyTotal = 50
	s.RateLimits.SevenDayUsed = 300
	s.RateLimits.SevenDayTotal = 1000

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output == "" {
		t.Error("expected non-empty output with rate limits")
	}
}

func TestRenderContextSize(t *testing.T) {
	s := state.New()
	s.Context.TotalTokens = 200000
	result := renderContextSize(s)

	if !strings.Contains(result, "⚡") {
		t.Error("expected lightning icon")
	}
	if !strings.Contains(result, "200k") {
		t.Errorf("expected '200k' in context size, got: %s", result)
	}
}

func TestRenderContextBar(t *testing.T) {
	tests := []struct {
		name       string
		percentage float64
	}{
		{"low usage", 25.0},
		{"medium usage", 70.0},
		{"high usage", 95.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := state.New()
			s.Context.Percentage = tt.percentage
			result := renderContextBar(s)

			if !strings.Contains(result, "🧠") {
				t.Error("expected brain icon")
			}
			if !strings.Contains(result, "%") {
				t.Error("expected percentage")
			}
		})
	}
}

func TestRenderIOTokens(t *testing.T) {
	s := state.New()
	s.Context.TotalInputTokens = 89000
	s.Context.TotalOutputTokens = 12000

	result := renderIOTokens(s)

	if !strings.Contains(result, "📥") {
		t.Error("expected input icon")
	}
	if !strings.Contains(result, "📤") {
		t.Error("expected output icon")
	}
	if !strings.Contains(result, "89k") {
		t.Errorf("expected '89k', got: %s", result)
	}
	if !strings.Contains(result, "12k") {
		t.Errorf("expected '12k', got: %s", result)
	}
}

func TestRenderCacheTokens(t *testing.T) {
	s := state.New()
	s.Context.CacheReadTokens = 45000
	s.Context.CacheCreateTokens = 23000

	result := renderCacheTokens(s)

	if !strings.Contains(result, "💾") {
		t.Error("expected cache icon")
	}
	if !strings.Contains(result, "R:45k") {
		t.Errorf("expected 'R:45k', got: %s", result)
	}
	if !strings.Contains(result, "W:23k") {
		t.Errorf("expected 'W:23k', got: %s", result)
	}
}

func TestRenderCost(t *testing.T) {
	s := state.New()
	s.Cost.TotalUSD = 0.1234

	result := renderCost(s)

	if !strings.Contains(result, "💰") {
		t.Error("expected cost icon")
	}
	if !strings.Contains(result, "$0.1234") {
		t.Errorf("expected '$0.1234', got: %s", result)
	}
}

func TestRenderTime(t *testing.T) {
	s := state.New()
	s.Cost.DurationMs = 154000

	result := renderTime(s)

	if !strings.Contains(result, "⏱") {
		t.Error("expected timer icon")
	}
	if !strings.Contains(result, "2m34s") {
		t.Errorf("expected '2m34s', got: %s", result)
	}
}

func TestRenderFileChanges(t *testing.T) {
	tests := []struct {
		name    string
		added   int
		removed int
		want    string
		empty   bool
	}{
		{"both", 45, 12, "📝", false},
		{"add only", 10, 0, "+10", false},
		{"remove only", 0, 5, "-5", false},
		{"none", 0, 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := state.New()
			s.Cost.LinesAdded = tt.added
			s.Cost.LinesRemoved = tt.removed

			result := renderFileChanges(s)

			if tt.empty {
				if result != "" {
					t.Errorf("expected empty, got: %s", result)
				}
				return
			}

			if !strings.Contains(result, tt.want) {
				t.Errorf("expected %q in result, got: %s", tt.want, result)
			}
		})
	}
}

func TestJoinSegmentsWithSpacing(t *testing.T) {
	segments := []string{"segment1", "segment2", "segment3"}

	result := joinSegments(segments)

	if !strings.Contains(result, "  │  ") {
		t.Errorf("Expected two-space separator, got: %s", result)
	}

	sepCount := strings.Count(result, "│")
	if sepCount != 2 {
		t.Errorf("Expected 2 separators, got %d", sepCount)
	}
}

func TestJoinSegmentsEdgeCases(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		result := joinSegments(nil)
		if result != "" {
			t.Errorf("expected empty, got: %s", result)
		}
	})

	t.Run("single segment", func(t *testing.T) {
		result := joinSegments([]string{"only"})
		if result != "only" {
			t.Errorf("expected 'only', got: %s", result)
		}
	})

	t.Run("all empty", func(t *testing.T) {
		result := joinSegments([]string{"", "", "  "})
		if result != "" {
			t.Errorf("expected empty after filtering, got: %s", result)
		}
	})

	t.Run("filters empty between segments", func(t *testing.T) {
		result := joinSegments([]string{"a", "", "b"})
		if strings.Count(result, "│") != 1 {
			t.Errorf("expected 1 separator, got: %s", result)
		}
	})
}

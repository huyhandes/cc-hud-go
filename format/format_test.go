package format

import "testing"

func TestTokens(t *testing.T) {
	tests := []struct {
		tokens int
		want   string
	}{
		{0, "0"},
		{500, "500"},
		{999, "999"},
		{1000, "1k"},
		{5000, "5k"},
		{200000, "200k"},
		{1000000, "1000k"},
	}

	for _, tt := range tests {
		got := Tokens(tt.tokens)
		if got != tt.want {
			t.Errorf("Tokens(%d) = %q, want %q", tt.tokens, got, tt.want)
		}
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		ms   int64
		want string
	}{
		{0, "0s"},
		{5000, "5s"},
		{45000, "45s"},
		{60000, "1m0s"},
		{154000, "2m34s"},
		{3600000, "60m0s"},
	}

	for _, tt := range tests {
		got := Duration(tt.ms)
		if got != tt.want {
			t.Errorf("Duration(%d) = %q, want %q", tt.ms, got, tt.want)
		}
	}
}

func TestCost(t *testing.T) {
	tests := []struct {
		usd  float64
		want string
	}{
		{0, "$0.0000"},
		{0.0234, "$0.0234"},
		{1.5, "$1.5000"},
		{0.0001, "$0.0001"},
	}

	for _, tt := range tests {
		got := Cost(tt.usd)
		if got != tt.want {
			t.Errorf("Cost(%f) = %q, want %q", tt.usd, got, tt.want)
		}
	}
}

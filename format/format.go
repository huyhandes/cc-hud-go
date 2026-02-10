package format

import "fmt"

// Tokens formats a token count for display (e.g. 5000 → "5k", 200000 → "200k", 500 → "500")
func Tokens(tokens int) string {
	if tokens >= 1000 {
		return fmt.Sprintf("%dk", tokens/1000)
	}
	return fmt.Sprintf("%d", tokens)
}

// Duration formats milliseconds into a human-readable duration (e.g. "2m34s", "45s")
func Duration(ms int64) string {
	sec := ms / 1000
	mins := sec / 60
	secs := sec % 60

	if mins > 0 {
		return fmt.Sprintf("%dm%ds", mins, secs)
	}
	return fmt.Sprintf("%ds", secs)
}

// Cost formats a USD cost value (e.g. 0.0234 → "$0.0234")
func Cost(usd float64) string {
	return fmt.Sprintf("$%.4f", usd)
}

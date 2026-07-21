package style

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderGradientBar renders a static gradient progress bar
// The gradient is always green -> yellow -> orange -> red (0-100%)
// Only the filled portion (based on percentage) is displayed
func RenderGradientBar(percentage float64, width int) string {
	if width <= 0 {
		width = 10
	}
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}

	filled := min(int(percentage/100*float64(width)), width)

	segments := make([]string, 0, width)

	for i := 0; i < width; i++ {
		if i < filled {
			positionPercent := (float64(i) / float64(width)) * 100
			color := getStaticGradientColor(positionPercent)
			segments = append(segments, renderer.NewStyle().Foreground(color).Render("█"))
		} else {
			segments = append(segments, renderer.NewStyle().Foreground(ColorMuted).Render("░"))
		}
	}

	return strings.Join(segments, "")
}

// getStaticGradientColor returns a color from the static gradient (0-100%)
// Gradient: green (0%) -> yellow (50%) -> orange (75%) -> red (100%)
func getStaticGradientColor(position float64) lipgloss.Color {
	var r, g, b uint8

	if position < 50 {
		t := position / 50
		r = lerp(0xa6, 0xee, t)
		g = lerp(0xda, 0xd4, t)
		b = lerp(0x95, 0x9f, t)
	} else if position < 75 {
		t := (position - 50) / 25
		r = lerp(0xee, 0xf5, t)
		g = lerp(0xd4, 0xa9, t)
		b = lerp(0x9f, 0x7f, t)
	} else {
		t := (position - 75) / 25
		r = lerp(0xf5, 0xed, t)
		g = lerp(0xa9, 0x87, t)
		b = lerp(0x7f, 0x96, t)
	}

	return lipgloss.Color(formatRGB(r, g, b))
}

func lerp(start, end uint8, t float64) uint8 {
	return uint8(float64(start) + (float64(end)-float64(start))*t)
}

func formatRGB(r, g, b uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// RenderContextBar renders a context-usage bar using absolute token thresholds.
// Bar fill maps to used/total. Color zones (independent of total):
//
//	chill green < 150k, orange 150-175k, red 175-230k, ULTRA RED >= 230k.
func RenderContextBar(usedTokens, totalTokens, width int) string {
	if width <= 0 {
		width = 10
	}
	if totalTokens <= 0 {
		return strings.Repeat(" ", width)
	}
	if usedTokens < 0 {
		usedTokens = 0
	}
	if usedTokens > totalTokens {
		usedTokens = totalTokens
	}

	filled := min(int(float64(usedTokens)/float64(totalTokens)*float64(width)), width)

	tokensPerCell := float64(totalTokens) / float64(width)
	segments := make([]string, 0, width)
	for i := 0; i < width; i++ {
		if i < filled {
			cellTokens := int(float64(i+1) * tokensPerCell)
			color := contextZoneColor(cellTokens)
			segments = append(segments, renderer.NewStyle().Foreground(color).Render("█"))
		} else {
			segments = append(segments, renderer.NewStyle().Foreground(ColorMuted).Render("░"))
		}
	}
	return strings.Join(segments, "")
}

func contextZoneColor(tokens int) lipgloss.Color {
	switch {
	case tokens >= 230_000:
		return lipgloss.Color("#ff0000")
	case tokens >= 175_000:
		return ColorDanger
	case tokens >= 150_000:
		return ColorWarning
	default:
		return ColorSuccess
	}
}

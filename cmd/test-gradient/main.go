package main

import (
	"fmt"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/style"
	"github.com/huyhandes/cc-hud-go/theme"
)

func main() {
	// Initialize style system
	cfg := config.Default()
	themeInstance := theme.LoadThemeFromConfig(cfg.Theme, cfg.Colors)
	style.Init(themeInstance)

	fmt.Println("Static Gradient Progress Bar Test")
	fmt.Println("===================================")
	fmt.Println()

	testPercentages := []float64{0, 2, 5, 10, 25, 40, 50, 60, 75, 85, 95, 100}

	for _, pct := range testPercentages {
		bar := style.RenderGradientBar(pct, 10)
		fmt.Printf("%3.0f%% │ %s\n", pct, bar)
	}

	fmt.Println("\n===================================")
	fmt.Println("Notice how the gradient transitions:")
	fmt.Println("  0-50%:  Green → Yellow")
	fmt.Println(" 50-75%:  Yellow → Orange")
	fmt.Println("75-100%:  Orange → Red")
	fmt.Println("\nOnly the filled portion is shown!")
}

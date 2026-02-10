package main

import (
	"fmt"

	"github.com/huybui/cc-hud-go/internal/oauth"
)

func main() {
	fmt.Println("Testing OAuth API...")

	usage, err := oauth.FetchUsage()
	if err != nil {
		fmt.Printf("Error fetching usage: %v\n", err)
		return
	}

	fmt.Printf("5-hour usage: %.1f%%\n", usage.FiveHour.Utilization)
	fmt.Printf("5-hour resets at: %s\n", usage.FiveHour.ResetsAt)
	fmt.Printf("7-day usage: %.1f%%\n", usage.SevenDay.Utilization)
	fmt.Printf("7-day resets at: %s\n", usage.SevenDay.ResetsAt)
}

package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// UsageResponse represents the API response from oauth/usage endpoint
type UsageResponse struct {
	FiveHour struct {
		Utilization float64   `json:"utilization"`
		ResetsAt    time.Time `json:"resets_at"`
	} `json:"five_hour"`
	SevenDay struct {
		Utilization float64   `json:"utilization"`
		ResetsAt    time.Time `json:"resets_at"`
	} `json:"seven_day"`
}

// GetAccessToken retrieves the OAuth access token from system keychain
func GetAccessToken() (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// macOS: use security command
		cmd = exec.Command("security", "find-generic-password", "-s", "Claude Code-credentials", "-w")
	case "linux":
		// Linux: use secret-tool
		cmd = exec.Command("secret-tool", "lookup", "service", "Claude Code-credentials")
	case "windows":
		// Windows: not yet supported
		return "", fmt.Errorf("windows keychain access not yet implemented")
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve access token: %w", err)
	}

	token := strings.TrimSpace(string(output))
	if token == "" {
		return "", fmt.Errorf("empty access token retrieved")
	}

	// Parse JSON (security command returns JSON with nested structure)
	if strings.HasPrefix(token, "{") {
		// Try nested structure first (claudeAiOauth.accessToken)
		var nestedCreds struct {
			ClaudeAiOauth struct {
				AccessToken  string `json:"accessToken"`
				RefreshToken string `json:"refreshToken"`
			} `json:"claudeAiOauth"`
		}
		if err := json.Unmarshal([]byte(token), &nestedCreds); err == nil && nestedCreds.ClaudeAiOauth.AccessToken != "" {
			return nestedCreds.ClaudeAiOauth.AccessToken, nil
		}

		// Try flat structure as fallback
		var creds struct {
			AccessToken string `json:"accessToken"`
		}
		if err := json.Unmarshal([]byte(token), &creds); err == nil && creds.AccessToken != "" {
			return creds.AccessToken, nil
		}
	}

	return token, nil
}

// FetchUsage retrieves rate limit usage from Anthropic OAuth API
func FetchUsage() (*UsageResponse, error) {
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", "https://api.anthropic.com/api/oauth/usage", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("anthropic-beta", "oauth-2025-04-20")
	req.Header.Set("Accept", "application/json")

	// Execute request
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch usage: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var usage UsageResponse
	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &usage, nil
}

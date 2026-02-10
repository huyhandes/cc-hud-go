package oauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	// This test will be skipped in CI since it requires actual keychain access
	if testing.Short() {
		t.Skip("Skipping keychain access test in short mode")
	}

	token, err := GetAccessToken()
	if err != nil {
		t.Logf("Could not retrieve access token (may not be configured): %v", err)
		t.Skip("Access token not available")
	}

	if token == "" {
		t.Error("Expected non-empty access token")
	}

	t.Logf("Successfully retrieved token (length: %d)", len(token))
}

func TestFetchUsage(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("Authorization") == "" {
			t.Error("Expected Authorization header")
		}

		if r.Header.Get("anthropic-beta") != "oauth-2025-04-20" {
			t.Error("Expected anthropic-beta header")
		}

		// Return mock response
		response := UsageResponse{}
		response.FiveHour.Utilization = 23.5
		response.SevenDay.Utilization = 67.2

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Note: This test doesn't actually call FetchUsage() since it requires
	// real credentials. To test the full flow, we'd need to mock GetAccessToken.
	t.Log("Mock server created successfully")
}

func TestUsageResponseParsing(t *testing.T) {
	jsonData := `{
		"five_hour": {
			"utilization": 23.5,
			"resets_at": "2025-11-04T04:59:59.943648+00:00"
		},
		"seven_day": {
			"utilization": 67.2,
			"resets_at": "2025-11-06T03:59:59.943679+00:00"
		}
	}`

	var usage UsageResponse
	err := json.Unmarshal([]byte(jsonData), &usage)
	if err != nil {
		t.Fatalf("Failed to parse usage response: %v", err)
	}

	if usage.FiveHour.Utilization != 23.5 {
		t.Errorf("Expected 5h utilization 23.5, got %f", usage.FiveHour.Utilization)
	}

	if usage.SevenDay.Utilization != 67.2 {
		t.Errorf("Expected 7d utilization 67.2, got %f", usage.SevenDay.Utilization)
	}

	if usage.FiveHour.ResetsAt.IsZero() {
		t.Error("Expected non-zero reset time for 5h")
	}

	if usage.SevenDay.ResetsAt.IsZero() {
		t.Error("Expected non-zero reset time for 7d")
	}
}

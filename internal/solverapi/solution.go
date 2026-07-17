package solverapi

import "github.com/aarock1234/unicap"

// mapStatus converts a provider status string to a unicap task status.
func mapStatus(status string) unicap.TaskStatus {
	switch status {
	case "processing":
		return unicap.TaskStatusProcessing
	case "ready":
		return unicap.TaskStatusReady
	case "failed":
		return unicap.TaskStatusFailed
	default:
		return unicap.TaskStatusPending
	}
}

// mapSolution extracts the common solution fields from a provider's raw
// solution payload. The full payload is always preserved in Extra.
func mapSolution(solution map[string]any) unicap.Solution {
	sol := unicap.Solution{Extra: solution}

	if token, ok := solution["gRecaptchaResponse"].(string); ok {
		sol.Token = token
	} else if token, ok := solution["token"].(string); ok {
		sol.Token = token
	}

	if cookie, ok := solution["cookie"].(string); ok {
		sol.Cookie = cookie
	}

	if text, ok := solution["text"].(string); ok {
		sol.Text = text
	}

	return sol
}

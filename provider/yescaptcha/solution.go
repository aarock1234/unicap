package yescaptcha

import (
	"encoding/json"

	"github.com/aarock1234/unicap/v3"
)

// challengeSolution is the part of YesCaptcha's Cloudflare challenge
// solution the shared decoder does not read: it names the user agent in
// snake case.
type challengeSolution struct {
	UserAgent string `json:"user_agent"`
}

// refineSolution fills solution.UserAgent from the snake case user_agent
// field when the shared decoder found no user agent. A payload that is not a
// solution object leaves solution unchanged.
func refineSolution(raw json.RawMessage, solution *unicap.Solution) {
	if solution.UserAgent != "" {
		return
	}

	var challenge challengeSolution
	if err := json.Unmarshal(raw, &challenge); err != nil {
		return
	}

	solution.UserAgent = challenge.UserAgent
}

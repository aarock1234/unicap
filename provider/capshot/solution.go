package capshot

import (
	"encoding/json"

	"github.com/aarock1234/unicap/v3"
)

// wafSolution is the part of Capshot's AWS WAF solution the shared decoder
// does not read: the token arrives in a field of its own.
type wafSolution struct {
	AWSWAFToken string `json:"awsWafToken"`
}

// refineSolution fills solution.Token from the awsWafToken field when the
// shared decoder found no token. A payload that is not a solution object
// leaves solution unchanged.
func refineSolution(raw json.RawMessage, solution *unicap.Solution) {
	if solution.Token != "" {
		return
	}

	var waf wafSolution
	if err := json.Unmarshal(raw, &waf); err != nil {
		return
	}

	solution.Token = waf.AWSWAFToken
}

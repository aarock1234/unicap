package nonecap

import (
	"encoding/json"
	"time"

	"github.com/aarock1234/unicap/v3"
)

// newResult builds the unicap result of a solve that is in progress or
// solved; raw is the solve object as received. The cost stays zero because
// NoneCap charges in credits, not US dollars. The solution is set only once
// the solve is solved.
func newResult(solve solveObject, raw json.RawMessage) *unicap.Result {
	result := &unicap.Result{
		TaskID:    solve.ID,
		Status:    taskStatus(solve.Status),
		CreatedAt: timestamp(solve.CreatedAt),
		SolvedAt:  timestamp(solve.FinishedAt),
	}

	if result.Ready() {
		result.Solution = newSolution(solve, raw)
	}

	return result
}

// newSolution converts a solved solve. NoneCap has no solution object, so
// Raw is the whole solve object, which also carries hCaptcha's resp_key and
// any warnings.
func newSolution(solve solveObject, raw json.RawMessage) unicap.Solution {
	return unicap.Solution{
		Token:     solve.Token,
		UserAgent: solve.UserAgent,
		Raw:       raw,
	}
}

// timestamp parses an RFC 3339 time. A missing or malformed time is the zero
// time, which unicap reads as "unknown".
func timestamp(text string) time.Time {
	parsed, err := time.Parse(time.RFC3339, text)
	if err != nil {
		return time.Time{}
	}

	return parsed
}

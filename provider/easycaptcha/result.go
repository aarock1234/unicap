package easycaptcha

import (
	"encoding/json"
	"time"

	"github.com/aarock1234/unicap/v3"
)

// newResult builds the unicap result of a task that is queued, running or
// completed. The cost stays zero because the task object carries none. The
// solution is decoded only once the task is completed.
func newResult(task taskObject) *unicap.Result {
	result := &unicap.Result{
		TaskID:    task.ID,
		Status:    taskStatus(task.Status),
		CreatedAt: timestamp(task.CreatedAt),
		SolvedAt:  timestamp(task.CompletedAt),
	}

	if result.Ready() {
		result.Solution = decodeSolution(task.Result)
	}

	return result
}

// decodeSolution converts the result of a completed task and keeps the
// result itself in Raw. The result of a token task is an object holding the
// token; one of another shape, such as a raw recognition task's, leaves the
// token empty.
func decodeSolution(raw json.RawMessage) unicap.Solution {
	var decoded tokenResult

	_ = json.Unmarshal(raw, &decoded)

	return unicap.Solution{
		Token: decoded.Token,
		Raw:   raw,
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

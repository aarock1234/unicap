package nopecha

import (
	"encoding/json"

	"github.com/aarock1234/unicap/v3"
)

// newResult builds the unicap result of a job that is not finished. The cost
// and the timestamps stay zero because NopeCHA charges in credits and
// reports no times.
func newResult(job jobRef, status unicap.TaskStatus) *unicap.Result {
	return &unicap.Result{
		TaskID: job.String(),
		Status: status,
	}
}

// newReadyResult builds the unicap result of a finished job; data is its
// solution as received.
func newReadyResult(job jobRef, data json.RawMessage) *unicap.Result {
	result := newResult(job, unicap.TaskStatusReady)
	result.Solution = decodeSolution(data)

	return result
}

// decodeSolution converts the solution of a finished job and keeps the data
// itself in Raw. A token endpoint answers with the token as a string; data
// of another shape, such as a recognition endpoint's, leaves the token
// empty.
func decodeSolution(data json.RawMessage) unicap.Solution {
	var token string

	_ = json.Unmarshal(data, &token)

	return unicap.Solution{
		Token: token,
		Raw:   data,
	}
}

package easycaptcha

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Paths of the native v1 API.
const (
	taskCollectionPath = "/v1/task"
	balancePath        = "/v1/balance"
)

// Task statuses that decide how a task is read. Every other status, such as
// "running", means the task is being solved.
const (
	statusQueued    = "queued"
	statusCompleted = "completed"
	statusFailed    = "failed"
)

// The HTTP statuses each endpoint answers with on success. Any other status
// is an error, even another 2xx.
var (
	// createStatuses are 200 for a task that ended while the request was
	// held and 202 for one that is queued or running.
	createStatuses = []int{
		http.StatusOK,
		http.StatusAccepted,
	}

	// fetchStatuses holds 200, which a fetched task comes with whatever its
	// state.
	fetchStatuses = []int{
		http.StatusOK,
	}

	// balanceStatuses holds 200, which the balance comes with.
	balanceStatuses = []int{
		http.StatusOK,
	}

	// feedbackStatuses holds 204, which recorded feedback is answered with,
	// without a body.
	feedbackStatuses = []int{
		http.StatusNoContent,
	}
)

// taskPath returns the path of one task. The ID is escaped so that it stays
// a single path segment.
func taskPath(taskID string) string {
	return taskCollectionPath + "/" + url.PathEscape(taskID)
}

// feedbackPath returns the path that takes feedback on one task.
func feedbackPath(taskID string) string {
	return taskPath(taskID) + "/feedback"
}

// taskObject is the task easycaptcha.to returns when one is created or
// fetched. Its timestamps are RFC 3339, and the fields that are null until
// the task ends decode as empty. Result is kept as received because its
// shape depends on the task type.
type taskObject struct {
	ID          string          `json:"task_id"`
	Status      string          `json:"status"`
	Result      json.RawMessage `json:"result"`
	Error       problem         `json:"error"`
	CreatedAt   string          `json:"created_at"`
	CompletedAt string          `json:"completed_at"`
}

// tokenResult is the result object of a completed token task.
type tokenResult struct {
	Token string `json:"token"`
}

// problem is the RFC 9457 problem details body of a request error, and the
// error object embedded in a failed task. Code is the short form of the
// problem type, and Detail is optional.
type problem struct {
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// balance is the account balance. Microcents is a signed decimal integer
// string in hundred-millionths of one unit of Currency.
type balance struct {
	Currency   string `json:"currency"`
	Microcents string `json:"balance_microcents"`
}

// feedbackRequest tells easycaptcha.to whether a task's result worked.
type feedbackRequest struct {
	Success bool `json:"success"`
}

package noptcha

import (
	"encoding/json"
	"regexp"
)

// Paths of the native v1 API.
const (
	tasksPath   = "/v1/tasks"
	balancePath = "/v1/balance"
)

// Final task statuses. Every other status, such as "processing", means the
// task is still being solved.
const (
	statusSolved = "solved"
	statusFailed = "failed"
)

// taskIDPattern is the form Noptcha documents for a task ID in a path.
var taskIDPattern = regexp.MustCompile(`^[0-9]{1,16}$`)

// taskPath returns the path of one task.
func taskPath(taskID string) string {
	return tasksPath + "/" + taskID
}

// feedbackPath returns the path that takes feedback on one task.
func feedbackPath(taskID string) string {
	return taskPath(taskID) + "/feedback"
}

// taskRequest is the create-task body: a native task type and the params
// object of that type.
type taskRequest[P any] struct {
	Type   string `json:"type"`
	Params P      `json:"params"`
}

// taskObject is the task Noptcha returns when one is created or fetched. Its
// timestamps are RFC 3339, and CompletedAt is empty while the task is
// processing.
type taskObject struct {
	ID          string          `json:"id"`
	Status      string          `json:"status"`
	Solution    json.RawMessage `json:"solution"`
	Error       taskError       `json:"error"`
	Cost        money           `json:"cost"`
	CreatedAt   string          `json:"created_at"`
	CompletedAt string          `json:"completed_at"`
}

// taskError explains why a task failed. It is empty unless the task's status
// is "failed".
type taskError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// money is an amount in US dollars, written as a decimal string with six
// places.
type money struct {
	Amount string `json:"amount"`
}

// problem is the RFC 9457 problem details body of a request error. Detail is
// optional.
type problem struct {
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// balance is the account balance. Available is what new tasks can spend, as
// a decimal string of US dollars.
type balance struct {
	Available string `json:"available"`
}

// feedbackRequest tells Noptcha whether a task's solution worked.
type feedbackRequest struct {
	Correct bool `json:"correct"`
}

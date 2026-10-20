package unicap

import (
	"time"
)

// TaskStatus is the state of a submitted task.
type TaskStatus string

const (
	// TaskStatusPending reports that the provider accepted the task but has
	// not started solving it.
	TaskStatusPending TaskStatus = "pending"

	// TaskStatusProcessing reports that the task is being solved.
	TaskStatusProcessing TaskStatus = "processing"

	// TaskStatusReady reports that the solution is available.
	TaskStatusReady TaskStatus = "ready"
)

// Result is the state of a submitted task and, once ready, its solution.
// Provider failures are returned as errors rather than as a result status.
type Result struct {
	// TaskID is the provider's identifier for the task.
	TaskID string

	// Status is the task's current state.
	Status TaskStatus

	// Solution is the solved captcha. It is set only when Status is
	// TaskStatusReady.
	Solution Solution

	// Cost is the amount charged in US dollars, or zero when the provider
	// does not report it.
	Cost float64

	// CreatedAt is when the provider created the task, or zero when unknown.
	CreatedAt time.Time

	// SolvedAt is when the provider finished the task, or zero when unknown.
	SolvedAt time.Time

	// Task is the submitted task. Client sets it so the result can be
	// reported later; it is nil for results fetched by task ID alone.
	Task Task
}

// Ready reports whether the solution is available.
func (r *Result) Ready() bool {
	return r.Status == TaskStatusReady
}

// Elapsed returns how long the provider took to solve the task, or zero when
// the provider does not report timestamps.
func (r *Result) Elapsed() time.Duration {
	if r.CreatedAt.IsZero() || r.SolvedAt.IsZero() {
		return 0
	}

	return r.SolvedAt.Sub(r.CreatedAt)
}

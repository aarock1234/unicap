package unicap

import (
	"context"
)

// Provider is a captcha solving service.
//
// Provider failures must be returned as *Error so callers and the poller can
// tell final failures apart from transient transport errors.
type Provider interface {
	// Name returns the provider's identifier, such as "capsolver".
	Name() string

	// CreateTask submits a task. The returned result carries the task ID and,
	// for providers that solve some tasks synchronously, may already be
	// ready.
	CreateTask(ctx context.Context, task Task) (*Result, error)

	// GetTaskResult returns the current state of a submitted task.
	GetTaskResult(ctx context.Context, taskID string) (*Result, error)
}

// BalanceChecker is implemented by providers that report the account
// balance.
type BalanceChecker interface {
	// Balance returns the account balance in US dollars.
	Balance(ctx context.Context) (float64, error)
}

// Reporter is implemented by providers that accept feedback on whether a
// solution worked. Reports help providers refund bad answers and improve
// worker quality.
type Reporter interface {
	// Report tells the provider whether result's solution was accepted by
	// the target site.
	Report(ctx context.Context, result *Result, correct bool) error
}

// TaskSupporter is implemented by providers that can tell in advance whether
// they handle a task type.
type TaskSupporter interface {
	// Supports reports whether the provider can solve tasks of the given
	// type.
	Supports(taskType TaskType) bool
}

// PollingAdvisor is implemented by providers that recommend their own
// polling pace. Client uses the recommendation unless WithPolling overrides
// it.
type PollingAdvisor interface {
	// RecommendedPolling returns the provider's recommended polling settings.
	RecommendedPolling() Polling
}

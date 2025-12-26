package upicap

import (
	"context"
)

// Provider defines the interface for captcha solving providers
type Provider interface {
	// CreateTask submits a captcha task and returns a task ID
	CreateTask(ctx context.Context, task Task) (string, error)

	// GetTaskResult retrieves the result for a given task ID
	GetTaskResult(ctx context.Context, taskID string) (*TaskResult, error)

	// Name returns the provider's identifier
	Name() string
}

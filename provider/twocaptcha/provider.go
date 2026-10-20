package twocaptcha

import (
	"context"

	"github.com/aarock1234/unicap/v3"
)

// Name returns the provider identifier.
func (c *Client) Name() string {
	return c.api.Name()
}

// CreateTask submits task.
func (c *Client) CreateTask(ctx context.Context, task unicap.Task) (*unicap.Result, error) {
	return c.api.CreateTask(ctx, task)
}

// GetTaskResult returns the current state of a submitted task.
func (c *Client) GetTaskResult(ctx context.Context, taskID string) (*unicap.Result, error) {
	return c.api.GetTaskResult(ctx, taskID)
}

// Balance returns the account balance in US dollars.
func (c *Client) Balance(ctx context.Context) (float64, error) {
	return c.api.Balance(ctx)
}

// Supports reports whether 2Captcha solves tasks of the given type.
func (c *Client) Supports(taskType unicap.TaskType) bool {
	return c.api.Supports(taskType)
}

// RecommendedPolling returns 2Captcha's recommended polling pace.
func (c *Client) RecommendedPolling() unicap.Polling {
	return c.api.RecommendedPolling()
}

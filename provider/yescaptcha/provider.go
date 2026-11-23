package yescaptcha

import (
	"context"

	"github.com/aarock1234/unicap/v3"
)

// Name returns the provider identifier.
func (c *Client) Name() string {
	return c.api.Name()
}

// CreateTask submits task. reCAPTCHA image classification tasks come back
// already ready.
func (c *Client) CreateTask(ctx context.Context, task unicap.Task) (*unicap.Result, error) {
	return c.api.CreateTask(ctx, task)
}

// GetTaskResult returns the current state of a submitted task. Query the API
// root that created the task; the two roots do not share tasks.
func (c *Client) GetTaskResult(ctx context.Context, taskID string) (*unicap.Result, error) {
	return c.api.GetTaskResult(ctx, taskID)
}

// Supports reports whether YesCaptcha solves tasks of the given type.
func (c *Client) Supports(taskType unicap.TaskType) bool {
	return c.api.Supports(taskType)
}

// RecommendedPolling returns YesCaptcha's recommended polling pace.
func (c *Client) RecommendedPolling() unicap.Polling {
	return c.api.RecommendedPolling()
}

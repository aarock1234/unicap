package solvecha

import (
	"context"

	"github.com/aarock1234/unicap/v3"
)

// Name returns the provider identifier.
func (c *Client) Name() string {
	return c.api.Name()
}

// CreateTask submits task. Picture and text tasks usually come back ready.
func (c *Client) CreateTask(ctx context.Context, task unicap.Task) (*unicap.Result, error) {
	return c.api.CreateTask(ctx, task)
}

// GetTaskResult returns the current state of a submitted task.
func (c *Client) GetTaskResult(ctx context.Context, taskID string) (*unicap.Result, error) {
	return c.api.GetTaskResult(ctx, taskID)
}

// Balance returns the prepaid balance. Plan quotas and the free quota are
// not part of it, so a zero balance does not mean that solving is blocked.
func (c *Client) Balance(ctx context.Context) (float64, error) {
	return c.api.Balance(ctx)
}

// Supports reports whether Solvecha solves tasks of the given type.
func (c *Client) Supports(taskType unicap.TaskType) bool {
	return c.api.Supports(taskType)
}

// RecommendedPolling returns the recommended polling pace.
func (c *Client) RecommendedPolling() unicap.Polling {
	return c.api.RecommendedPolling()
}

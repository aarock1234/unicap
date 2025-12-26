package unicap

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

// Client is the main interface for solving captchas
type Client struct {
	provider Provider
	logger   *slog.Logger
	poller   *Poller
}

// NewClient creates a new captcha solving client
func NewClient(provider Provider, opts ...Option) (*Client, error) {
	if provider == nil {
		return nil, fmt.Errorf("provider cannot be nil")
	}

	c := &Client{
		provider: provider,
		logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.poller == nil {
		c.poller = NewPoller(provider, DefaultPollerConfig())
	}

	return c, nil
}

// Solve submits a task and automatically polls for the result (synchronous)
func (c *Client) Solve(ctx context.Context, task Task) (*Solution, error) {
	taskID, err := c.provider.CreateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("creating task: %w", err)
	}

	c.logger.InfoContext(ctx, "task created",
		slog.String("task_id", taskID),
		slog.String("task_type", string(task.Type())),
		slog.String("provider", c.provider.Name()),
	)

	result, err := c.poller.Poll(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("polling task %s: %w", taskID, err)
	}

	return &result.Solution, nil
}

// CreateTask submits a task without polling (asynchronous)
func (c *Client) CreateTask(ctx context.Context, task Task) (string, error) {
	taskID, err := c.provider.CreateTask(ctx, task)
	if err != nil {
		return "", fmt.Errorf("creating task: %w", err)
	}

	c.logger.InfoContext(ctx, "task created",
		slog.String("task_id", taskID),
		slog.String("task_type", string(task.Type())),
		slog.String("provider", c.provider.Name()),
	)

	return taskID, nil
}

// GetTaskResult retrieves a task result by ID (asynchronous)
func (c *Client) GetTaskResult(ctx context.Context, taskID string) (*TaskResult, error) {
	result, err := c.provider.GetTaskResult(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("getting task result %s: %w", taskID, err)
	}

	return result, nil
}

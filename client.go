// Package unicap is a universal captcha solving SDK. A Client submits tasks
// from the tasks package to any Provider, polls for the solution, and
// returns typed results and errors that are consistent across providers.
package unicap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

// Client submits captcha tasks to a provider and retrieves their solutions.
// It is safe for concurrent use.
type Client struct {
	provider Provider
	logger   *slog.Logger

	polling      Polling
	pollingSet   bool
	solveRetries int
}

// New creates a client for provider. It returns an error when an option
// supplies unusable settings.
func New(provider Provider, opts ...Option) (*Client, error) {
	if provider == nil {
		return nil, ErrNilProvider
	}

	c := &Client{
		provider: provider,
		logger:   slog.New(slog.DiscardHandler),
	}

	for _, opt := range opts {
		opt(c)
	}

	if !c.pollingSet {
		c.polling = recommendedPolling(provider)
	}

	if err := c.polling.Validate(); err != nil {
		return nil, err
	}

	if c.solveRetries < 0 {
		return nil, errors.New("solve retries must not be negative")
	}

	return c, nil
}

// Provider returns the client's provider.
func (c *Client) Provider() Provider {
	return c.provider
}

// Solve submits task and waits until its solution is ready, the context is
// canceled, or polling times out. With WithSolveRetries it resubmits the task
// after retryable failures.
func (c *Client) Solve(ctx context.Context, task Task) (*Result, error) {
	for attempt := 0; ; attempt++ {
		result, err := c.solveOnce(ctx, task)
		if err == nil {
			return result, nil
		}

		if attempt >= c.solveRetries || !IsRetryable(err) || ctx.Err() != nil {
			return nil, err
		}

		c.logger.WarnContext(ctx, "retrying task after retryable failure",
			slog.String("provider", c.provider.Name()),
			slog.String("task_type", string(task.Type())),
			slog.Int("attempt", attempt+1),
			slog.Any("error", err),
		)
	}
}

// CreateTask validates and submits task without waiting. The result carries
// the task ID; for synchronous recognition tasks it may already be ready.
func (c *Client) CreateTask(ctx context.Context, task Task) (*Result, error) {
	if task == nil {
		return nil, fmt.Errorf("task is nil: %w", ErrInvalidTask)
	}

	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("validating %s task: %w", task.Type(), err)
	}

	result, err := c.provider.CreateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("creating %s task: %w", task.Type(), err)
	}

	result.Task = task

	c.logger.DebugContext(ctx, "task created",
		slog.String("provider", c.provider.Name()),
		slog.String("task_type", string(task.Type())),
		slog.String("task_id", result.TaskID),
		slog.String("status", string(result.Status)),
	)

	return result, nil
}

// GetTaskResult returns the current state of a task without waiting.
func (c *Client) GetTaskResult(ctx context.Context, taskID string) (*Result, error) {
	result, err := c.provider.GetTaskResult(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("getting task result %s: %w", taskID, err)
	}

	return result, nil
}

// Wait polls a submitted task until its solution is ready, the context is
// canceled, or polling times out.
func (c *Client) Wait(ctx context.Context, taskID string) (*Result, error) {
	return c.wait(ctx, taskID)
}

// Balance returns the provider account balance in US dollars. It returns
// ErrUnsupportedOperation when the provider does not report balances.
func (c *Client) Balance(ctx context.Context) (float64, error) {
	checker, ok := c.provider.(BalanceChecker)
	if !ok {
		return 0, fmt.Errorf("%s balance: %w", c.provider.Name(), ErrUnsupportedOperation)
	}

	balance, err := checker.Balance(ctx)
	if err != nil {
		return 0, fmt.Errorf("getting balance: %w", err)
	}

	return balance, nil
}

// Report tells the provider whether result's solution was accepted by the
// target site. It returns ErrUnsupportedOperation when the provider does not
// accept reports.
func (c *Client) Report(ctx context.Context, result *Result, correct bool) error {
	reporter, ok := c.provider.(Reporter)
	if !ok {
		return fmt.Errorf("%s reports: %w", c.provider.Name(), ErrUnsupportedOperation)
	}

	if result == nil || result.TaskID == "" {
		return fmt.Errorf("report needs a result with a task id: %w", ErrInvalidTask)
	}

	if err := reporter.Report(ctx, result, correct); err != nil {
		return fmt.Errorf("reporting task %s: %w", result.TaskID, err)
	}

	return nil
}

// Supports reports whether the provider can solve tasks of the given type.
// Providers that do not declare their capabilities are assumed to support
// every type.
func (c *Client) Supports(taskType TaskType) bool {
	supporter, ok := c.provider.(TaskSupporter)
	if !ok {
		return true
	}

	return supporter.Supports(taskType)
}

// solveOnce submits task and waits for its solution.
func (c *Client) solveOnce(ctx context.Context, task Task) (*Result, error) {
	created, err := c.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	if created.Ready() {
		return created, nil
	}

	result, err := c.wait(ctx, created.TaskID)
	if err != nil {
		return nil, err
	}

	result.Task = task

	return result, nil
}

// recommendedPolling returns the provider's polling recommendation, or the
// defaults when it has none.
func recommendedPolling(provider Provider) Polling {
	advisor, ok := provider.(PollingAdvisor)
	if !ok {
		return DefaultPolling()
	}

	return advisor.RecommendedPolling()
}

package unicap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// wait polls the provider until the task is ready. Transient failures are
// tolerated up to Polling.MaxTransientErrors in a row; any other provider
// error ends the wait immediately.
func (c *Client) wait(ctx context.Context, taskID string) (*Result, error) {
	if c.polling.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.polling.Timeout)
		defer cancel()
	}

	interval := c.polling.Interval
	transientErrors := 0

	timer := time.NewTimer(c.polling.InitialDelay)
	defer timer.Stop()

	for {
		if err := sleep(ctx, timer); err != nil {
			return nil, fmt.Errorf("waiting for task %s: %w", taskID, err)
		}

		result, err := c.provider.GetTaskResult(ctx, taskID)

		switch {
		case err != nil && ctx.Err() != nil:
			return nil, fmt.Errorf("waiting for task %s: %w", taskID, contextError(ctx))

		case err != nil && isTransient(err) && transientErrors < c.polling.MaxTransientErrors:
			transientErrors++

			c.logger.DebugContext(ctx, "transient poll failure, retrying",
				slog.String("task_id", taskID),
				slog.Int("consecutive_errors", transientErrors),
				slog.Any("error", err),
			)

		case err != nil:
			return nil, fmt.Errorf("waiting for task %s: %w", taskID, err)

		case result.Ready():
			c.logger.DebugContext(ctx, "task ready",
				slog.String("provider", c.provider.Name()),
				slog.String("task_id", taskID),
			)

			return result, nil

		default:
			transientErrors = 0

			c.logger.DebugContext(ctx, "task not ready",
				slog.String("task_id", taskID),
				slog.String("status", string(result.Status)),
				slog.Duration("next_check", interval),
			)
		}

		timer.Reset(interval)
		interval = c.polling.next(interval)
	}
}

// sleep blocks until timer fires or ctx is done.
func sleep(ctx context.Context, timer *time.Timer) error {
	select {
	case <-ctx.Done():
		return contextError(ctx)
	case <-timer.C:
		return nil
	}
}

// contextError translates a deadline into ErrTimeout while keeping the
// context error inspectable.
func contextError(ctx context.Context) error {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return errors.Join(ErrTimeout, ctx.Err())
	}

	return ctx.Err()
}

package unicap

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"
)

// maxPollErrors is the number of consecutive transient result-fetch failures
// tolerated before polling gives up.
const maxPollErrors = 3

// Poller repeatedly checks a provider for task completion.
type Poller struct {
	provider Provider
	config   PollerConfig
	logger   *slog.Logger
}

// PollerConfig defines polling behavior.
type PollerConfig struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Timeout         time.Duration
	Multiplier      float64
}

// DefaultPollerConfig returns sensible defaults for polling.
func DefaultPollerConfig() PollerConfig {
	return PollerConfig{
		InitialInterval: 2 * time.Second,
		MaxInterval:     15 * time.Second,
		Timeout:         5 * time.Minute,
		Multiplier:      1.5,
	}
}

// PollerOption configures a Poller.
type PollerOption func(*Poller)

// WithPollerLogger sets the poller's logger.
func WithPollerLogger(logger *slog.Logger) PollerOption {
	return func(p *Poller) {
		if logger != nil {
			p.logger = logger
		}
	}
}

// NewPoller creates a poller for the given provider and config.
func NewPoller(provider Provider, config PollerConfig, opts ...PollerOption) *Poller {
	p := &Poller{
		provider: provider,
		config:   config,
		logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Poll blocks until the task reaches a terminal state, the context is
// cancelled, or the configured timeout elapses. Transient result-fetch
// failures are tolerated up to maxPollErrors before the error is returned.
func (p *Poller) Poll(ctx context.Context, taskID string) (*TaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, p.config.Timeout)
	defer cancel()

	interval := p.config.InitialInterval
	consecutiveErrors := 0

	// Fire immediately on the first iteration, then back off between checks so
	// an already-solved task returns without waiting a full interval.
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil, fmt.Errorf("polling task %s: %w", taskID, ErrTimeout)
			}

			return nil, fmt.Errorf("polling task %s: %w", taskID, ctx.Err())
		case <-timer.C:
		}

		result, err := p.provider.GetTaskResult(ctx, taskID)
		if err != nil {
			consecutiveErrors++
			if consecutiveErrors > maxPollErrors {
				return nil, fmt.Errorf("polling task %s: %w", taskID, err)
			}

			p.logger.DebugContext(ctx, "transient poll error, retrying",
				slog.String("task_id", taskID),
				slog.Int("consecutive_errors", consecutiveErrors),
				slog.Any("error", err),
			)

			timer.Reset(interval)

			continue
		}

		consecutiveErrors = 0

		switch result.Status {
		case TaskStatusReady:
			p.logger.InfoContext(ctx, "task completed",
				slog.String("task_id", taskID),
				slog.String("provider", p.provider.Name()),
			)

			return result, nil
		case TaskStatusFailed:
			if result.Error != nil {
				return nil, result.Error
			}

			return nil, fmt.Errorf("polling task %s: %w", taskID, ErrInvalidTask)
		case TaskStatusPending, TaskStatusProcessing:
			interval = nextInterval(interval, p.config)

			p.logger.DebugContext(ctx, "task still processing",
				slog.String("task_id", taskID),
				slog.String("status", string(result.Status)),
				slog.Duration("next_check", interval),
			)

			timer.Reset(interval)
		}
	}
}

// nextInterval grows the poll interval by the configured multiplier, capped at
// MaxInterval.
func nextInterval(current time.Duration, config PollerConfig) time.Duration {
	next := time.Duration(float64(current) * config.Multiplier)
	if next > config.MaxInterval {
		return config.MaxInterval
	}

	return next
}

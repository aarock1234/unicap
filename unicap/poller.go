package unicap

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

// Poller handles automatic polling for task results
type Poller struct {
	provider Provider
	config   PollerConfig
	logger   *slog.Logger
}

// PollerConfig defines polling behavior
type PollerConfig struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Timeout         time.Duration
	Multiplier      float64
}

// DefaultPollerConfig returns sensible defaults for polling
func DefaultPollerConfig() PollerConfig {
	return PollerConfig{
		InitialInterval: 2 * time.Second,
		MaxInterval:     15 * time.Second,
		Timeout:         5 * time.Minute,
		Multiplier:      1.5,
	}
}

// NewPoller creates a new poller with the given provider and config
func NewPoller(provider Provider, config PollerConfig) *Poller {
	return &Poller{
		provider: provider,
		config:   config,
		logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

// Poll continuously checks for task completion
func (p *Poller) Poll(ctx context.Context, taskID string) (*TaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, p.config.Timeout)
	defer cancel()

	interval := p.config.InitialInterval

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("polling task %s: %w", taskID, ErrTimeout)
		case <-time.After(interval):
			result, err := p.provider.GetTaskResult(ctx, taskID)
			if err != nil {
				return nil, fmt.Errorf("getting task result: %w", err)
			}

			switch result.Status {
			case TaskStatusReady:
				p.logger.InfoContext(ctx, "task completed",
					slog.String("task_id", taskID),
					slog.String("provider", p.provider.Name()),
				)
				return result, nil
			case TaskStatusFailed:
				p.logger.ErrorContext(ctx, "task failed",
					slog.String("task_id", taskID),
					slog.Any("error", result.Error),
				)
				return nil, result.Error
			case TaskStatusPending, TaskStatusProcessing:
				interval = time.Duration(float64(interval) * p.config.Multiplier)
				if interval > p.config.MaxInterval {
					interval = p.config.MaxInterval
				}

				p.logger.DebugContext(ctx, "task still processing",
					slog.String("task_id", taskID),
					slog.String("status", string(result.Status)),
					slog.Duration("next_check", interval),
				)
			}
		}
	}
}

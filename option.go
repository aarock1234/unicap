package unicap

import (
	"log/slog"
)

// Option configures a Client.
type Option func(*Client)

// WithLogger sets the logger for task lifecycle events. By default the
// client discards its logs.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// WithPolling overrides the polling settings. Without it the client uses the
// provider's recommendation, or DefaultPolling when the provider has none.
func WithPolling(polling Polling) Option {
	return func(c *Client) {
		c.polling = polling
		c.pollingSet = true
	}
}

// WithSolveRetries makes Solve submit a fresh task up to retries more times
// when an attempt fails with a retryable error such as ErrUnsolvable or
// ErrUnavailable. Each attempt may be billed by the provider. Default: 0.
func WithSolveRetries(retries int) Option {
	return func(c *Client) {
		c.solveRetries = retries
	}
}
